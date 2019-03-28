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

type col64 struct {
	R, G, B float64
}

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

	table := map[int]col64{}

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

		table[i] = col64{R: r, G: g, B: b}
		i++
	}

	space := &image.NRGBA{}
	model := space.ColorModel()
	N := float64(32) // LUT_3D_SIZE

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := src.At(x, y)
			c := model.Convert(px).(color.NRGBA)

			r := math.Floor((float64(c.R) / 255.0) * (N - 1))
			g := math.Floor((float64(c.G) / 255.0) * (N - 1))
			b := math.Floor((float64(c.B) / 255.0) * (N - 1))

			i := r + N*g + N*N*b

			l := table[int(i)]

			o := color.NRGBA{
				R: uint8(l.R * 255),
				G: uint8(l.G * 255),
				B: uint8(l.B * 255),
				A: 255,
			}

			out.Set(x, y, o)
		}
	}

	return out, nil
}
