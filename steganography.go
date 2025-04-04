package gostegano

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"iter"
	"net/http"
	"os"
	"strings"
)

const paddingSize = 4

// ピクセルにバイトデータを埋め込む
func embedByteInPixel(c color.Color, data byte) color.NRGBA {
	r, g, b, a := c.RGBA()

	embeddedR := uint8(r&^0b11) | data>>6
	embeddedG := uint8(g&^0b111) | (data>>3)&0b111
	embeddedB := uint8(b&^0b111) | data&0b111

	return color.NRGBA{R: embeddedR, G: embeddedG, B: embeddedB, A: uint8(a)}
}

// EmbedDataInImage 画像にデータを埋め込む
func EmbedDataInImage(img image.Image, data []byte) (image.Image, error) {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	if !CanFitDataInImage(len(data), width*height) {
		return nil, errors.New("データサイズが画像の容量を超えています")
	}

	newImg := image.NewNRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Src)

	// データサイズを先頭4バイトに格納する
	padding := make([]byte, paddingSize)
	binary.BigEndian.PutUint32(padding, uint32(len(data)))
	data = append(padding, data...)

	for i, v := range data {
		x, y := i%width, i/width%height
		newImg.Set(x, y, embedByteInPixel(newImg.At(x, y), v))
	}

	return newImg, nil
}

// SaveEncodedImage 埋め込まれたデータをPNGファイルとして保存する
func SaveEncodedImage(r io.Reader, data []byte, fileName string) error {
	if fileName == "" || !strings.HasSuffix(fileName, ".png") {
		return errors.New("出力ファイルは .png 形式でなければなりません")
	}

	reader := ToBufferedReader(r)
	img, _, err := image.Decode(reader)
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return errors.New("対応していない画像フォーマットです")
		}
		return err
	}

	newImg, err := EmbedDataInImage(img, data)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	return encoder.Encode(f, newImg)
}

// ピクセルからバイトデータを抽出する
func extractByteFromPixel(c color.Color) byte {
	r, g, b, _ := c.RGBA()
	return (byte(r) << 6) | ((byte(g) << 3) & 0b00111000) | (byte(b) & 0b00000111)
}

// ExtractDataFromImage 画像からデータを抽出する
func ExtractDataFromImage(img image.Image) []byte {
	dataSize := GetEmbeddedDataSize(img)

	decodedData := make([]byte, paddingSize+dataSize)
	for i, v := range IterateImagePixels(img) {
		if i >= len(decodedData) {
			break
		}
		decodedData[i] = extractByteFromPixel(v)
	}

	return decodedData[paddingSize:]
}

// ReadAndExtractData Reader からデータを抽出する
func ReadAndExtractData(r io.Reader) ([]byte, error) {
	reader := ToBufferedReader(r)
	img, err := png.Decode(reader)
	if err != nil {
		return nil, err
	}
	return ExtractDataFromImage(img), nil
}

// CanFitDataInImage 画像にデータを埋め込めるか確認する
func CanFitDataInImage(length, totalPixels int) bool {
	return length > 0 && length <= totalPixels
}

// IterateImagePixels 画像のすべてのピクセルを走査する
func IterateImagePixels(img image.Image) iter.Seq2[int, color.Color] {
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

// GetEmbeddedDataSize 埋め込まれたデータのサイズを取得する
func GetEmbeddedDataSize(img image.Image) int {
	padding := make([]byte, paddingSize)
	for i, v := range IterateImagePixels(img) {
		if i >= paddingSize {
			break
		}
		padding[i] = extractByteFromPixel(v)
	}
	return int(binary.BigEndian.Uint32(padding))
}

// OpenImageSource ファイルまたはURLを開く
func OpenImageSource(s string) (io.ReadCloser, error) {
	switch {
	case IsValidImageURL(s):
		resp, err := http.Get(s)
		if err != nil || resp.StatusCode != 200 {
			return nil, fmt.Errorf("%s を取得できません", s)
		}
		return resp.Body, nil
	case IsSupportedImageFile(s):
		return os.Open(s)
	default:
		return nil, errors.New("無効なファイルまたはURLです")
	}
}

// ToBufferedReader io.Reader を bufio.Reader に変換する
func ToBufferedReader(r io.Reader) *bufio.Reader {
	if reader, ok := r.(*bufio.Reader); ok {
		return reader
	}
	return bufio.NewReader(r)
}

// IsSupportedImageFile 画像ファイル形式を判定する
func IsSupportedImageFile(s string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range extensions {
		if strings.HasSuffix(s, ext) {
			return true
		}
	}
	return false
}

// IsValidImageURL URL かどうかを判定する
func IsValidImageURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
