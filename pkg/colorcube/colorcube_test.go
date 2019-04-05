package colorcube

import (
	"image/color"
	"reflect"
	"testing"
)

func TestCube_Get(t *testing.T) {
	cube := New(32)
	want := color.NRGBA{255, 255, 255, 255}
	cube.Set(1, 2, 3, want)
	got := cube.Get(1, 2, 3)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Cube.Get() = %v, want %v", got, want)
	}
}
