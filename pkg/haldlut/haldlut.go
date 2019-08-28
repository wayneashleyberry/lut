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

	for z := 0; z < cube.Size; z++ {
		for x := 0; x < cube.Size; x++ {
			for y := 0; y < cube.Size; y++ {
				imgx := (z % 8 * cube.Size) + x
				imgy := (z / 8 * cube.Size) + y
				rgb := cube.Get(x, y, z)
				out.SetNRGBA(imgx, imgy, color.NRGBA{
					R: uint8(rgb[0] * 0xff),
					G: uint8(rgb[1] * 0xff),
					B: uint8(rgb[2] * 0xff),
					A: 0xff,
				})
			}
		}
	}

	return out
}

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
	nrgba := space.ColorModel()

	space2 := &image.RGBA{}
	rgba := space2.ColorModel()

	width, height := bounds.Dx(), bounds.Dy()
	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				// not all images use the same colour space, so ensure we convert
				// them all to nrgba to be consistent with our output
				px := src.At(x, y)
				c := nrgba.Convert(px).(color.NRGBA)

				// Android reference implementation:
				// float x = clamp(g % 8 * 64 + r, 0.0, 511.0);
				// float y = clamp(floor(g/8.0f) + b * 8, 0.0, 511.0);
				// uchar4 lutPoint = rsGetElementAt_uchar4(lut, x, y);

				// find the location of the pixel in our lookup table
				lutx := clamp(int(c.G%8*64+c.R), 0, 511)
				luty := clamp(int(math.Floor(float64(c.G)/8)+float64(c.B)*8), 0, 511)

				pixel := effect.At(lutx, luty)
				lut := rgba.Convert(pixel).(color.RGBA)

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

func clamp(val, min, max int) int {
	if val > max {
		return max
	}
	if val < min {
		return min
	}
	return val
}
