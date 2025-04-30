package gostegano

import (
	"image/color"
	"io"
	"os"
)

func decodePixel(encodedPixel color.Color) (decodeData byte) {
	switch c := encodedPixel.(type) {
	case color.RGBA:
		return (c.R << 6) | ((c.G &^ 0b11111000) << 3) | (c.B &^ 0b11111000)
	case color.NRGBA:
		return (c.R << 6) | ((c.G &^ 0b11111000) << 3) | (c.B &^ 0b11111000)
	case color.RGBA64:
		return (uint8(c.R) << 6) | ((uint8(c.G) &^ 0b11111000) << 3) | (uint8(c.B) &^ 0b11111000)
	case color.NRGBA64:
		return (uint8(c.R) << 6) | ((uint8(c.G) &^ 0b11111000) << 3) | (uint8(c.B) &^ 0b11111000)
	default:
		panic("not a support format.")
	}
}

type DecodeResult struct {
	decodedData []byte
	header      []byte
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
