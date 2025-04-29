package gostegano

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
)

func encodePixel(c color.Color, data byte) color.NRGBA {
	r, g, b, a := c.RGBA()

	embeddedR := uint8(r&^0b11) | data>>6
	embeddedG := uint8(g&^0b111) | (data>>3)&0b111
	embeddedB := uint8(b&^0b111) | data&0b111

	return color.NRGBA{R: embeddedR, G: embeddedG, B: embeddedB, A: uint8(a)}
}

type EncodeResult struct {
	targetImage *image.NRGBA
	buffer      *bytes.Reader
}

func (result *EncodeResult) Read(p []byte) (n int, err error) {
	if result.buffer == nil {
		if err = result.prepareReader(); err != nil {
			return 0, err
		}
	}
	return result.buffer.Read(p)
}

func (result *EncodeResult) prepareReader() error {
	var buf bytes.Buffer
	err := png.Encode(&buf, result.targetImage)
	if err != nil {
		return err
	}
	result.buffer = bytes.NewReader(buf.Bytes())
	return nil
}

func (result *EncodeResult) WriteTo(w io.Writer) (n int64, err error) {
	countingWriter := &byteCounter{Writer: w}
	err = png.Encode(countingWriter, result.targetImage)
	if err != nil {
		return 0, err
	}
	return countingWriter.n, nil
}

func (result *EncodeResult) SaveFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = result.WriteTo(f)
	return err
}

type byteCounter struct {
	io.Writer
	n int64
}

func (b *byteCounter) Write(p []byte) (int, error) {
	n, err := b.Writer.Write(p)
	b.n += int64(n)
	return n, err
}
