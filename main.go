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
	cmd := &cobra.Command{
		Use: "lut",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	cmd.AddCommand(
		apply.Command(),
		convert.Command(),
	)

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("lut v%s (%s)\n", version, date)
		},
	})

	t := time.Now()

	if err := cmd.Execute(); err != nil {
		util.Exit(err)
	}

	fmt.Println(time.Since(t))
}
