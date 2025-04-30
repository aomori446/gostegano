package gostegano

import (
	"image/color"
	"math/rand"
	"testing"
)

func TestDecodePixel(t *testing.T) {
	tests := []struct {
		name string
		c    color.NRGBA
		want byte
	}{
		{
			name: "All 0",
			c: color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0x00,
			},
			want: 0b00000000,
		},
		{
			name: "All 1",
			c: color.NRGBA{
				R: 0xFF,
				G: 0xFF,
				B: 0xFF,
				A: 0xFF,
			},
			want: 0b11111111,
		},
		{
			name: "Mid Values",
			c: color.NRGBA{
				R: 0xAA, // 10101010
				G: 0xBB, // 10111011
				B: 0xCC, // 11001100
				A: 0xDD,
			},
			want: decodePixel(color.NRGBA{R: 0xAA, G: 0xBB, B: 0xCC, A: 0xDD}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := decodePixel(test.c)
			if got != test.want {
				t.Errorf("decodePixel(%v): want %08b, got %08b", test.c, test.want, got)
			}
		})
	}
}

func TestEncodeAndDecodePixel(t *testing.T) {
	for data := range byte(0xFF) {
		originalPixel := color.NRGBA{R: uint8(rand.Intn(0xFF)), G: uint8(rand.Intn(0xFF)), B: uint8(rand.Intn(0xFF)), A: uint8(rand.Intn(0xFF))}
		encodedPixel := encodePixel(originalPixel, data)
		decodedData := decodePixel(encodedPixel)
		if decodedData != data {
			t.Errorf("data = %d; originalPixel = %v; encodedPixel => %d; decodedData => %d\n", data, originalPixel, encodedPixel, decodedData)
		}
	}
}
