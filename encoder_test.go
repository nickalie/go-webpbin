package webpbin

import (
	"testing"
	"image/jpeg"
	"os"
	"github.com/stretchr/testify/assert"
	"bytes"
	"golang.org/x/image/webp"
)

func TestEncode(t *testing.T) {
	f, err := os.Open("source.jpg")
	assert.Nil(t, err)
	imgSource, err := jpeg.Decode(f)
	assert.Nil(t, err)
	var b bytes.Buffer
	err = Encode(&b, imgSource)
	assert.Nil(t, err)
	imgTarget, err := webp.Decode(bytes.NewReader(b.Bytes()))
	assert.Nil(t, err)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}
