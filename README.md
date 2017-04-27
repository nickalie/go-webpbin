# WebP Encoder/Decoder for Golang

[![](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/nickalie/go-webpbin)
[![](https://circleci.com/gh/nickalie/go-webpbin.png?circle-token=ebaa6a739ac4dc96dcb167e0700dcc699409f672)](https://circleci.com/gh/nickalie/go-webpbin)

WebP Encoder/Decoder for Golang based on official libwebp distribution

## Install

```go get -u github.com/nickalie/go-webpbin```

## Example of usage

```
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

```
err := webpbin.NewCWebP().
		Quality(80).
		InputFile("image.png").
		OutputFile("image.webp").
		Run()
```

## DWebP

DWebP is a wrapper for *dwebp* command line tool.

Example to convert image.webp to image.png:

```
err := webpbin.NewDWebP().
		InputFile("image.webp").
		OutputFile("image.png").
		Run()
```

## libwebp distribution

Under the hood library uses [official libwebp distribution](https://storage.googleapis.com/downloads.webmproject.org/releases/webp/index.html), so if you're going to use it on not supported platform (arm or alpine), you need to build libwebp from sources and call ```webpbin.SkipDownload()``` method.

Snippet to build libweb on alpine:

```
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
