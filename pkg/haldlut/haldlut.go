package haldlut

import (
	"errors"
	"image"
	"math"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/overhq/lut/pkg/colorcube"
	"github.com/overhq/lut/pkg/parallel"
)

// FromColorCube will create an image from a color cube
func FromColorCube(cube colorcube.Cube) image.Image {
	out := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{cube.Size * 8, cube.Size * 8},
	})

	return out
}

// Parse implementation
func Parse(src image.Image) (colorcube.Cube, error) {
	return colorcube.Cube{}, errors.New("not yet implemented")
}

// Apply colour transformations to an image from the provided lookup table
func Apply(src, effect image.Image, intensity float64) (image.Image, error) {
	if intensity < 0 || intensity > 1 {
		return src, errors.New("intensity must be between 0 and 1")
	}

	bounds := src.Bounds()

	out := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	width, height := bounds.Dx(), bounds.Dy()
	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				px := src.At(x, y)

				c, _ := colorful.MakeColor(px)
				r, g, b := c.RGB255()

				lutx := (g%8)*64 + r
				luty := math.Floor(float64(g)/8) + float64(b)*8

				at := effect.At(int(lutx), int(luty))
				o, _ := colorful.MakeColor(at)

				out.Set(x, y, o)
			}
		}
	})

	return out, nil
}
