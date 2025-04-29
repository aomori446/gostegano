package gostegano

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"iter"
	"strings"
)

func validatePayloadSize(img image.Image, body []byte) error {
	imageSize := img.Bounds().Dx() * img.Bounds().Dy()
	if len(body)+headerSize > imageSize {
		return errors.New("body is too large to Encode into targetImage")
	}
	return nil
}

func validateBodySize(img image.Image, bodySize int) error {
	if bodySize > img.Bounds().Dx()*img.Bounds().Dy()-headerSize {
		return fmt.Errorf("invalid body size: %d", bodySize)
	}
	return nil
}

func getBodySize(sourceImage image.Image) (int, error) {
	header := make([]byte, headerSize)
	for i, pixel := range iteratePixelN(sourceImage, headerSize) {
		header[i] = decodePixel(pixel)
	}

	if string(header[:4]) != magicBytes {
		return -1, fmt.Errorf("no embedded data found")
	}

	bodySize := int(binary.BigEndian.Uint32(header[4:]))
	if err := validateBodySize(sourceImage, bodySize); err != nil {
		return -1, err
	}

	return bodySize, nil
}

func copyImage(sourceImage image.Image) (targetImage *image.NRGBA) {
	targetImage = image.NewNRGBA(sourceImage.Bounds())
	draw.Draw(targetImage, targetImage.Bounds(), sourceImage, sourceImage.Bounds().Min, draw.Src)
	return
}

func newHeader(bodySize int) (header []byte) {
	header = append([]byte(magicBytes), make([]byte, 4)...)
	binary.BigEndian.PutUint32(header[4:], uint32(bodySize))
	return
}

func iteratePixel(img image.Image) iter.Seq2[int, color.Color] {
	return func(yield func(int, color.Color) bool) {
		bounds := img.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				index := x + y*bounds.Dx()
				if !yield(index, img.At(x, y)) {
					return
				}
			}
		}
	}
}

func iteratePixelN(img image.Image, n int) iter.Seq2[int, color.Color] {
	return func(yield func(int, color.Color) bool) {
		bounds := img.Bounds()

		if n > bounds.Dx()*bounds.Dy() {
			n = bounds.Dx() * bounds.Dy()
		}

		count := 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				if count >= n {
					return
				}
				index := x + y*bounds.Dx()
				if !yield(index, img.At(x, y)) {
					return
				}
				count++
			}
		}
	}
}

func isSameImage(img1, img2 image.Image) bool {
	if img1.Bounds() != img2.Bounds() {
		return false
	}

	bounds := img1.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img1.At(x, y) != img2.At(x, y) {
				return false
			}
		}
	}

	return true
}

func IsSupportedImageFile(s string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range extensions {
		if strings.HasSuffix(s, ext) {
			return true
		}
	}
	return false
}

func IsValidImageURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
