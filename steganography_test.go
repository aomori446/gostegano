package gostegano

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	// Step 1: Create a blank image (100x100)
	width, height := 100, 100
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Fill image with white
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}

	// Step 2: Encode message into image
	message := "Hello, Stegano!"
	steg := NewStegano(img)

	encoded, err := steg.Encode([]byte(message))
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	// Step 3: Convert encoded image to PNG and read back (simulate saving/loading)
	var buf bytes.Buffer
	if err = png.Encode(&buf, encoded.image); err != nil {
		t.Fatalf("Failed to encode image to PNG: %v", err)
	}

	decodedSteg, err := NewSteganoFrom(&buf)
	if err != nil {
		t.Fatalf("Failed to decode PNG back to image: %v", err)
	}

	// Step 4: Decode message
	result, err := decodedSteg.Decode()
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	got := result.ToString()
	if got != message {
		t.Errorf("Decoded message mismatch.\nGot:  %q\nWant: %q", got, message)
	}
}
