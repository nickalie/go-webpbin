package webpbin

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/webp"
)

func TestDecode(t *testing.T) {
	f, err := os.Open("source.webp")
	assert.Nil(t, err)
	imgSource, err := Decode(f)
	assert.Nil(t, err)
	f.Seek(0, 0)
	imgTarget, err := webp.Decode(f)
	assert.Nil(t, err)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}
