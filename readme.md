> LUT contains command-line tools and Go packages for applying colour adjustments to images using lookup tables.

[![Go Report Card](https://goreportcard.com/badge/github.com/wayneashleyberry/lut)](https://goreportcard.com/report/github.com/wayneashleyberry/lut)

### Motivation and caveats

I wrote this command-line tool and package to inspect LUT's locally, and provide basic server-side rendering. This is most likely a bad idea, and you should probably be doing colour manipulations in OpenGL or a similar graphics programming framework.

### Example

```sh
go run main.go apply testdata/sample.jpg --lut testdata/filter.png --out testdata/output.jpg
```

| Input                                                                                     | LUT                                                                                       | Output                                                                                    |
| ----------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| ![...](https://raw.githubusercontent.com/wayneashleyberry/lut/master/testdata/sample.jpg) | ![...](https://raw.githubusercontent.com/wayneashleyberry/lut/master/testdata/filter.png) | ![...](https://raw.githubusercontent.com/wayneashleyberry/lut/master/testdata/output.jpg) |

### Usage

```
Usage:
lut [flags]
lut [command]

Available Commands:
apply Adjust image colour according to a LUT
help Help about any command

Flags:
-h, --help help for lut

Use "lut [command] --help" for more information about a command.
```

### Installation

This project uses [Go modules](https://blog.golang.org/modules2019), so make sure to clone it outside of your `$GOPATH`. You will need at least Go 1.11.

```sh
git clone git@github.com:wayneashleyberry/lut.git
cd lut
go run main.go
```
