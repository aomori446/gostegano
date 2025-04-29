package gostegano

import (
	"image"
	"image/color"
	"io"
	"os"
)

func decode(sourceImage image.Image) (result DecodeResult, err error) {
	bodySize, err := getBodySize(sourceImage)
	if err != nil {
		return
	}

	body := make([]byte, bodySize)
	for index, pixel := range iteratePixel(sourceImage) {
		if index >= headerSize+bodySize {
			break
		}
		if index < headerSize {
			continue
		}
		body[index-headerSize] = decodePixel(pixel)
	}

	return DecodeResult{
		DecodedData: body,
	}, nil
}

func decodePixel(c color.Color) byte {
	r, g, b, _ := c.RGBA()
	return (byte(r) << 6) | ((byte(g) << 3) & 0b00111000) | (byte(b) & 0b00000111)
}

type DecodeResult struct {
	DecodedData []byte
	readOffset  int
}

func (result *DecodeResult) Read(p []byte) (n int, err error) {
	if result.readOffset >= len(result.DecodedData) {
		return 0, io.EOF
	}

	n = copy(p, result.DecodedData[result.readOffset:])
	result.readOffset += n
	return n, nil
}

func (result *DecodeResult) WriteTo(w io.Writer) (n int64, err error) {
	written, err := w.Write(result.DecodedData)
	if err != nil {
		return 0, err
	}
	return int64(written), nil
}

func (result *DecodeResult) SaveFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = result.WriteTo(f)
	return err
}
