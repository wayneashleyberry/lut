package cubelut

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/wayneashleyberry/lut/pkg/colorcube"
	"github.com/wayneashleyberry/lut/pkg/parallel"
	"github.com/wayneashleyberry/lut/pkg/util"
)

// CubeFile implementation.
type CubeFile struct {
	Dimensions int
	DomainMax  []float64 // DOMAIN_MAX
	DomainMin  []float64 // DOMAIN_MIN
	Size       int       // LUT_3D_SIZE
	Title      string    // TITLE
	R          []float64
	G          []float64
	B          []float64
}

// FromColorCube will create a cube file from a color cube.
func FromColorCube(cube colorcube.Cube) CubeFile {
	r := make([]float64, cube.Size*cube.Size*cube.Size)
	g := make([]float64, cube.Size*cube.Size*cube.Size)
	b := make([]float64, cube.Size*cube.Size*cube.Size)

	for i := 0; i < cube.Size*cube.Size*cube.Size; i++ {
		x := i % cube.Size
		y := i / cube.Size % cube.Size
		z := i / cube.Size / cube.Size
		rgb := cube.Get(x, y, z)
		r[i] = rgb[0]
		g[i] = rgb[1]
		b[i] = rgb[2]
	}

	return CubeFile{
		Dimensions: 3,
		DomainMax:  cube.DomainMax,
		DomainMin:  cube.DomainMin,
		Size:       cube.Size,
		R:          r,
		G:          g,
		B:          b,
	}
}

// Parse will parse an io.Reader and return a CubeFile.
func Parse(r io.Reader) (CubeFile, error) {
	o := CubeFile{}

	// Defaults
	o.Dimensions = 1
	o.DomainMin = []float64{0, 0, 0}
	o.DomainMax = []float64{1, 1, 1}

	i := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

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

			n, err := strconv.ParseInt(s, 0, 64)
			if err != nil {
				return o, err
			}

			o.Size = int(n)
			o.Dimensions = 3
			o.R = make([]float64, n*n*n)
			o.G = make([]float64, n*n*n)
			o.B = make([]float64, n*n*n)

			continue
		}

		rgb := util.ParseFloats(line, 16)
		if len(rgb) == 3 {
			o.R[i] = rgb[0]
			o.G[i] = rgb[1]
			o.B[i] = rgb[2]
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

	return o, nil
}

// Cube will convert a cube file into a color cube.
func (cf CubeFile) Cube() colorcube.Cube {
	cube := colorcube.New(cf.Size, cf.DomainMin, cf.DomainMax)

	for i := 0; i < cf.Size*cf.Size*cf.Size; i++ {
		x := i % cf.Size
		y := i / cf.Size % cf.Size
		z := i / cf.Size / cf.Size
		cube.Set(x, y, z, []float64{cf.R[i], cf.G[i], cf.B[i]})
	}

	return cube
}

// Bytes implementation.
func (cf CubeFile) Bytes() []byte {
	var b bytes.Buffer

	fmt.Fprintf(&b, `TITLE "%s"
LUT_3D_SIZE %d
DOMAIN_MIN %.1f %.1f %.1f
DOMAIN_MAX %.1f %.1f %.1f
`, cf.Title, cf.Size, cf.DomainMin[0], cf.DomainMin[1], cf.DomainMin[2], cf.DomainMax[0], cf.DomainMax[1], cf.DomainMax[2])

	for i := range cf.R {
		fmt.Fprintf(&b, "%.6f %.6f %.6f\n", cf.R[i], cf.G[i], cf.B[i])
	}

	return b.Bytes()
}

// Apply implementation.
func (cf CubeFile) Apply(src image.Image, intensity float64) (image.Image, error) {
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

	N := float64(cf.Size)

	width, height := bounds.Dx(), bounds.Dy()
	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				px := src.At(x, y)
				c := model.Convert(px).(color.NRGBA)

				r := math.Floor((float64(c.R) / 255.0) * (N - 1))
				g := math.Floor((float64(c.G) / 255.0) * (N - 1))
				b := math.Floor((float64(c.B) / 255.0) * (N - 1))

				i := r + N*g + N*N*b

				lr, lg, lb := uint8(cf.R[int(i)]*255), uint8(cf.G[int(i)]*255), uint8(cf.B[int(i)]*255)

				o := color.NRGBA{}
				o.R = uint8(float64(c.R)*(1-intensity) + float64(lr)*intensity)
				o.G = uint8(float64(c.G)*(1-intensity) + float64(lg)*intensity)
				o.B = uint8(float64(c.B)*(1-intensity) + float64(lb)*intensity)
				o.A = c.A

				out.Set(x, y, o)
			}
		}
	})

	return out, nil
}
