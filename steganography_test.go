package gostegano

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestEmbedByteInPixel(t *testing.T) {
	tests := []struct {
		name  string
		color color.Color
		data  byte
		wantR byte
		wantG byte
		wantB byte
	}{
		{
			name:  "All zero, data 0xFF",
			color: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			data:  0xFF,
			wantR: 0b11,
			wantG: 0b111,
			wantB: 0b111,
		},
		{
			name:  "All zero, data 0x00",
			color: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			data:  0x00,
			wantR: 0b00,
			wantG: 0b000,
			wantB: 0b000,
		},
		{
			name:  "Mixed type RGBA64, data 0xAA",
			color: color.RGBA64{R: 0, G: 0, B: 0, A: 0},
			data:  0xAA,
			wantR: 0b10,
			wantG: 0b101,
			wantB: 0b010,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EmbedByteInPixel(tt.color, tt.data)
			if got.R != tt.wantR || got.G != tt.wantG || got.B != tt.wantB {
				t.Errorf("got RGB = %08b %08b %08b, want %08b %08b %08b",
					got.R, got.G, got.B, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestEmbedDataInImage(t *testing.T) {
	makeBaseImage := func() *image.RGBA {
		img := image.NewRGBA(image.Rect(0, 0, 256, 256))
		draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{A: 255}}, image.Point{}, draw.Src)
		return img
	}

	tests := []struct {
		name    string
		input   image.Image
		data    []byte
		want    image.Image
		wantErr bool
	}{
		{
			name:  "basic embedding",
			input: makeBaseImage(),
			data:  []byte{0xAA, 0xBB, 0xCC},
			want: func() image.Image {
				base := makeBaseImage()
				want := image.NewNRGBA(base.Bounds())
				draw.Draw(want, want.Bounds(), base, image.Point{}, draw.Src)
				padding := []byte{0x00, 0x00, 0x00, 0x03}
				payload := append(padding, 0xAA, 0xBB, 0xCC)
				for i, b := range payload {
					x := i % want.Bounds().Dx()
					y := i / want.Bounds().Dx()
					want.Set(x, y, EmbedByteInPixel(want.At(x, y), b))
				}
				return want
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EmbedDataInImage(tt.input, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmbedDataInImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !imagesEqual(got, tt.want) {
				t.Error("EmbedDataInImage() image mismatch")
			}
		})
	}
}

func imagesEqual(img1, img2 image.Image) bool {
	if !img1.Bounds().Eq(img2.Bounds()) {
		return false
	}
	for y := img1.Bounds().Min.Y; y < img1.Bounds().Max.Y; y++ {
		for x := img1.Bounds().Min.X; x < img1.Bounds().Max.X; x++ {
			if img1.At(x, y) != img2.At(x, y) {
				return false
			}
		}
	}
	return true
}
