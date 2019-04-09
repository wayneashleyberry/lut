package main

import (
	"fmt"
	"os"

	"github.com/overhq/lut/cmd/apply"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "lut",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	cmd.AddCommand(apply.Command())

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
