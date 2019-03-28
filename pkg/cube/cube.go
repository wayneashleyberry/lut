package cube

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
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

	for i, line := range strings.Split(file, "\n") {
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
	}

	space := &image.NRGBA{}
	model := space.ColorModel()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := src.At(x, y)
			c := model.Convert(px).(color.NRGBA)

			fmt.Printf("%+v\n", c)
			os.Exit(1)
		}
	}

	return out, nil
}
