package main

import (
	"bufio"
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

func main() {
	cmd := &cobra.Command{
		Use: "lut",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	var lutfile, outfile string
	var intensity float64

	apply := &cobra.Command{
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

				lut, err := cubelut.Parse(r)
				if err != nil {
					exit(err)
				}

				cube := lut.ColorCube()

				img, err := transform.Image(srcimg, cube, intensity)
				if err != nil {
					exit(err)
				}

				out = img

			default:
				lutimg, err := util.ReadImage(lutfile)
				if err != nil {
					exit(err)
				}

				img, err := imagelut.Apply(srcimg, lutimg, intensity)
				if err != nil {
					exit(err)
				}

				out = img
			}

			if err := util.WriteImage(outfile, out); err != nil {
				exit(err)
			}
		},
	}

	apply.Flags().Float64VarP(&intensity, "intensity", "", 1, "Intensity of the applied effect")
	apply.Flags().StringVarP(&lutfile, "lut", "", "", "Path to LUT [required]")
	apply.Flags().StringVarP(&outfile, "out", "o", "", "Path to write output [required]")

	_ = apply.MarkFlagRequired("lut")
	_ = apply.MarkFlagRequired("out")

	cmd.AddCommand(apply)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
