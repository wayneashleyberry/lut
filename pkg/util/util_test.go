package util

import (
	"reflect"
	"testing"
)

func TestParseFloats(t *testing.T) {
	type args struct {
		in      string
		bitSize int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			args: args{
				in:      "1.0",
				bitSize: 8,
			},
			want: []float64{1.0},
		},
		{
			args: args{
				in:      "1.1 1.2 1.3",
				bitSize: 8,
			},
			want: []float64{1.1, 1.2, 1.3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFloats(tt.args.in, tt.args.bitSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFloats() = %v, want %v", got, tt.want)
			}
		})
	}
}
