package gostegano

import (
	"image/color"
	"io"
	"os"
)

func decodePixel(c color.Color) byte {
	r, g, b, _ := c.RGBA()
	return (byte(r) << 6) | ((byte(g) << 3) & 0b00111000) | (byte(b) & 0b00000111)
}

type DecodeResult struct {
	decodedData []byte
	readOffset  int
}

func (result *DecodeResult) Read(p []byte) (n int, err error) {
	if result.readOffset >= len(result.decodedData) {
		return 0, io.EOF
	}

	n = copy(p, result.decodedData[result.readOffset:])
	result.readOffset += n
	return n, nil
}

func (result *DecodeResult) WriteTo(w io.Writer) (n int64, err error) {
	written, err := w.Write(result.decodedData)
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

func (result *DecodeResult) ToString() string {
	return string(result.decodedData)
}
