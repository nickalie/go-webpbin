package webpbin

import (
	"image"
	"io"
)

func Decode(r io.Reader) (image.Image, error) {
	return NewDWebP().Input(r).Run()
}