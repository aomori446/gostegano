// Package gostegano provides lightweight steganography encoding/decoding
// using image files (PNG, JPG, JPEG). It allows embedding and extracting
// text messages from images by manipulating pixel color bits.
package gostegano

import (
	"fmt"
	"image"
	"io"
)

const (
	// magicBytes defines the file signature "GOST" used as a header prefix.
	magicBytes = "GOST"

	// headerSize is the total size of the header (magicBytes + body length).
	headerSize = 4 + len(magicBytes)
)

// Stegano represents a steganography wrapper for an image.
type Stegano struct {
	image image.Image
}

// NewStegano creates a Stegano instance from a given image.
func NewStegano(sourceImage image.Image) *Stegano {
	return &Stegano{image: sourceImage}
}

// NewSteganoFrom decodes an image from an io.Reader and creates a Stegano instance.
func NewSteganoFrom(reader io.Reader) (*Stegano, error) {
	sourceImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return &Stegano{image: sourceImage}, nil
}

// Encode embeds the given body into the image by manipulating
// the RGB bits of each pixel. It returns a new image with the hidden message.
//
// The encoding process includes:
//  1. Prepending a header (magic + length)
//  2. Iterating over the image pixels
//  3. Modifying RGB bits to embed each byte of payload
func (s *Stegano) Encode(body []byte) (*EncodeResult, error) {
	if err := validateBodySize(s.image, len(body)); err != nil {
		return nil, err
	}

	targetImage := copyImage(s.image)
	payload := loadPayload(body)

	for index, pixel := range iteratePixel(targetImage, 0, len(payload)) {
		x := index % targetImage.Bounds().Dx()
		y := index / targetImage.Bounds().Dx()
		data := payload[index]
		targetImage.Set(x, y, encodePixel(pixel, data))
	}

	return &EncodeResult{image: targetImage}, nil
}

// Decode extracts the hidden data embedded in an image by reading the pixels
// and decoding the RGB bits that store the message.
//
// The decoding process includes:
//  1. Extracting the header to determine body size
//  2. Iterating pixels to reconstruct the byte stream
func (s *Stegano) Decode() (*DecodeResult, error) {
	bodySize, err := getBodySize(s.image)
	if err != nil {
		return nil, err
	}

	body := make([]byte, bodySize)
	for index, pixel := range iteratePixel(s.image, headerSize, headerSize+bodySize) {
		body[index] = decodePixel(pixel)
	}

	return &DecodeResult{
		decodedData: body,
	}, nil
}
