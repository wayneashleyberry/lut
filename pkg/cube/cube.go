package cube

import (
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/wayneashleyberry/lut/pkg/util"
)

// Apply implementation
func Apply(srcfile, lutfile string) (image.Image, error) {
	src, err := util.ReadImage(srcfile)
	if err != nil {
		return nil, err
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

	table := map[int][]float32{}

	// size := float32(32)
	// dmin := []float32{0, 0, 0}
	// dmax := []float32{1, 1, 1}

	i := 0

	for _, line := range strings.Split(file, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
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

		table[i] = []float32{
			float32(r),
			float32(g),
			float32(b),
		}

		i++
	}

	space := &image.NRGBA{}
	model := space.ColorModel()

	n := float32(32)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := src.At(x, y)
			c := model.Convert(px).(color.NRGBA)

			r := float32(float32(c.R)/255.0) * (n - 1)
			g := float32(float32(c.G)/255.0) * (n - 1)
			b := float32(float32(c.B)/255.0) * (n - 1)

			idx := int(math.Floor(float64(r + n*g + n*n*b)))

			lutr := table[idx][0]
			lutg := table[idx][1]
			lutb := table[idx][2]

			o := color.NRGBA{
				R: uint8(lutr * 255),
				G: uint8(lutg * 255),
				B: uint8(lutb * 255),
				A: 255,
			}

			out.Set(x, y, o)
		}
	}

	return out, nil
}
