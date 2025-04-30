# gostegano

A lightweight steganography tool written in Go for embedding and extracting text messages in images.

## Features

- ✅ Embed messages into images (steganography)
- ✅ Extract embedded messages from images
- ✅ Load images from local files or URLs
- ✅ Optionally delete the source image after use
- ✅ **Supports PNG, JPG, and JPEG formats**

## Installation

```bash
go install github.com/aomori446/gostegano/cmd/gostegano@latest
```
OR

```bash
git clone https://github.com/aomori446/gostegano.git
cd gostegano/cmd/gostegano
go build -o gostegano
```

## Usage

### Embed a message into an image

```bash
./gostegano -en -from input.png -msg "Secret message" -to output.png
```

- `-en`: Enable encode mode  
- `-from`: Source image (local file or URL)  
- `-msg`: Message to embed  
- `-to`: Output image file  

> ✅ **Supported input/output formats**: PNG, JPG, JPEG

### Extract a message from an image

```bash
./gostegano -de -from output.png
```

- `-de`: Enable decode mode  
- `-from`: Image file with an embedded message  

### Optional: Remove the source file after processing

```bash
./gostegano -en -from input.jpg -msg "Secret message" -to output.jpg -rm
./gostegano -de -from output.jpg -rm
```

## License

[MIT LICENSE](https://github.com/aomori446/gostegano/blob/main/LICENSE)

