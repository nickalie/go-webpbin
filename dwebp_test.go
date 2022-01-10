package webpbin

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"golang.org/x/image/webp"
	"image/png"
)

func TestVersionDWebP(t *testing.T) {
	c := NewDWebP()
	r, err := c.Version()
	assert.Nil(t, err)
	if _, ok := os.LookupEnv("DOCKER_ARM_TEST"); !ok {
		assert.Equal(t, "1.2.0", r)
	}
}

func TestDecodeReader(t *testing.T) {
	c := NewDWebP()
	f, err := os.Open("source.webp")
	assert.Nil(t, err)
	defer f.Close()
	c.Input(f)
	c.OutputFile("target.png")
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	validatePng(t)
}

func TestDecodeFile(t *testing.T) {
	c := NewDWebP()
	c.InputFile("source.webp")
	c.OutputFile("target.png")
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	validatePng(t)
}

func TestDecodeImage(t *testing.T) {
	c := NewDWebP()
	f, err := os.Open("source.webp")
	assert.Nil(t, err)
	defer f.Close()
	imgSource, err := webp.Decode(f)
	assert.Nil(t, err)
	f.Seek(0, 0)
	c.Input(f)
	imgTarget, err := c.Run()
	assert.Nil(t, err)
	assert.NotNil(t, imgTarget)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}

func TestDecodeWriter(t *testing.T) {
	f, err := os.Create("target.png")
	assert.Nil(t, err)
	defer f.Close()
	c := NewDWebP()
	c.InputFile("source.webp")
	c.Output(f)
	img, err := c.Run()
	assert.Nil(t, err)
	assert.Nil(t, img)
	f.Close()
	validatePng(t)
}

func validatePng(t *testing.T) {
	defer os.Remove("target.png")
	fSource, err := os.Open("source.webp")
	assert.Nil(t, err)
	imgSource, err := webp.Decode(fSource)
	assert.Nil(t, err)
	fTarget, err := os.Open("target.png")
	assert.Nil(t, err)
	defer fTarget.Close()
	imgTarget, err := png.Decode(fTarget)
	assert.Nil(t, err)
	assert.Equal(t, imgSource.Bounds(), imgTarget.Bounds())
}
