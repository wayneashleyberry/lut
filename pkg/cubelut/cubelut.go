package cubelut

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
	"math"
	"strconv"
	"strings"
)

// CubeFile implementation
type CubeFile struct {
	Title string  // TITLE
	Size  float64 // LUT_3D_SIZE
	Table map[int][]float64
}

// Parse will parse an io.Reader and return a CubeFile
func Parse(r io.Reader) (CubeFile, error) {
	o := CubeFile{}

	table := map[int][]float64{}

	i := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "LUT_3D_SIZE") {
			s := strings.ReplaceAll(line, "LUT_3D_SIZE ", "")
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return o, err
			}

			o.Size = n
		}

		if strings.HasPrefix(line, "TITLE") {
			s := strings.ReplaceAll(line, "TITLE", "")
			s = strings.ReplaceAll(s, `"`, "")
			o.Title = strings.TrimSpace(s)
		}

		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			continue
		}

		r, err := strconv.ParseFloat(parts[0], 16)
		if err != nil {
			return o, err
		}

		g, err := strconv.ParseFloat(parts[1], 16)
		if err != nil {
			return o, err
		}

		b, err := strconv.ParseFloat(parts[2], 16)
		if err != nil {
			return o, err
		}

		table[i] = []float64{r, g, b}
		i++
	}

	if err := scanner.Err(); err != nil {
		return o, err
	}

	if o.Size == 0 {
		return o, errors.New("invalid lut size")
	}

	o.Table = table

	return o, nil
}

// Apply implementation
func Apply(src image.Image, cube CubeFile, intensity float64) (image.Image, error) {
	if intensity < 0 || intensity > 1 {
		return src, errors.New("intensity must be between 0 and 1")
	}

	bounds := src.Bounds()

	out := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	space := &image.NRGBA{}
	model := space.ColorModel()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := src.At(x, y)
			c := model.Convert(px).(color.NRGBA)

			r := math.Floor((float64(c.R) / 255.0) * (cube.Size - 1))
			g := math.Floor((float64(c.G) / 255.0) * (cube.Size - 1))
			b := math.Floor((float64(c.B) / 255.0) * (cube.Size - 1))

			i := r + cube.Size*g + cube.Size*cube.Size*b

			row := cube.Table[int(i)]

			lr, lg, lb := uint8(row[0]*255), uint8(row[1]*255), uint8(row[2]*255)

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
