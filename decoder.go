package webpbin

import (
	"image"
	"io"
)

//Reads a WebP image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	return NewDWebP().Input(r).Run()
}