package convert

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/overhq/lut/pkg/colorcube"
	"github.com/overhq/lut/pkg/cubelut"
	"github.com/overhq/lut/pkg/imagelut"
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

			fmt.Println("PARSED")

			switch strings.ToLower(path.Ext(out)) {
			case ".cube":
				f := cubelut.FromColorCube(cube)

				fmt.Println("CREATED CUBE FILE")

				filename := filepath.Base(in)
				extension := filepath.Ext(filename)
				name := filename[0 : len(filename)-len(extension)]

				f.Title = name

				err := ioutil.WriteFile(out, f.Bytes(), 0777)
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
