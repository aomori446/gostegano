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

type Steganography struct {
	sourceImage image.Image
	targetImage *image.NRGBA // Only used in Encode mode
}

func NewSteganography(sourceImage image.Image) *Steganography {
	return &Steganography{sourceImage: sourceImage}
}

func NewSteganographyFromReader(reader io.Reader) (*Steganography, error) {
	sourceImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode targetImage: %w", err)
	}
	return &Steganography{sourceImage: sourceImage}, nil
}
