package gostegano

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"io"
)

/*
Header:
   +--------+-----------+
   | "GOST" | len(body) |
   +--------+-----------+
   |   4    |     4     |
   +--------+-----------+
*/

const (
	magicBytes = "GOST"
	headerSize = 4 + len(magicBytes)
)

type Stegano struct {
	sourceImage image.Image
	targetImage *image.NRGBA // Only used in Encode mode
}

func NewStegano(sourceImage image.Image) *Stegano {
	return &Stegano{sourceImage: sourceImage}
}

func NewSteganoFrom(reader io.Reader) (*Stegano, error) {
	sourceImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to Decode targetImage: %w", err)
	}
	return &Stegano{sourceImage: sourceImage}, nil
}

func (s *Stegano) Decode() (result *DecodeResult, err error) {
	bodySize, err := getBodySize(s.sourceImage)
	if err != nil {
		return
	}

	body := make([]byte, bodySize)
	for index, pixel := range iteratePixel(s.sourceImage) {
		if index >= headerSize+bodySize {
			break
		}
		if index < headerSize {
			continue
		}
		body[index-headerSize] = decodePixel(pixel)
	}

	return &DecodeResult{
		decodedData: body,
	}, nil
}

func (s *Stegano) Encode(body []byte) (*EncodeResult, error) {
	if err := validatePayloadSize(s.sourceImage, body); err != nil {
		return nil, err
	}

	targetImage := copyImage(s.sourceImage)

	header := newHeader(len(body))
	payload := append(header, body...)

	for index, pixel := range iteratePixel(s.sourceImage) {
		x := index % s.sourceImage.Bounds().Dx()
		y := index / s.sourceImage.Bounds().Dx()

		if index >= len(payload) {
			break
		}
		targetImage.Set(x, y, encodePixel(pixel, payload[index]))
	}

	return &EncodeResult{targetImage: targetImage}, nil
}
