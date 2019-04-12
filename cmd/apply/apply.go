package apply

import (
	"bufio"
	"errors"
	"image"
	"os"
	"path"
	"strings"

	"github.com/overhq/lut/pkg/cubelut"
	"github.com/overhq/lut/pkg/imagelut"
	"github.com/overhq/lut/pkg/trilinear"
	"github.com/overhq/lut/pkg/util"
	"github.com/spf13/cobra"
)

// Error types
var (
	ErrInvalidInterpolation = errors.New("invalid interpolation, accepted values are `none`, `tri` and `interp`")
)

// Command will create a new "apply" command
func Command() *cobra.Command {
	var lutfile, outfile string
	var intensity float64
	var interp string

	cmd := &cobra.Command{
		Use:   "apply [source.png] --lut sepia.png --out image.png --interp none",
		Short: "Adjust image colour according to a LUT",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			srcimg, err := util.ReadImage(args[0])
			if err != nil {
				util.Exit(err)
			}

			var out image.Image

			switch strings.ToLower(path.Ext(lutfile)) {
			case ".cube":
				file, err := os.Open(lutfile)
				if err != nil {
					util.Exit(err)
				}
				defer file.Close()

				r := bufio.NewReader(file)

				cubefile, err := cubelut.Parse(r)
				if err != nil {
					util.Exit(err)
				}

				switch interp {
				case "tri":
					cube := cubefile.Cube()

					img, err := trilinear.Interpolate(srcimg, cube, intensity)
					if err != nil {
						util.Exit(err)
					}

					out = img
				case "tetra":
					util.TODO()
				case "none":
					util.TODO()
				default:
					util.Exit(ErrInvalidInterpolation)
				}
			case ".png":
				fallthrough
			case ".jpg":
				fallthrough
			case ".jpeg":
				lutimg, err := util.ReadImage(lutfile)
				if err != nil {
					util.Exit(err)
				}

				switch interp {
				case "tri":
					cube, err := imagelut.Parse(lutimg)
					if err != nil {
						util.Exit(err)
					}

					img, err := trilinear.Interpolate(srcimg, cube, intensity)
					if err != nil {
						util.Exit(err)
					}

					out = img
				case "tetra":
					util.TODO()
				case "none":
					img, err := imagelut.Apply(srcimg, lutimg, intensity)
					if err != nil {
						util.Exit(err)
					}

					out = img
				default:
					util.Exit(ErrInvalidInterpolation)
				}

			default:
				util.Exit(errors.New("unsupported file type"))
			}

			if err := util.WriteImage(outfile, out); err != nil {
				util.Exit(err)
			}
		},
	}

	cmd.Flags().Float64VarP(&intensity, "intensity", "", 1, "Intensity of the applied effect")
	cmd.Flags().StringVarP(&interp, "interp", "i", "tri", "Interpolation")

	// Required flags
	cmd.Flags().StringVarP(&lutfile, "lut", "", "", "Path to LUT [required]")
	cmd.Flags().StringVarP(&outfile, "out", "o", "", "Path to write output [required]")

	_ = cmd.MarkFlagRequired("lut")
	_ = cmd.MarkFlagRequired("out")

	return cmd
}
