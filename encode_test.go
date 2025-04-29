package gostegano

import (
	"image/color"
	"testing"
)

func TestEncodePixel(t *testing.T) {
	tests := []struct {
		name     string
		original color.NRGBA
		data     byte
		want     color.NRGBA
	}{
		{
			name:     "Embed 0 into black",
			original: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
			data:     0b00000000,
			want:     color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		},
		{
			name:     "Embed 0xFF into white",
			original: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
			data:     0b11111111,
			want: color.NRGBA{
				R: (0xFF &^ 0b11) | 0b11,
				G: (0xFF &^ 0b111) | 0b111,
				B: (0xFF &^ 0b111) | 0b111,
				A: 0xFF,
			},
		},
		{
			name:     "Embed pattern into gray",
			original: color.NRGBA{R: 0xAA, G: 0xBB, B: 0xCC, A: 0xFF},
			data:     0b10101010,
			want: color.NRGBA{
				R: (0xAA &^ 0b11) | 0b10,
				G: (0xBB &^ 0b111) | 0b101,
				B: (0xCC &^ 0b111) | 0b010,
				A: 0xFF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodePixel(tt.original, tt.data)
			if got != tt.want {
				t.Errorf(
					"encodePixel(%v, %08b):\n  got  R:%08b G:%08b B:%08b\n  want R:%08b G:%08b B:%08b",
					tt.original, tt.data,
					got.R, got.G, got.B,
					tt.want.R, tt.want.G, tt.want.B,
				)
			}
		})
	}
}
