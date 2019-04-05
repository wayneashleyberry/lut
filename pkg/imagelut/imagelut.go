package imagelut

import (
	"errors"
	"image"
	"image/color"

	"github.com/overhq/lut/pkg/colorcube"
)

// Parse implementation
func Parse(src image.Image) (colorcube.Cube, error) {
	// hardcoded defaults
	size := 64
	dmin := []float64{0, 0, 0}
	dmax := []float64{1, 1, 1}

	cube := colorcube.New(size, dmin, dmax)

	bounds := src.Bounds()
	if bounds.Max.X != 512 || bounds.Max.Y != 512 {
		return cube, errors.New("invalid image size")
	}

	space := &image.NRGBA{}
	model := space.ColorModel()

	for z := 0; z < size; z++ {
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				imgx := (z % 8 * 64) + x
				imgy := (z / 8 * 64) + y
				px := src.At(imgx, imgy)
				c := model.Convert(px).(color.NRGBA)

				cube.Set(x, y, z, []float64{
					float64(c.R) / 0xff,
					float64(c.G) / 0xff,
					float64(c.B) / 0xff,
				})
			}
		}
	}

	return cube, nil
}
