// Package colorcube implements a very simple data structure for a 3d color cube
// and contains helper methods for getting and settings values at specific
// points. The underlying data structure is a multidimensional map which needs
// to be initialised to a specific size to keep things efficient.
package colorcube

import (
	"image/color"
)

// Cube implementation
type Cube struct {
	Size int
	Data map[int]map[int]map[int]color.Color
}

// New will create a new Cube struct with the given size
func New(size int) Cube {
	d := make(map[int]map[int]map[int]color.Color, size)

	for x := 0; x < size; x++ {
		yy := make(map[int]map[int]color.Color, size)
		for y := 0; y < size; y++ {
			z := make(map[int]color.Color, size)
			yy[y] = z
		}
		d[x] = yy
	}

	return Cube{
		Size: size,
		Data: d,
	}
}

// Get will return the color at a given point
func (c Cube) Get(x, y, z int) color.Color {
	return c.Data[x][y][z]
}

// Set will set a color for the given point
func (c Cube) Set(x, y, z int, val color.Color) {
	c.Data[x][y][z] = val
}
