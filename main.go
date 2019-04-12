package main

import (
	"fmt"
	"time"

	"github.com/overhq/lut/cmd/apply"
	"github.com/overhq/lut/cmd/convert"
	"github.com/overhq/lut/pkg/util"
	"github.com/spf13/cobra"
)

var version string
var date string

func main() {
	root := &cobra.Command{
		Use: "lut",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	root.AddCommand(
		apply.Command(),
		convert.Command(),
	)

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("lut v%s (%s)\n", version, date)
		},
	})

	var verbose bool
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	t := time.Now()

	if err := root.Execute(); err != nil {
		util.Exit(err)
	}

	if verbose {
		fmt.Println(time.Since(t))
	}
}
