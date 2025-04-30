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

/*
Process:

"A.png".openFile().Decode() => originalImage

originalImage.copy() => newImage

newImage.At(x,y) => color.encodeColor(*data) => encodedColor

newImage.Set(x,y,encodedColor).Encode().saveFile() => "A_copy.png"

"A_copy.png".openFile().Decode() => newImage.At(x,y) => color.decodeColor() => *data

*/

const (
	magicBytes = "GOST"
	headerSize = 4 + len(magicBytes)
)

type Stegano struct {
	image image.Image
}

func NewStegano(sourceImage image.Image) *Stegano {
	return &Stegano{image: sourceImage}
}

func NewSteganoFrom(reader io.Reader) (*Stegano, error) {
	sourceImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to Decode image: %w", err)
	}
	return &Stegano{image: sourceImage}, nil
}

func (s *Stegano) Decode() (result *DecodeResult, err error) {
	bodySize, err := getBodySize(s.image)
	if err != nil {
		return
	}

	body := make([]byte, bodySize)
	for index, pixel := range iteratePixel(s.image, headerSize, headerSize+bodySize) {
		body[index] = decodePixel(pixel)
	}

	return &DecodeResult{
		decodedData: body,
	}, nil
}

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
