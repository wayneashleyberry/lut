package lut

import (
	"errors"
	"image"
	"image/color"
	"math"
)

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

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// not all images use the same colour space, so ensure we convert
			// them all to nrgba to be consistent with our output
			px := src.At(x, y)
			c := space.ColorModel().Convert(px).(color.NRGBA)

			// find the location of the pixel in our lookup table
			lutx := int((c.B/4%8)*64 + c.R/4)
			luty := int(math.Floor(float64(c.B/4)/8)*64 + float64(c.G/4))

			lut := effect.At(lutx, luty).(color.RGBA)

			// create our output colour, adjusted according to the intensity
			o := color.NRGBA{}
			o.R = uint8(float64(c.R)*(1-intensity) + float64(lut.R)*intensity)
			o.G = uint8(float64(c.G)*(1-intensity) + float64(lut.G)*intensity)
			o.B = uint8(float64(c.B)*(1-intensity) + float64(lut.B)*intensity)
			o.A = c.A

			out.Set(x, y, o)
		}
	}

	return out, nil
}
