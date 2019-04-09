package convert

import (
	"errors"
	"path"
	"strings"

	"github.com/overhq/lut/pkg/util"
	"github.com/spf13/cobra"
)

// Command will create a new convert command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert [source.png] target.cube",
		Short: "Convert a LUT file to a different format",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			in := args[0]
			out := args[1]

			switch strings.ToLower(path.Ext(in)) {
			case ".cube":
				// TODO
			case ".png":
			case ".jpeg":
			case ".jpg":
				// TODO
			default:
				util.Exit(errors.New("unsupported file type: " + in))
			}

			switch strings.ToLower(path.Ext(out)) {
			case ".cube":
				// TODO
			case ".png":
			case ".jpeg":
			case ".jpg":
				// TODO
			default:
				util.Exit(errors.New("unsupported file type: " + out))
			}
		},
	}

	return cmd
}
