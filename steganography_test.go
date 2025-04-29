package gostegano

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
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
		err:         nil,
		targetImage: wantImage,
	}

	// Assertions
	if !isSameImage(result.image, wantResult.targetImage) {
		t.Errorf("Image output does not match expected targetImage")
	}
	if result.Err != nil {
		t.Errorf("Unexpected error occurred: %v", result.Err)
	}
	if wantResult.err != nil {
		t.Errorf("Unexpected error in expected result: %v", wantResult.err)
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

	f, err := os.CreateTemp("", "test_*.png")
	if err != nil {
		t.Errorf("Create Temp png file failed :%v", err)
	}

	if err = png.Encode(f, sourceImage); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	f.Close()

	f1, err := os.Open(f.Name())
	if err != nil {
		t.Errorf("open file %s failed", f1.Name())
	}
	defer os.Remove(f1.Name())

	steg, err := NewSteganographyFromReader(f1)
	if err != nil {
		t.Errorf("create steganography failed")
	}
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
		t.Errorf("Decoded data mismatched. Got: %s, Want: %s", string(decoded), string(originalData))
	}
}

func TestDecodePixel(t *testing.T) {
	tests := []struct {
		name     string
		input    color.NRGBA
		expected byte
	}{
		{
			name:     "All zeros",
			input:    color.NRGBA{R: 0b00000000, G: 0b00000000, B: 0b00000000, A: 255},
			expected: 0b00000000,
		},
		{
			name:     "Max value",
			input:    color.NRGBA{R: 0b00000011, G: 0b00000111, B: 0b00000111, A: 255},
			expected: 0b11111111,
		},
		{
			name:     "Random middle value",
			input:    color.NRGBA{R: 0b00000010, G: 0b00000101, B: 0b00000011, A: 255},
			expected: 0b10101011,
		},
		{
			name:     "Another pattern",
			input:    color.NRGBA{R: 0b00000001, G: 0b00000011, B: 0b00000010, A: 255},
			expected: 0b01011010,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodePixel(tt.input)
			if got != tt.expected {
				t.Errorf("decodePixel() = %08b, want %08b", got, tt.expected)
			}
		})
	}
}
