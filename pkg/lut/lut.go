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
			// r, g, b, a := c.RGBA()
			// if a == 0 {
			// 	out.Set(x, y, c)
			// 	continue
			// }
			// rr, gg, bb := c.RGB255()

			// rr, gg, bb := c.RGB255()
			// c := src.At(x, y)
			// r, g, b, a := c.RGBA()
			// rr := uint8(float64(r*255.0) + 0.5)
			// gg := uint8(float64(g*255.0) + 0.5)
			// bb := uint8(float64(b*255.0) + 0.5)
			// r, g, b, a := c.R, c.G, c.B, c.A
			// rr := uint8(r / 4 >> 8)
			// gg := uint8(g / 4 >> 8)
			// bb := uint8(b / 4 >> 8)

			lutx := int((c.B%8)*64 + c.R)
			luty := int(math.Floor(float64(c.B)/8)*64 + float64(c.G))

			lut := effect.At(lutx, luty).(color.RGBA)
			// lr, lg, lb, _ := lut.RGBA()

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

			// if a == 0 {
			// 	out.Set(x, y, color.RGBA{0, 0, 0, 0})
			// 	continue
			// }

			// if a == 0xffff {
			// 	out.Set(x, y, color.NRGBA{
			// 		R: uint8(lr >> 8),
			// 		G: uint8(lg >> 8),
			// 		B: uint8(lb >> 8),
			// 		A: 0xff,
			// 	})
			// 	continue
			// }

			// xr := (lr * 0xffff) / a
			// xg := (lg * 0xffff) / a
			// xb := (lb * 0xffff) / a

			// out.Set(x, y, color.NRGBA{
			// 	R: uint8(xr >> 8),
			// 	G: uint8(xg >> 8),
			// 	B: uint8(xb >> 8),
			// 	A: uint8(a >> 8),
			// })
		}
	}

	return out, nil
}
