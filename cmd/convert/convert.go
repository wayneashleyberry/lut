package convert

import "github.com/spf13/cobra"

// Command will create a new convert command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert [source.png] target.cube",
		Short: "Convert a LUT file to a different format",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			//
		},
	}

	return cmd
}
