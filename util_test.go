package gostegano

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestIteratePixel(t *testing.T) {
	tests := []struct {
		name       string
		img        image.Image
		start, end int
		want       map[int]color.Color
	}{
		{
			name:  "[2x2][start:0,end:4][RGBA:0x00]",
			img:   image.NewNRGBA(image.Rect(0, 0, 2, 2)),
			start: 0,
			end:   4,
			want: map[int]color.Color{
				0: color.NRGBA{},
				1: color.NRGBA{},
				2: color.NRGBA{},
				3: color.NRGBA{},
			},
		},
		{
			name: "[2x2][start:0,end:4][RGBA:0xFF]",
			img: func() image.Image {
				img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
				draw.Draw(img, img.Bounds(), &image.Uniform{C: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}}, image.Point{}, draw.Src)
				return img
			}(),
			start: 0,
			end:   4,
			want: map[int]color.Color{
				0: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
				1: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
				2: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
				3: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
			},
		},
		{
			name: "[100x100][start:25,end:30][Black]",
			img: func() image.Image {
				img := image.NewNRGBA(image.Rect(0, 0, 100, 100))
				draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)
				return img
			}(),
			start: 25,
			end:   30,
			want: map[int]color.Color{
				25: color.Black,
				26: color.Black,
				27: color.Black,
				28: color.Black,
				29: color.Black,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for relIndex, pixel := range iteratePixel(test.img, test.start, test.end) {
				actualIndex := test.start + relIndex
				wantPixel, ok := test.want[actualIndex]
				if !ok {
					t.Errorf("unexpected pixel at index %d", actualIndex)
					continue
				}

				r1, g1, b1, a1 := wantPixel.RGBA()
				r2, g2, b2, a2 := pixel.RGBA()

				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Errorf("at index %d: want RGBA(%d,%d,%d,%d), got RGBA(%d,%d,%d,%d)",
						actualIndex, r1, g1, b1, a1, r2, g2, b2, a2)
				}
			}
		})
	}
}
