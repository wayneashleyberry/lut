// Package trilinear implements trilinear interpolation
package trilinear

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/wayneashleyberry/lut/pkg/colorcube"
	"github.com/wayneashleyberry/lut/pkg/parallel"
)

// bits per channel (we're assuming 8 bits).
const bpc = 0xff

// Interpolate will apply color transformations to the provided image using
// trilinear interpolation (taking the intensity multiplier into account).
func Interpolate(src image.Image, cube colorcube.Cube, intensity float64) (image.Image, error) {
	if intensity < 0 || intensity > 1 {
		return src, errors.New("intensity must be between 0 and 1")
	}

	bounds := src.Bounds()

	out := image.NewNRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{bounds.Max.X, bounds.Max.Y},
	})

	k := (float64(cube.Size) - 1) / bpc

	space := &image.NRGBA{}
	model := space.ColorModel()

	dKR := cube.DomainMax[0] - cube.DomainMin[0]
	dKG := cube.DomainMax[1] - cube.DomainMin[1]
	dKB := cube.DomainMax[2] - cube.DomainMin[2]

	width, height := bounds.Dx(), bounds.Dy()
	parallel.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				px := src.At(x, y)
				c := model.Convert(px).(color.NRGBA)

				rgb := getFromRGBTrilinear(
					int(c.R),
					int(c.G),
					int(c.B),
					cube.Size,
					k,
					cube,
				)

				o := color.NRGBA{}
				o.R = uint8(float64(c.R)*(1-intensity) + float64(toIntCh(rgb[0]*dKR))*intensity)
				o.G = uint8(float64(c.G)*(1-intensity) + float64(toIntCh(rgb[1]*dKG))*intensity)
				o.B = uint8(float64(c.B)*(1-intensity) + float64(toIntCh(rgb[2]*dKB))*intensity)
				o.A = c.A

				out.Set(x, y, o)
			}
		}
	})

	return out, nil
}

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
	switch {
	case x <= 0:
		return 0
	case x >= bpc:
		return bpc
	default:
		return x
	}
}

func toIntCh(x float64) int {
	return clampToChannelSize(int(math.Floor(x * float64(bpc))))
}

func getFromRGBTrilinear(r, g, b, size int, k float64, cube colorcube.Cube) []float64 {
	iR := float64(r) * k

	var fR1 int
	if iR >= float64(size)-1 {
		fR1 = clampToChannelSize(size - 1)
	} else {
		fR1 = clampToChannelSize(int(math.Floor(iR + 1)))
	}

	var fR0 int
	if iR > 0 {
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
	if iG > 0 {
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
	if iB > 0 {
		fB0 = clampToChannelSize(int(math.Floor(iB - 1)))
	}

	c000 := cube.Get(fR0, fG0, fB0)
	c010 := cube.Get(fR0, fG1, fB0)
	c001 := cube.Get(fR0, fG0, fB1)
	c011 := cube.Get(fR0, fG1, fB1)
	c101 := cube.Get(fR1, fG0, fB1)
	c100 := cube.Get(fR1, fG0, fB0)
	c110 := cube.Get(fR1, fG1, fB0)
	c111 := cube.Get(fR1, fG1, fB1)

	rx := trilerp(
		iR, iG, iB, c000[0], c001[0], c010[0], c011[0],
		c100[0], c101[0], c110[0], c111[0],
		float64(fR0), float64(fR1), float64(fG0), float64(fG1), float64(fB0), float64(fB1),
	)

	gx := trilerp(
		iR, iG, iB, c000[1], c001[1], c010[1], c011[1],
		c100[1], c101[1], c110[1], c111[1],
		float64(fR0), float64(fR1), float64(fG0), float64(fG1), float64(fB0), float64(fB1),
	)

	bx := trilerp(
		iR, iG, iB, c000[2], c001[2], c010[2], c011[2],
		c100[2], c101[2], c110[2], c111[2],
		float64(fR0), float64(fR1), float64(fG0), float64(fG1), float64(fB0), float64(fB1),
	)

	return []float64{rx, gx, bx}
}
