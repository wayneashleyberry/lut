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

	"github.com/wayneashleyberry/lut/pkg/util"
)

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

// Apply implementation
func Apply(src image.Image, lut CubeFile, intensity float64) (image.Image, error) {
	if intensity < 0 || intensity > 1 {
		return src, errors.New("intensity must be between 0 and 1")
	}

	bounds := src.Bounds()

	out := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	// Build a color cube

	var cube Table3D

	for i := 0; i < lut.Size*lut.Size*lut.Size; i++ {
		x := i % lut.Size
		y := i / lut.Size % lut.Size
		z := i / lut.Size / lut.Size
		cube[x][y][z] = []float64{lut.R[i], lut.G[i], lut.B[i]}
	}

	// fmt.Println(cube[0][0][0])
	// fmt.Println(cube[31][0][0])
	// fmt.Println(cube[0][31][0])
	// fmt.Println(cube[0][0][31])

	k := (float64(lut.Size) - 1) / bpc

	space := &image.NRGBA{}
	model := space.ColorModel()

	dKR := lut.DomainMax[0] - lut.DomainMin[0]
	dKG := lut.DomainMax[1] - lut.DomainMin[1]
	dKB := lut.DomainMax[2] - lut.DomainMin[2]

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := src.At(x, y)
			c := model.Convert(px).(color.NRGBA)

			rgb := getFromRGBTrilinear(
				int(c.R),
				int(c.G),
				int(c.B),
				lut.Size,
				k,
				cube,
			)

			o := color.NRGBA{}
			o.R = uint8(toIntCh(rgb[0] * dKR))
			o.G = uint8(toIntCh(rgb[1] * dKG))
			o.B = uint8(toIntCh(rgb[2] * dKB))
			o.A = c.A

			// fmt.Println(c, o)

			out.Set(x, y, o)
		}
	}

	return out, nil
}

// int _toIntCh(double x) => _clampToChannelSize((x * bpc).floor());
// func toIntCh(x float64) int {
// return clampToChannelSize(int(math.Floor((x * float64(bpc)))))
// }

// Table3D Implementation
type Table3D [32][32][32][]float64

func trilerp(x, y, z, c000, c001, c010, c011, c100, c101, c110, c111, x0, x1, y0, y1, z0, z1 float64) float64 {
	xd := (x - x0) / (x1 - x0)
	yd := (y - y0) / (y1 - y0)
	zd := (z - z0) / (z1 - z0)

	c00 := c000*(1.0-xd) + c100*xd
	c01 := c001*(1.0-xd) + c101*xd
	c10 := c010*(1.0-xd) + c110*xd
	c11 := c011*(1.0-xd) + c111*xd

	c0 := c00*(1.0-yd) + c10*yd
	c1 := c01*(1.0-yd) + c11*yd

	c := c0*(1.0-zd) + c1*zd

	return c
}

func clampToChannelSize(x int) int {
	if x < 0 {
		return 0
	}
	if x > bpc {
		return bpc
	}
	return x
}

func toIntCh(x float64) int {
	return clampToChannelSize(int(math.Floor(x * float64(bpc))))
}

func getFromRGBTrilinear(r, g, b, size int, k float64, table3D Table3D) []float64 {
	iR := float64(r) * k

	var fR1 int
	if iR >= float64(size)-1 {
		fR1 = clampToChannelSize(size - 1)
	} else {
		fR1 = clampToChannelSize(int(math.Floor(iR + 1)))
	}

	var fR0 int
	if iR <= 0 {
		fR0 = 0
	} else {
		fR0 = clampToChannelSize(int(math.Floor(iR - 1)))
	}

	iG := float64(g) * k

	var fG1 int
	if iG >= float64(size)-1 {
		fG1 = clampToChannelSize(size - 1)
	} else {
		fG1 = clampToChannelSize(int(math.Floor(iG + 1)))
	}

	var fG0 int
	if iG <= 0 {
		fG0 = 0
	} else {
		fG0 = clampToChannelSize(int(math.Floor(iG - 1)))
	}

	iB := float64(b) * k

	var fB1 int
	if iB >= float64(size)-1 {
		fB1 = clampToChannelSize(size - 1)
	} else {
		fB1 = clampToChannelSize(int(math.Floor(iB + 1)))
	}

	var fB0 int
	if iB <= 0 {
		fB0 = 0
	} else {
		fB0 = clampToChannelSize(int(math.Floor(iB - 1)))
	}

	// fmt.Println(fR0, fR1, fG0, fG1, fB0, fB1)

	c000 := table3D[fR0][fG0][fB0]
	c010 := table3D[fR0][fG1][fB0]
	c001 := table3D[fR0][fG0][fB1]
	c011 := table3D[fR0][fG1][fB1]
	c101 := table3D[fR1][fG0][fB1]
	c100 := table3D[fR1][fG0][fB0]
	c110 := table3D[fR1][fG1][fB0]
	c111 := table3D[fR1][fG1][fB1]

	rx := trilerp(
		iR, iG, iB, c000[0], c001[0], c010[0], c011[0],
		c100[0], c101[0], c110[0], c111[0], float64(fR0), float64(fR1), float64(fG0), float64(fG1), float64(fB0), float64(fB1),
	)

	gx := trilerp(
		iR, iG, iB, c000[1], c001[1], c010[1], c011[1],
		c100[1], c101[1], c110[1], c111[1], float64(fR0), float64(fR1), float64(fG0), float64(fG1), float64(fB0), float64(fB1),
	)

	bx := trilerp(
		iR, iG, iB, c000[1], c001[1], c010[1], c011[1],
		c100[1], c101[1], c110[1], c111[1], float64(fR0), float64(fR1), float64(fG0), float64(fG1), float64(fB0), float64(fB1),
	)

	return []float64{rx, gx, bx}
}
