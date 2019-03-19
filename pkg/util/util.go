// Package util contains... utilities, don't @ me.
package util

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
)

// ReadImage will try and read any supported image type from a filename.
func ReadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := path.Ext(filename)
	switch ext {
	case ".jpg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	default:
		return nil, errors.New("unsupported output type")
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
