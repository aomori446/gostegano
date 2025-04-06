package main

import (
	"flag"
	"github.com/aomori446/gostegano"
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
	flag.BoolVar(&removeAfterUse, "rm", false, "Delete the file after usage.")

	flag.BoolVar(&decodeMode, "decode", false, "Decode a PNG image.")
	flag.BoolVar(&encodeMode, "encode", false, "Encode a message into an image.")

	flag.StringVar(&message, "message", "", "Message to encode.")
	flag.StringVar(&target, "target", "", "Target output file name.")

	flag.Parse()
}

func main() {
	if decodeMode == encodeMode {
		flag.Usage()
		return
	}

	if !gostegano.IsSupportedImageFile(source) && !gostegano.IsValidImageURL(source) {
		flag.Usage()
		return
	}

}
