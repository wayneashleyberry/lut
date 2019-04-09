package apply

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"os"
	"path"
	"strings"

	"github.com/overhq/lut/pkg/cubelut"
	"github.com/overhq/lut/pkg/imagelut"
	"github.com/overhq/lut/pkg/transform"
	"github.com/overhq/lut/pkg/util"
	"github.com/spf13/cobra"
)

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// Command will create a new "apply" command
func Command() *cobra.Command {
	var lutfile, outfile string
	var intensity float64

	cmd := &cobra.Command{
		Use:   "apply [source.png] --lut sepia.png --out image.png",
		Short: "Adjust image colour according to a LUT",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			srcimg, err := util.ReadImage(args[0])
			if err != nil {
				exit(err)
			}

			var out image.Image

			switch strings.ToLower(path.Ext(lutfile)) {
			case ".cube":
				file, err := os.Open(lutfile)
				if err != nil {
					exit(err)
				}
				defer file.Close()

				r := bufio.NewReader(file)

				cubefile, err := cubelut.Parse(r)
				if err != nil {
					exit(err)
				}

				cube := cubefile.Cube()

				img, err := transform.Image(srcimg, cube, intensity)
				if err != nil {
					exit(err)
				}

				out = img
			case ".png":
			case ".jpg":
			case ".jpeg":
				lutimg, err := util.ReadImage(lutfile)
				if err != nil {
					exit(err)
				}

				cube, err := imagelut.Parse(lutimg)
				if err != nil {
					exit(err)
				}

				img, err := transform.Image(srcimg, cube, intensity)
				if err != nil {
					exit(err)
				}

				out = img
			default:
				exit(errors.New("unsupported file type"))
			}

			if err := util.WriteImage(outfile, out); err != nil {
				exit(err)
			}
		},
	}

	cmd.Flags().Float64VarP(&intensity, "intensity", "", 1, "Intensity of the applied effect")
	cmd.Flags().StringVarP(&lutfile, "lut", "", "", "Path to LUT [required]")
	cmd.Flags().StringVarP(&outfile, "out", "o", "", "Path to write output [required]")

	_ = cmd.MarkFlagRequired("lut")
	_ = cmd.MarkFlagRequired("out")

	return cmd
}