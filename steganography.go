package gostegano

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
)

const (
	magicBytes = "GOST"
	headerSize = 4 + len(magicBytes)
)

type Steganography struct {
	sourceImage image.Image
	targetImage *image.NRGBA // Only used in Encode mode
}

func NewSteganography(img image.Image) *Steganography {
	return &Steganography{sourceImage: img}
}

func NewSteganographyFromReader(reader io.Reader) (*Steganography, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return &Steganography{sourceImage: img}, nil
}

func (s *Steganography) Decode() ([]byte, error) {
	header := make([]byte, headerSize)
	for i, pixel := range IteratePixels(s.sourceImage) {
		if i >= headerSize {
			break
		}
		header[i] = s.decodePixel(pixel)
	}
	if string(header[:4]) != magicBytes {
		return nil, fmt.Errorf("no embedded data found")
	}
	bodySize := int(binary.BigEndian.Uint32(header[4:]))
	if bodySize > s.sourceImage.Bounds().Dx()*s.sourceImage.Bounds().Dy()-headerSize {
		return nil, fmt.Errorf("invalid body size: %d", bodySize)
	}

	body := make([]byte, bodySize)
	for i, pixel := range IteratePixels(s.sourceImage) {
		if i >= headerSize+bodySize {
			break
		}
		if i < headerSize {
			continue
		}
		body[i-headerSize] = s.decodePixel(pixel)
	}

	return body, nil
}

func (s *Steganography) decodePixel(c color.Color) byte {
	r, g, b, _ := c.RGBA()
	return (byte(r) << 6) | ((byte(g) << 3) & 0b00111000) | (byte(b) & 0b00000111)
}

func (s *Steganography) Encode(data []byte) EncodeResult {
	imageSize := s.sourceImage.Bounds().Dx() * s.sourceImage.Bounds().Dy()
	if len(data)+headerSize > imageSize {
		return EncodeResult{
			image: nil,
			err:   errors.New("data is too large to encode"),
		}
	}
	s.targetImage = image.NewNRGBA(s.sourceImage.Bounds())
	draw.Draw(s.targetImage, s.targetImage.Bounds(), s.sourceImage, s.sourceImage.Bounds().Min, draw.Src)

	header := append([]byte(magicBytes), make([]byte, 4)...)
	binary.BigEndian.PutUint32(header[4:], uint32(len(data)))
	encodedData := append(header, data...)

	for i, pixel := range IteratePixels(s.sourceImage) {
		x := i % s.targetImage.Bounds().Dx()
		y := i / s.targetImage.Bounds().Dx()

		if i >= len(encodedData) {
			break
		}
		s.targetImage.Set(x, y, s.encodePixel(pixel, encodedData[i]))
	}

	return EncodeResult{image: s.targetImage, err: nil}
}

func (s *Steganography) encodePixel(c color.Color, data byte) color.NRGBA {
	r, g, b, a := c.RGBA()

	embeddedR := uint8(r&^0b11) | data>>6
	embeddedG := uint8(g&^0b111) | (data>>3)&0b111
	embeddedB := uint8(b&^0b111) | data&0b111

	return color.NRGBA{R: embeddedR, G: embeddedG, B: embeddedB, A: uint8(a)}
}

type EncodeResult struct {
	image *image.NRGBA
	err   error
}

func (e EncodeResult) SaveToFile(fileName string) error {
	if e.err != nil {
		return e.err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	return encoder.Encode(file, e.image)
}

func (e EncodeResult) ToReader() (io.Reader, error) {
	if e.err != nil {
		return nil, e.err
	}
	buffer := new(bytes.Buffer)
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	err := encoder.Encode(buffer, e.image)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
