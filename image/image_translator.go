package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
)

type InvalidFormatError struct {
}

func (e *InvalidFormatError) Error() string {
	return fmt.Sprintf("Invalid image format")
}

func bufferStartsWith(buffer []byte, needle []byte) bool {
	if len(buffer) < len(needle) {
		return false
	}
	for i := range needle {
		if buffer[i] != needle[i] {
			return false
		}
	}
	return true
}

func isJpegFile(buffer []byte) bool {
	return bufferStartsWith(buffer, []byte{0xff, 0xd8, 0xff, 0xe0})
}

func isPngFile(buffer []byte) bool {
	return bufferStartsWith(buffer, []byte{0x89, 0x50, 0x4e, 0x47})
}

func isGifFile(buffer []byte) bool {
	return bufferStartsWith(buffer, []byte{0x47, 0x49, 0x46, 0x38})
}

func isBmpFile(buffer []byte) bool {
	return bufferStartsWith(buffer, []byte{0x42, 0x4d})
}

func TranslatePngToJpeg(buffer []byte) ([]byte, error) {
	decode, err := png.Decode(bytes.NewReader(buffer))
	if err != nil {
		return []byte{}, nil
	}

	newImg := image.NewRGBA(decode.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	draw.Draw(newImg, newImg.Bounds(), decode, decode.Bounds().Min, draw.Over)

	buff := new(bytes.Buffer)
	writer := io.Writer(buff)

	var opt jpeg.Options
	opt.Quality = 100

	err = jpeg.Encode(writer, newImg, &opt)
	if err != nil {
		return []byte{}, nil
	}

	return buff.Bytes(), nil
}

func TranslateToReadableImageForAll(buffer []byte) ([]byte, error) {
	if isJpegFile(buffer) {
		return buffer, nil
	}
	if isPngFile(buffer) {
		toJpeg, err := TranslatePngToJpeg(buffer)
		return toJpeg, err
	}
	if isGifFile(buffer) {
		return buffer, nil
	}
	return []byte{}, &InvalidFormatError{}
}
