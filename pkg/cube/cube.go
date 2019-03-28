package cube

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
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

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := src.At(x, y)
			c := model.Convert(px).(color.NRGBA)

			r := float32(float32(c.R) / 255.0)
			g := float32(float32(c.G) / 255.0)
			b := float32(float32(c.B) / 255.0)

			ri := int(math.Floor(float64(r) * float64(len(table)-1)))
			gi := int(math.Floor(float64(g) * float64(len(table)-1)))
			bi := int(math.Floor(float64(b) * float64(len(table)-1)))

			lutr := table[ri][0]
			lutg := table[gi][1]
			lutb := table[bi][2]

			o := color.NRGBA{
				R: uint8(lutr * 255),
				G: uint8(lutg * 255),
				B: uint8(lutb * 255),
				A: 255,
			}

			fmt.Println(r, g, b)
			fmt.Println(lutr, lutg, lutb)
			os.Exit(1)

			out.Set(x, y, o)
		}
	}

	return out, nil
}
