package gostegano

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"iter"
	"slices"
)

func validateBodySize(img image.Image, bodySize int) error {
	if bodySize+headerSize > img.Bounds().Dx()*img.Bounds().Dy() {
		return fmt.Errorf("invalid body size: %d", bodySize)
	}
	return nil
}

func getBodySize(sourceImage image.Image) (int, error) {
	header := make([]byte, headerSize)
	for index, pixel := range iteratePixel(sourceImage, 0, headerSize) {
		header[index] = decodePixel(pixel)
	}

	if !slices.Equal(header[:4], []byte(magicBytes)) {
		return 0, fmt.Errorf("no embedded data found")
	}

	bodySize := int(binary.BigEndian.Uint32(header[4:]))
	if err := validateBodySize(sourceImage, bodySize); err != nil {
		return 0, err
	}

	return bodySize, nil
}

func copyImage(sourceImage image.Image) (targetImage *image.NRGBA) {
	targetImage = image.NewNRGBA(sourceImage.Bounds())
	draw.Draw(targetImage, targetImage.Bounds(), sourceImage, sourceImage.Bounds().Min, draw.Src)
	return
}

func loadPayload(body []byte) (payload []byte) {
	header := append([]byte(magicBytes), make([]byte, 4)...)
	binary.BigEndian.PutUint32(header[4:], uint32(len(body)))
	return append(header, body...)
}

func iteratePixel(img image.Image, start, end int) iter.Seq2[int, color.Color] {
	return func(yield func(int, color.Color) bool) {
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()

		if start < 0 || start > end || end > width*height {
			return
		}

		index := 0
		for i := start; i < end; i++ {
			x := i % width
			y := i / width
			if y >= height {
				break
			}
			if !yield(index, img.At(x, y)) {
				return
			}
			index++
		}
	}
}
