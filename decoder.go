package webpbin

import (
	"image"
	"io"
)

// Decode reads a WebP image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	return NewDWebP().Input(r).Run()
}
