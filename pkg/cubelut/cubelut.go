package cubelut

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/overhq/lut/pkg/colorcube"
	"github.com/overhq/lut/pkg/util"
)

// bits per channel (we're assuming 8-bits)
const bpc = 0xff

// CubeFile implementation
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

// Parse will parse an io.Reader and return a CubeFile
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

// ColorCube will create a new color cube from a lut file
func (lut CubeFile) ColorCube() colorcube.Cube {
	cube := colorcube.New(lut.Size, lut.DomainMin, lut.DomainMax)

	for i := 0; i < lut.Size*lut.Size*lut.Size; i++ {
		x := i % lut.Size
		y := i / lut.Size % lut.Size
		z := i / lut.Size / lut.Size
		cube.Set(x, y, z, []float64{lut.R[i], lut.G[i], lut.B[i]})
	}

	return cube
}
