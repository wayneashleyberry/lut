package cubelut

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	lut, err := os.Open("./testdata/testlut.cube")
	if err != nil {
		t.Fatal("could not open file")
	}

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    CubeFile
		wantErr bool
	}{
		{
			args: args{
				r: lut,
			},
			wantErr: false,
			want: CubeFile{
				Dimensions: 3,
				DomainMax:  []float64{2.0, 2.0, 2.0},
				DomainMin:  []float64{0.0, 0.0, 0.0},
				Size:       2,
				Title:      "Hello, World!",
				R:          []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
				G:          []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
				B:          []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
