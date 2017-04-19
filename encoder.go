package webpbin

import (
	"image"
	"io"
)

type Encoder struct {
	Quality uint
}

func (e *Encoder) Encode(w io.Writer, m image.Image) error {
	c := NewCWebP()
	c.Quality(e.Quality)
	c.InputImage(m)
	c.Output(w)
	return c.Run()

}

func Encode(w io.Writer, m image.Image) error {
	e := &Encoder{Quality: 75}
	return e.Encode(w, m)
}
