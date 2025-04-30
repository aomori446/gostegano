package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"gostegano"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	decodeMode     bool
	encodeMode     bool
	source         string
	message        string
	target         string
	removeAfterUse bool
)

func init() {
	flag.StringVar(&source, "from", "", "Specify URL or file path.")
	flag.StringVar(&target, "to", "", "Target output file name.")

	flag.BoolVar(&decodeMode, "de", false, "Decode a PNG image.")
	flag.BoolVar(&encodeMode, "en", false, "Encode a message into an image.")

	flag.StringVar(&message, "msg", "", "Message to encode.")
	flag.BoolVar(&removeAfterUse, "rm", false, "Delete the file after usage.")

	flag.Parse()
}

func main() {
	if decodeMode == encodeMode {
		flag.Usage()
		return
	}

	reader, err := FromSource(source)
	if err != nil {
		fmt.Printf("Failed to read source: %v\n", err)
		return
	}

	ste, err := gostegano.NewSteganoFrom(reader)
	if err != nil {
		fmt.Printf("Failed to create Stegano from source: %v\n", err)
		return
	}

	switch {
	case decodeMode:
		result, err := ste.Decode()
		if err != nil {
			fmt.Printf("Failed to decode image: %v\n", err)
			return
		}

		fmt.Printf("Decoded message: %s\n", result.ToString())
	case encodeMode:
		if message == "" {
			fmt.Println("Message to encode cannot be empty.")
			return
		}

		result, err := ste.Encode([]byte(message))
		if err != nil {
			fmt.Printf("Failed to encode image: %v\n", err)
			return
		}

		if err = result.SaveFile(target); err != nil {
			fmt.Printf("Failed to save %s : %v\n", target, err)
			return
		}

		fmt.Printf("<%s> encoded and saved to %s\n", message, target)
	}

	if removeAfterUse {
		err = RemoveSource(source)
		if err != nil {
			fmt.Printf("Failed to remove source file: %v\n", err)
			return
		}

		if IsSupportedImageFile(source) && !IsValidImageURL(source) {
			fmt.Printf("Source file %s removed.\n", source)
		}
	}
}

func FromSource(source string) (io.Reader, error) {
	switch {
	case IsSupportedImageFile(source) && !IsValidImageURL(source):
		f, err := os.Open(source)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		data, _ := io.ReadAll(f)

		return bytes.NewReader(data), nil
	case IsValidImageURL(source):
		resp, err := http.Get(source)
		if err != nil || resp.StatusCode != 200 {
			return nil, err
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)

		return bytes.NewReader(data), nil
	default:
		return nil, errors.New("not a supported source format")
	}
}

func IsSupportedImageFile(filename string) bool {
	extensions := []string{".jpg", ".jpeg", ".png"}
	for _, ext := range extensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func IsValidImageURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func RemoveSource(source string) error {
	if IsSupportedImageFile(source) && !IsValidImageURL(source) {
		return os.Remove(source)
	}
	return nil
}
