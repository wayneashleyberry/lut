package cubelut

import (
	"errors"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type col64 struct {
	R, G, B float64
}

// Apply implementation
func Apply(src image.Image, lutfile string, intensity float64) (image.Image, error) {
	if intensity < 0 || intensity > 1 {
		return src, errors.New("intensity must be between 0 and 1")
	}

	bounds := src.Bounds()

	out := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	b, err := ioutil.ReadFile(lutfile)
	if err != nil {
		return out, err
	}

	file := string(b)

	table := map[int]col64{}

	i := 0

	var n float64 // LUT_3D_SIZE

	for _, line := range strings.Split(file, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "LUT_3D_SIZE") {
			s := strings.ReplaceAll(line, "LUT_3D_SIZE ", "")
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return out, err
			}

			n = f
		}

		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			continue
		}

		r, err := strconv.ParseFloat(parts[0], 32)
		if err != nil {
			return out, err
		}

		g, err := strconv.ParseFloat(parts[1], 32)
		if err != nil {
			return out, err
		}

		b, err := strconv.ParseFloat(parts[2], 32)
		if err != nil {
			return out, err
		}

		table[i] = col64{R: r, G: g, B: b}
		i++
	}

	if n == 0 {
		return src, errors.New("invalid lut size")
	}

	space := &image.NRGBA{}
	model := space.ColorModel()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := src.At(x, y)
			c := model.Convert(px).(color.NRGBA)

			r := math.Floor((float64(c.R) / 255.0) * (n - 1))
			g := math.Floor((float64(c.G) / 255.0) * (n - 1))
			b := math.Floor((float64(c.B) / 255.0) * (n - 1))

			i := r + n*g + n*n*b

			l := table[int(i)]

			lr, lg, lb := uint8(l.R*255), uint8(l.G*255), uint8(l.B*255)

			o := color.NRGBA{}
			o.R = uint8(float64(c.R)*(1-intensity) + float64(lr)*intensity)
			o.G = uint8(float64(c.G)*(1-intensity) + float64(lg)*intensity)
			o.B = uint8(float64(c.B)*(1-intensity) + float64(lb)*intensity)
			o.A = c.A

			out.Set(x, y, o)
		}
	}

	return out, nil
}
