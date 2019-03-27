package lut

import (
	"errors"
	"image"
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

	// intensity := uint32(amount)
	// fmt.Println(amount, intensity, 0xFFFFFFFF, intensity == 0xFFFFFFFF)

	bounds := src.Bounds()

	out := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			rr, gg, bb := uint8(r/4>>8), uint8(g/4>>8), uint8(b/4>>8)

			lutx := (bb%8)*64 + rr
			luty := math.Floor(float64(bb)/8)*64 + float64(gg)

			lut := effect.At(int(lutx), int(luty))
			lr, lg, lb, _ := lut.RGBA()

			o := col32{}
			o.R = uint32(float64(r)*(1-intensity) + float64(lr)*intensity)
			o.G = uint32(float64(g)*(1-intensity) + float64(lg)*intensity)
			o.B = uint32(float64(b)*(1-intensity) + float64(lb)*intensity)
			o.A = a

			out.Set(x, y, o)
		}
	}

	return out, nil
}
