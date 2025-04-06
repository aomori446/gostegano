package gostegano

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestSteganography_Encode(t *testing.T) {
	testData := []byte("happy")
	sourceImage := image.NewRGBA(image.Rect(0, 0, 32, 32))
	draw.Draw(sourceImage, sourceImage.Bounds(), &image.Uniform{C: color.NRGBA{A: 255}}, image.Point{}, draw.Src)

	steganography := NewSteganography(sourceImage)
	result := steganography.Encode(testData)

	wantImage := image.NewNRGBA(sourceImage.Bounds())
	draw.Draw(wantImage, wantImage.Bounds(), sourceImage, sourceImage.Bounds().Min, draw.Src)

	// Construct expected pixels: magic bytes + size header + body
	expectedData := []byte{
		'G', 'O', 'S', 'T',
		0x00, 0x00, 0x00, 0x05,
		'h', 'a', 'p', 'p', 'y',
	}
	for i, b := range expectedData {
		SetPixelByIndex(wantImage, i, encodePixel(color.NRGBA{A: 255}, b))
	}

	wantResult := EncodeResult{
		Err:   nil,
		image: wantImage,
	}

	// Assertions
	if !IsSameImage(result.image, wantResult.image) {
		t.Errorf("Image output does not match expected image")
	}
	if result.Err != nil {
		t.Errorf("Unexpected error occurred: %v", result.Err)
	}
	if wantResult.Err != nil {
		t.Errorf("Unexpected error in expected result: %v", wantResult.Err)
	}
}

func SetPixelByIndex(img *image.NRGBA, index int, col color.NRGBA) {
	bounds := img.Bounds()
	width := bounds.Dx()
	if index < 0 || index >= width*bounds.Dy() {
		panic("index out of bounds")
	}
	x := index % width
	y := index / width
	img.Set(x, y, col)
}

func TestSteganography_Decode(t *testing.T) {
	// Given
	originalData := []byte("test message")
	width, height := 16, 16
	sourceImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(sourceImage, sourceImage.Bounds(), &image.Uniform{C: color.NRGBA{A: 255}}, image.Point{}, draw.Src)

	// Encode the data into the image
	steg := NewSteganography(sourceImage)
	result := steg.Encode(originalData)
	if result.Err != nil {
		t.Fatalf("Encode failed: %v", result.Err)
	}

	// When
	steg2 := NewSteganography(result.image)
	decoded, err := steg2.Decode()

	// Then
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if string(decoded) != string(originalData) {
		t.Errorf("Decoded data mismatch. Got: %s, Want: %s", string(decoded), string(originalData))
	}
}
