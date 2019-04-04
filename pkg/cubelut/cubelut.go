package cubelut

import (
	"bufio"
	"errors"
	"image"
	"io"
	"strconv"
	"strings"

	"github.com/wayneashleyberry/lut/pkg/util"
)

// CubeFile implementation
type CubeFile struct {
	Dimensions int
	DomainMax  []float64 // DOMAIN_MAX
	DomainMin  []float64 // DOMAIN_MIN
	Size       float64   // LUT_3D_SIZE
	Table      map[int][]float64
	Title      string // TITLE
}

// Parse will parse an io.Reader and return a CubeFile
func Parse(r io.Reader) (CubeFile, error) {
	o := CubeFile{}

	// Defaults
	o.Dimensions = 1
	o.DomainMin = []float64{0, 0, 0}
	o.DomainMax = []float64{1, 1, 1}

	table := map[int][]float64{}

	i := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "TITLE") {
			s := strings.ReplaceAll(line, "TITLE", "")
			s = strings.ReplaceAll(s, `"`, "")
			o.Title = strings.TrimSpace(s)
			continue
		}

		if strings.HasPrefix(line, "DOMAIN_MIN") {
			s := strings.ReplaceAll(line, "DOMAIN_MIN", "")
			min := util.ParseFloats(s, 8)
			if len(min) != 3 {
				return o, errors.New("invalid domain min values")
			}
			o.DomainMin = min
			continue
		}

		if strings.HasPrefix(line, "DOMAIN_MAX") {
			s := strings.ReplaceAll(line, "DOMAIN_MAX", "")
			max := util.ParseFloats(s, 8)
			if len(max) != 3 {
				return o, errors.New("invalid domain max values")
			}
			o.DomainMax = max
			continue
		}

		if strings.HasPrefix(line, "LUT_3D_SIZE") {
			s := strings.ReplaceAll(line, "LUT_3D_SIZE ", "")
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return o, err
			}

			o.Size = n
			o.Dimensions = 3
			continue
		}

		rgb := util.ParseFloats(line, 16)
		if len(rgb) == 3 {
			table[i] = rgb
			i++
			continue
		}
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

	// space := &image.NRGBA{}
	// model := space.ColorModel()

	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
	// 		px := src.At(x, y)
	// 		c := model.Convert(px).(color.NRGBA)

	// 		r := math.Floor((float64(c.R) / 0xff) * (cube.Size - 1))
	// 		g := math.Floor((float64(c.G) / 0xff) * (cube.Size - 1))
	// 		b := math.Floor((float64(c.B) / 0xff) * (cube.Size - 1))

	// 		i := r + cube.Size*g + cube.Size*cube.Size*b

	// 		row := cube.Table[int(i)]

	// 		lr, lg, lb := row[0]*0xff, row[1]*0xff, row[2]*0xff

	// 		o := color.NRGBA{}
	// 		o.R = uint8(float64(c.R)*(1-intensity) + lr*intensity)
	// 		o.G = uint8(float64(c.G)*(1-intensity) + lg*intensity)
	// 		o.B = uint8(float64(c.B)*(1-intensity) + lb*intensity)
	// 		o.A = c.A

	// 		out.Set(x, y, o)
	// 	}
	// }

	return out, nil
}
