package haldlut

import (
	"errors"
	"image"
	"image/color"
	"math"

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

	space := &image.NRGBA{}
	model := space.ColorModel()

	width, height := bounds.Dx(), bounds.Dy()
	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				// not all images use the same colour space, so ensure we convert
				// them all to nrgba to be consistent with our output
				px := src.At(x, y)
				c := model.Convert(px).(color.NRGBA)

				lutx := (c.G/4%8)*64 + (c.R / 4)
				luty := math.Floor(float64(c.G/4)/8) + float64(c.B/4)*8

				l := effect.At(int(lutx), int(luty))
				lut := model.Convert(l).(color.NRGBA)

				// create our output colour, adjusted according to the intensity
				o := color.NRGBA{}
				o.R = uint8(float64(c.R)*(1-intensity) + float64(lut.R)*intensity)
				o.G = uint8(float64(c.G)*(1-intensity) + float64(lut.G)*intensity)
				o.B = uint8(float64(c.B)*(1-intensity) + float64(lut.B)*intensity)
				o.A = c.A

				out.Set(x, y, o)
			}
		}
	})

	return out, nil
}
