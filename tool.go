package gostegano

import (
	"image"
	"image/color"
	"iter"
	"strings"
)

func IteratePixels(img image.Image) iter.Seq2[int, color.Color] {
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

func IsSameImage(img1, img2 image.Image) bool {
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
