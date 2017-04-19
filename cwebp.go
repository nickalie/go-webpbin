package webpbin

import (
	"github.com/nickalie/go-binwrapper"
	"io"
	"errors"
	"os"
	"image"
	"fmt"
	"bytes"
)

type cropInfo struct {
	x      int
	y      int
	width  int
	height int
}

//CWebP compresses an image using the WebP format. Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
type CWebP struct {
	*binwrapper.BinWrapper
	inputFile  string
	inputImage image.Image
	input      io.Reader
	outputFile string
	output     io.Writer
	quality    int
	crop *cropInfo
}

//Creates new CWebP instance
func NewCWebP() *CWebP {
	bin := &CWebP{
		BinWrapper: createBinWrapper(),
		quality:    -1,
	}
	bin.ExecPath("cwebp")

	return bin
}

//Returns cwebp's version number
func (c *CWebP) Version() (string, error) {
	return version(c.BinWrapper)
}

//Sets image file to convert
//Input or InputImage called before will be ignored
func (c *CWebP) InputFile(file string) *CWebP {
	c.input = nil
	c.inputImage = nil
	c.inputFile = file
	return c
}

//Sets reader to convert
//InputFile or InputImage called before will be ignored
func (c *CWebP) Input(reader io.Reader) *CWebP {
	c.inputFile = ""
	c.inputImage = nil
	c.input = reader
	return c
}

//Sets image to convert
//InputFile or Input called before will be ignored
func (c *CWebP) InputImage(img image.Image) *CWebP {
	c.inputFile = ""
	c.input = nil
	c.inputImage = img
	return c
}

//Specify the name of the output WebP file
//Output called before will be ignored
func (c *CWebP) OutputFile(file string) *CWebP {
	c.output = nil
	c.outputFile = file
	return c
}

//Specify writer to write webp file content
//OutputFile called before will be ignored
func (c *CWebP) Output(writer io.Writer) *CWebP {
	c.outputFile = ""
	c.output = writer
	return c
}

//Specify the compression factor for RGB channels between 0 and 100. The default is 75.
//
//A small factor produces a smaller file with lower quality. Best quality is achieved by using a value of 100.
func (c *CWebP) Quality(quality uint) *CWebP {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

//Crop the source to a rectangle with top-left corner at coordinates (x, y) and size width x height. This cropping area must be fully contained within the source rectangle.
func (c *CWebP) Crop(x, y, width, height int) *CWebP {
	c.crop = &cropInfo{x, y, width, height}
	return c
}

//Runs compression
func (c *CWebP) Run() error {
	defer c.BinWrapper.Reset()

	if c.quality > -1 {
		c.Arg("-q", fmt.Sprintf("%d", c.quality))
	}

	if c.crop != nil {
		c.Arg("-crop", fmt.Sprintf("%d", c.crop.x), fmt.Sprintf("%d", c.crop.y), fmt.Sprintf("%d", c.crop.width), fmt.Sprintf("%d", c.crop.height))
	}

	input, err := c.getInput()

	if err != nil {
		return err
	}

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	err = c.Arg(input).
		Arg("-o", output).
		Run()

	if err != nil {
		return errors.New(string(c.StdErr))
	}

	if c.inputFile == "" {
		os.Remove(input)
	}

	if c.output != nil {
		b := bytes.NewReader(c.BinWrapper.StdOut)
		_, err = io.Copy(c.output, b)
		return err
	}

	return nil
}

//Resets all parameters to default values
func (c *CWebP) Reset() *CWebP {
	c.crop = nil
	c.quality = -1
	return c
}

func (c *CWebP) getInput() (string, error) {
	if c.input != nil {
		return createFileFromReader(c.input)

	} else if c.inputImage != nil {
		return createFileFromImage(c.inputImage)
	} else if c.inputFile != "" {
		return c.inputFile, nil
	} else {
		return "", errors.New("Undefined input")
	}
}

func (c *CWebP) getOutput() (string, error) {
	if c.output != nil {
		return "-", nil
	} else if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "", errors.New("Undefined output")
	}
}
