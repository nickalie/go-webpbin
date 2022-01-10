# WebP Encoder/Decoder for Golang

[![](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/nickalie/go-webpbin)
[![](https://circleci.com/gh/nickalie/go-webpbin.png?circle-token=ebaa6a739ac4dc96dcb167e0700dcc699409f672)](https://circleci.com/gh/nickalie/go-webpbin)

WebP Encoder/Decoder for Golang based on official libwebp distribution

## Install

```go get -u github.com/nickalie/go-webpbin```

## Available env
All env can be override with option functions.

|Name|Default|Desscription|
|-----|------|------------|
|SKIP_DOWNLOAD|`false`|Download webp bin automatically. Since there is no precompiled file for alpine, **THE SKIP_DOWNLOAD MUST BE true AND ASSIGN A SOURCE FOR RUN.**|
|VENDOR_PATH|`.bin/webp`|When there is no lib within and `SKIP_DOWNLAOD` is not `true`, it'll be downloaded.|
|LIBWEBP_VERSION|`1.2.0`|The latest version for now. (2021/07/16)|


## Example of usage

```go
package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"github.com/nickalie/go-webpbin"
)

func main() {
	const width, height = 256, 256

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	f, err := os.Create("image.webp")
	if err != nil {
		log.Fatal(err)
	}

	if err := webpbin.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
```

## CWebP

CWebP is a wrapper for *cwebp* command line tool.

Example to convert image.png to image.webp:

```go
err := webpbin.NewCWebP().
		Quality(80).
		InputFile("image.png").
		OutputFile("image.webp").
		Run()
```

## DWebP

DWebP is a wrapper for *dwebp* command line tool.

Example to convert image.webp to image.png:

```go
err := webpbin.NewDWebP().
		InputFile("image.webp").
		OutputFile("image.png").
		Run()
```

## libwebp distribution

Under the hood library uses [official libwebp distribution](https://storage.googleapis.com/downloads.webmproject.org/releases/webp/index.html), so if you're going to use it on not supported platform (arm or alpine), you need to build libwebp from sources and set ```SKIP_DOWNLOAD=true```.

Snippet to build libweb on alpine:

```sh
apk add --no-cache --update libpng-dev libjpeg-turbo-dev giflib-dev tiff-dev autoconf automake make gcc g++ wget

wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-0.6.0.tar.gz && \
tar -xvzf libwebp-0.6.0.tar.gz && \
mv libwebp-0.6.0 libwebp && \
rm libwebp-0.6.0.tar.gz && \
cd /libwebp && \
./configure && \
make && \
make install && \
rm -rf libwebp
```
