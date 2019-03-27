package lut

import (
	"errors"
	"image"
	"image/color"
	"math"
)

type col32 struct {
	R, G, B, A uint32
}

func (c col32) RGBA() (uint32, uint32, uint32, uint32) {
	return c.R, c.G, c.B, c.A
}

// Apply implementation
func Apply(src, effect image.Image, intensity float64) (image.Image, error) {
	if intensity < 0 || intensity > 1 {
		return src, errors.New("intensity must be between 0 and 1")
	}

	bounds := src.Bounds()

	out := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.At(x, y).(color.NRGBA)

			lutx := int((c.B/4%8)*64 + c.R/4)
			luty := int(math.Floor(float64(c.B/4)/8)*64 + float64(c.G/4))

			lut := effect.At(lutx, luty).(color.RGBA)

			o := color.NRGBA{}
			o.R = lut.R
			o.G = lut.G
			o.B = lut.B
			o.A = c.A

			out.Set(x, y, o)

			// if c.A != 0 {
			// 	fmt.Printf(
			// 		"x: %d y: %d col: %+v lut x: %d lut y: %d out: %+v\n",
			// 		x, y, c, lutx, luty, o,
			// 	)
			// }
		}
	}

	return out, nil
}
