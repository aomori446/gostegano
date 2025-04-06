package gostegano

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"iter"
	"net/http"
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

func FetchFrom(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("can not fetch picture from %s\n", url)
	}

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bts), nil
}
