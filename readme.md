> LUT contains command-line tools and Go packages for applying colour adjustments to images using lookup tables.

[![Go Report Card](https://goreportcard.com/badge/github.com/wayneashleyberry/lut)](https://goreportcard.com/report/github.com/wayneashleyberry/lut)
[![wercker status](https://app.wercker.com/status/3d33abdf103b7aba4e1b7d6283912523/s/master "wercker status")](https://app.wercker.com/project/byKey/3d33abdf103b7aba4e1b7d6283912523)

### Motivation and caveats

I wrote this command-line tool and package to inspect LUT's locally, and provide basic server-side rendering. This is most likely a bad idea, and you should probably be doing colour manipulations in OpenGL or a similar graphics programming framework.

Currently only 512px by 512px LUT's are supported — future work might be done to support other dimenions.

### Example

```sh
go run main.go apply testdata/sample.jpg --lut testdata/filter.png --out testdata/output.jpg
```

| Input                                                                                                     | LUT                                                                                                  | Output                                                                                           |
| --------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| ![an unfiltered image](https://raw.githubusercontent.com/wayneashleyberry/lut/master/testdata/sample.jpg) | ![a lookup table](https://raw.githubusercontent.com/wayneashleyberry/lut/master/testdata/filter.png) | ![the result](https://raw.githubusercontent.com/wayneashleyberry/lut/master/testdata/output.jpg) |

### Usage

```
Usage:
  lut [flags]
  lut [command]

Available Commands:
  apply       Adjust image colour according to a LUT
  help        Help about any command

Flags:
  -h, --help   help for lut

Use "lut [command] --help" for more information about a command.
```

[Binaries are available for you to download](https://github.com/wayneashleyberry/lut/releases/latest) if you don't want to [write Go code](https://golang.org/doc/code.html).

### Installation

This project uses [Go modules](https://blog.golang.org/modules2019), so make sure to clone it outside of your `$GOPATH`. You will need at least Go 1.11.

```sh
git clone git@github.com:wayneashleyberry/lut.git
cd lut
go run main.go
```
