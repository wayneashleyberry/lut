package split

import (
	"image"
	"image/color"

	"github.com/wayneashleyberry/lut/pkg/util"
)

// Image will split an image into red, green, blue and alpha images.
func Image(src image.Image) error {
	bounds := src.Bounds()

	red := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	green := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	blue := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	alpha := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			red.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(r >> 8),
				B: uint8(r >> 8),
				A: 255,
			})
			green.Set(x, y, color.RGBA{
				R: uint8(g >> 8),
				G: uint8(g >> 8),
				B: uint8(g >> 8),
				A: 255,
			})
			blue.Set(x, y, color.RGBA{
				R: uint8(b >> 8),
				G: uint8(b >> 8),
				B: uint8(b >> 8),
				A: 255,
			})
			alpha.Set(x, y, color.RGBA{
				R: uint8(a >> 8),
				G: uint8(a >> 8),
				B: uint8(a >> 8),
				A: 255,
			})
		}
	}

	if err := util.WriteImage("red.jpg", red); err != nil {
		return err
	}

	if err := util.WriteImage("green.jpg", green); err != nil {
		return err
	}

	if err := util.WriteImage("blue.jpg", blue); err != nil {
		return err
	}

	if err := util.WriteImage("alpha.jpg", alpha); err != nil {
		return err
	}

	return nil
}
