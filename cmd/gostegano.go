package main

import (
	"flag"
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
	flag.StringVar(&source, "source", "", "Specify the source URL or file path.")
	flag.StringVar(&target, "target", "", "Target output file name.")

	flag.BoolVar(&decodeMode, "decode", false, "Decode a PNG image.")
	flag.BoolVar(&encodeMode, "encode", false, "Encode a message into an image.")

	flag.StringVar(&message, "message", "", "Message to encode.")
	flag.BoolVar(&removeAfterUse, "rm", false, "Delete the file after usage.")

	flag.Parse()
}

/*func main() {
	if decodeMode == encodeMode {
		flag.Usage()
		return
	}

	if !gostegano.IsSupportedImageFile(source) && !gostegano.IsValidImageURL(source) {
		flag.Usage()
		return
	}

	var reader io.Reader
	switch {
	case gostegano.IsSupportedImageFile(source):
		f, err := os.Open(source)
		if err != nil {
			fmt.Printf("Failed to open file: %s\n", source)
			return
		}
		defer f.Close()

		bts, err := io.ReadAll(f)
		if err != nil {
			fmt.Printf("Failed to read file: %s\n", source)
			return
		}
		reader = bytes.NewReader(bts)
	case gostegano.IsValidImageURL(source):
		resp, err := http.Get(source)
		if err != nil || resp.StatusCode != 200 {
			fmt.Printf("Failed to fetch image from URL: %s\n", source)
			return
		}
		defer resp.Body.Close()

		bts, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body from URL: %s\n", source)
			return
		}
		reader = bytes.NewReader(bts)
	}

	stegano, err := gostegano.NewSteganoFromReader(reader)
	if err != nil {
		fmt.Printf("Failed to create Stegano from source: %v\n", err)
		return
	}

	switch {
	case decodeMode:
		data, err := gostegano
		if err != nil {
			fmt.Printf("Failed to decode image: %v\n", err)
			return
		}
		fmt.Printf("Decoded message: %s\n", string(data))
	case encodeMode:
		if message == "" {
			fmt.Println("Message to encode cannot be empty.")
			return
		}

		if err = stegano.Encode([]byte(message)).SaveToFile(target); err != nil {
			fmt.Printf("Failed to save encoded image: %v\n", err)
			return
		}
		fmt.Printf("Message encoded and saved to %s\n", target)
	}

	if removeAfterUse {
		err = os.Remove(source)
		if err != nil {
			fmt.Printf("Failed to remove source file: %v\n", err)
		} else {
			fmt.Printf("Source file %s removed.\n", source)
		}
	}
}*/
