package webpbin

import (
	"errors"
	"fmt"
	"github.com/nickalie/go-binwrapper"
	"image"
	"io"
)

type cropInfo struct {
	x      int
	y      int
	width  int
	height int
}

// CWebP compresses an image using the WebP format. Input format can be either PNG, JPEG, TIFF, WebP or raw Y'CbCr samples.
// https://developers.google.com/speed/webp/docs/cwebp
type CWebP struct {
	*binwrapper.BinWrapper
	inputFile  string
	inputImage image.Image
	input      io.Reader
	outputFile string
	output     io.Writer
	quality    int
	crop       *cropInfo
}

// NewCWebP creates new CWebP instance.
func NewCWebP() *CWebP {
	bin := &CWebP{
		BinWrapper: createBinWrapper(),
		quality:    -1,
	}
	bin.ExecPath("cwebp")

	return bin
}

// Version returns cwebp version.
func (c *CWebP) Version() (string, error) {
	return version(c.BinWrapper)
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *CWebP) InputFile(file string) *CWebP {
	c.input = nil
	c.inputImage = nil
	c.inputFile = file
	return c
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *CWebP) Input(reader io.Reader) *CWebP {
	c.inputFile = ""
	c.inputImage = nil
	c.input = reader
	return c
}

// InputImage sets image to convert.
// InputFile or Input called before will be ignored.
func (c *CWebP) InputImage(img image.Image) *CWebP {
	c.inputFile = ""
	c.input = nil
	c.inputImage = img
	return c
}

// OutputFile specify the name of the output WebP file.
// Output called before will be ignored.
func (c *CWebP) OutputFile(file string) *CWebP {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specify writer to write webp file content.
// OutputFile called before will be ignored.
func (c *CWebP) Output(writer io.Writer) *CWebP {
	c.outputFile = ""
	c.output = writer
	return c
}

// Quality specify the compression factor for RGB channels between 0 and 100. The default is 75.
//
// A small factor produces a smaller file with lower quality. Best quality is achieved by using a value of 100.
func (c *CWebP) Quality(quality uint) *CWebP {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

// Crop the source to a rectangle with top-left corner at coordinates (x, y) and size width x height. This cropping area must be fully contained within the source rectangle.
func (c *CWebP) Crop(x, y, width, height int) *CWebP {
	c.crop = &cropInfo{x, y, width, height}
	return c
}

// Run starts cwebp with specified parameters.
func (c *CWebP) Run() error {
	defer c.BinWrapper.Reset()

	if c.quality > -1 {
		c.Arg("-q", fmt.Sprintf("%d", c.quality))
	}

	if c.crop != nil {
		c.Arg("-crop", fmt.Sprintf("%d", c.crop.x), fmt.Sprintf("%d", c.crop.y), fmt.Sprintf("%d", c.crop.width), fmt.Sprintf("%d", c.crop.height))
	}

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	c.Arg("-o", output)

	err = c.setInput()

	if err != nil {
		return err
	}

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

// Reset all parameters to default values
func (c *CWebP) Reset() *CWebP {
	c.crop = nil
	c.quality = -1
	return c
}

func (c *CWebP) setInput() error {
	if c.input != nil {
		c.Arg("--").Arg("-")
		c.StdIn(c.input)
	} else if c.inputImage != nil {
		r, err := createReaderFromImage(c.inputImage)

		if err != nil {
			return err
		}

		c.Arg("--").Arg("-")
		c.StdIn(r)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
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
