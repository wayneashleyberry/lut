package lut

import (
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
func Apply(src, effect image.Image) (image.Image, error) {
	bounds := src.Bounds()

	out := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			rr, gg, bb := uint8(r/4>>8), uint8(g/4>>8), uint8(b/4>>8)

			lutx := (bb%8)*64 + rr
			luty := math.Floor(float64(bb)/8)*64 + float64(gg)

			lut := effect.At(int(lutx), int(luty))
			lr, lg, lb, la := lut.RGBA()

			out.Set(x, y, col32{lr, lg, lb, la})
		}
	}

	return out, nil
}
