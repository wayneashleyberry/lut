package main

import (
	"fmt"

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

	if err := root.Execute(); err != nil {
		util.Exit(err)
	}
}
