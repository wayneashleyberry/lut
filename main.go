package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wayneashleyberry/lut/pkg/lut"
	"github.com/wayneashleyberry/lut/pkg/util"
)

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	cmd := &cobra.Command{
		Use: "lut",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	var lutfile, outfile string

	apply := &cobra.Command{
		Use:   "apply [source.png] --lut sepia.png --out image.png",
		Short: "Adjust image colour according to a LUT",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			srcimg, err := util.ReadImage(args[0])
			if err != nil {
				exit(err)
			}

			lutimg, err := util.ReadImage(lutfile)
			if err != nil {
				exit(err)
			}

			img, err := lut.Apply(srcimg, lutimg)
			if err != nil {
				exit(err)
			}

			if err := util.WriteImage(outfile, img); err != nil {
				exit(err)
			}
		},
	}

	apply.Flags().StringVarP(&lutfile, "lut", "", "", "Path to LUT [required]")
	apply.MarkFlagRequired("lut")

	apply.Flags().StringVarP(&outfile, "out", "o", "", "Path to write output [required]")
	apply.MarkFlagRequired("out")

	cmd.AddCommand(apply)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
