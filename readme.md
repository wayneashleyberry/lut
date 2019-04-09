> LUT contains command-line tools and Go packages for applying colour adjustments to images using lookup tables.

[![wercker status](https://app.wercker.com/status/d6c0d4f2a9fbe670e8a1b11ad161a053/s/master "wercker status")](https://app.wercker.com/project/byKey/d6c0d4f2a9fbe670e8a1b11ad161a053)

### Motivation and caveats

This command-line tool and packages were written to inspect LUT's locally, and provide basic server-side rendering. This is most likely a bad idea, and you should probably be doing colour manipulations in OpenGL or a similar graphics programming framework.

There are tons of free LUT's available online, [LUTHOUSE](https://www.luthouse.com/free-luts) has some of my favourite.

### Quickstart

> :warning: You will need a `HOMEBREW_GITHUB_API_TOKEN` environment variable set to get access to our private homebrew tap.

```sh
brew tap overhq/homebrew-tap
brew install overhq/tap/lut
```

[Binaries are available for you to download](https://github.com/overhq/lut/releases/latest) if you don't want to [write Go code](https://golang.org/doc/code.html).

### Usage

```
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

This project uses [Go modules](https://blog.golang.org/modules2019), so make sure to clone it outside of your `$GOPATH`. You will need at least Go 1.11.

```sh
git clone git@github.com:overhq/lut.git
cd lut
go run main.go
```

### Supported Features

- 3D LUT's stored in the [`.cube` format](https://wwwimages2.adobe.com/content/dam/acom/en/products/speedgrade/cc/pdfs/cube-lut-specification-1.0.pdf)
- Image LUT's stored in 512x512 `jpeg` or `png` images
- Filter intensity
- Trilinear interpolation

### Not yet supported

- Image LUT's of arbitrary sizes
- 2d `.cube` files
