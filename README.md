<p align="center">
  <img width="256" height="256" src="https://user-images.githubusercontent.com/727262/81205527-bb7a4000-8fc2-11ea-8ccf-46f91ab08c91.png">
</p>

> LUT contains command-line tools and Go packages for applying colour adjustments to images using lookup tables. [I gave a talk on the subject at London Gophers.](https://www.youtube.com/watch?v=KVmDATg2mCE)

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wayneashleyberry/lut?tab=overview)
![Go](https://github.com/wayneashleyberry/lut/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/wayneashleyberry/lut)](https://goreportcard.com/report/github.com/wayneashleyberry/lut)
![CodeQL](https://github.com/wayneashleyberry/lut/workflows/CodeQL/badge.svg)

### Motivation

This command-line tool and packages were written to inspect LUT's locally, and provide basic server-side rendering. This is most likely a bad idea, and you should probably be doing colour manipulations in OpenGL or a similar graphics programming framework.

There are tons of free LUT's available online, [luthouse.com](https://www.luthouse.com/free-luts) is a great example.

### Usage

```sh
Usage:
  lut [flags]
  lut [command]

Available Commands:
  apply       Adjust image colour according to a LUT
  convert     Convert a LUT file to a different format
  help        Help about any command
  version     Print version information

Flags:
  -h, --help   help for lut

Use "lut [command] --help" for more information about a command.
```

### Installation

This project uses [Go modules](https://blog.golang.org/modules2019), so make sure to clone it outside of your `$GOPATH`. You will need at least Go 1.12.

```sh
git clone git@github.com:wayneashleyberry/lut.git
cd lut
go mod download
go run main.go
```

### Supported Features

- 3D LUT's stored in the [`.cube` format](https://wwwimages2.adobe.com/content/dam/acom/en/products/speedgrade/cc/pdfs/cube-lut-specification-1.0.pdf) (recommended)
- Squar image LUT's stored in 512x512 `jpeg` or `png` images
- Filter intensity
- Trilinear interpolation

### Not yet supported

- Image LUT's of arbitrary sizes
- 2d `.cube` files
