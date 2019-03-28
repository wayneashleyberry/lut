package cube

import (
	"image"
	"image/color"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/sgreben/piecewiselinear"
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

	i := 0

	xaxis := []float64{}
	raxis := []float64{}
	gaxis := []float64{}
	baxis := []float64{}

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

	for n := 0; n < i; n++ {
		row := table[n]
		xaxis = append(xaxis, float64(n)/float64(len(table)))
		raxis = append(raxis, float64(row[0]))
		gaxis = append(gaxis, float64(row[1]))
		baxis = append(baxis, float64(row[2]))
	}

	fr := piecewiselinear.Function{
		X: xaxis,
		Y: raxis,
	}

	fg := piecewiselinear.Function{
		X: xaxis,
		Y: gaxis,
	}

	fb := piecewiselinear.Function{
		X: xaxis,
		Y: baxis,
	}

	space := &image.NRGBA{}
	model := space.ColorModel()

	// n := float64(32)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := src.At(x, y)
			c := model.Convert(px).(color.NRGBA)

			r := float64(float64(c.R) / 255.0)
			g := float64(float64(c.G) / 255.0)
			b := float64(float64(c.B) / 255.0)

			estr := fr.At(r)
			estg := fg.At(g)
			estb := fb.At(b)

			o := color.NRGBA{
				R: uint8(estr * 255),
				G: uint8(estg * 255),
				B: uint8(estb * 255),
				A: 255,
			}

			out.Set(x, y, o)
		}
	}

	return out, nil
}
