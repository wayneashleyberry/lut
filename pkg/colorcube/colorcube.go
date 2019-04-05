// Package colorcube implements a very simple data structure for a 3d color cube
// and contains helper methods for getting and settings values at specific
// points. The underlying data structure is a multidimensional map which needs
// to be initialised to a specific size to keep things efficient.
package colorcube

// Cube implementation
type Cube struct {
	Size int
	Data map[int]map[int]map[int][]float64
}

// New will create a new Cube struct with the given size
func New(size int) Cube {
	d := make(map[int]map[int]map[int][]float64, size)

	for x := 0; x < size; x++ {
		yy := make(map[int]map[int][]float64, size)
		for y := 0; y < size; y++ {
			z := make(map[int][]float64, size)
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
func (c Cube) Get(x, y, z int) []float64 {
	return c.Data[x][y][z]
}

// Set will set a color for the given point
func (c Cube) Set(x, y, z int, val []float64) {
	c.Data[x][y][z] = val
}
