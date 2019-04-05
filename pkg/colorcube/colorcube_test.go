package colorcube

import (
	"reflect"
	"testing"
)

func TestCube_Get(t *testing.T) {
	cube := New(32, []float64{0, 0, 0}, []float64{1, 1, 1})
	want := []float64{1, 1, 1}
	cube.Set(1, 2, 3, want)
	got := cube.Get(1, 2, 3)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Cube.Get() = %v, want %v", got, want)
	}
}
