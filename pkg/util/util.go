// Package util contains... utilities, don't @ me.
package util

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strconv"
	"strings"
)

// Exit will shut down the process with a simple error message
// and the correct error code
func Exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// ReadImage will try and read any supported image type from a filename.
func ReadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".jpg":
		fallthrough
	case ".jpeg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	default:
		return nil, errors.New("unsupported output type: " + filename)
	}
}

// WriteImage will try to write an image.Image to the given file name.
func WriteImage(filename string, img image.Image) error {
	ext := path.Ext(filename)
	switch ext {
	case ".jpg":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		return jpeg.Encode(f, img, &jpeg.Options{
			Quality: 100,
		})
	case ".png":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		return png.Encode(f, img)
	default:
		return errors.New("unsupported output type")
	}
}

// ParseFloats will parse space delimted floats from a string
func ParseFloats(in string, bitSize int) []float64 {
	s := strings.TrimSpace(in)
	parts := strings.Split(s, " ")
	o := []float64{}
	for _, part := range parts {
		f, err := strconv.ParseFloat(part, bitSize)
		if err == nil {
			o = append(o, f)
		}
	}
	return o
}
