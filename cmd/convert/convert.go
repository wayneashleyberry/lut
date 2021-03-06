package convert

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wayneashleyberry/lut/pkg/colorcube"
	"github.com/wayneashleyberry/lut/pkg/cubelut"
	"github.com/wayneashleyberry/lut/pkg/imagelut"
	"github.com/wayneashleyberry/lut/pkg/util"
)

// Command will create a new convert command.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert [source.png] target.cube",
		Short: "Convert a LUT file to a different format",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			in := args[0]
			out := args[1]

			var cube colorcube.Cube

			switch strings.ToLower(path.Ext(in)) {
			case ".cube":
				file, err := os.Open(in)
				if err != nil {
					util.Exit(err)
				}
				defer file.Close()

				r := bufio.NewReader(file)

				cubefile, err := cubelut.Parse(r)
				if err != nil {
					util.Exit(err)
				}

				cube = cubefile.Cube()
			case ".png":
				lutimg, err := util.ReadImage(in)
				if err != nil {
					util.Exit(err)
				}

				c, err := imagelut.Parse(lutimg)
				if err != nil {
					util.Exit(err)
				}

				cube = c
			default:
				util.Exit(errors.New("unsupported file type: " + in))
			}

			switch strings.ToLower(path.Ext(out)) {
			case ".cube":
				f := cubelut.FromColorCube(cube)

				filename := filepath.Base(in)
				extension := filepath.Ext(filename)
				name := filename[0 : len(filename)-len(extension)]

				f.Title = name

				err := ioutil.WriteFile(out, f.Bytes(), 0600)
				if err != nil {
					util.Exit(err)
				}
			case ".png":
				img := imagelut.FromColorCube(cube)

				err := util.WriteImage(out, img)
				if err != nil {
					util.Exit(err)
				}
			default:
				util.Exit(errors.New("unsupported file type: " + out))
			}
		},
	}

	return cmd
}
