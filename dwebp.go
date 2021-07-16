package webpbin

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"

	"github.com/nickalie/go-binwrapper"
)

// DWebP wraps dwebp tool used for decompression of WebP files into PNG.
// https://developers.google.com/speed/webp/docs/dwebp
type DWebP struct {
	*binwrapper.BinWrapper
	inputFile  string
	input      io.Reader
	outputFile string
	output     io.Writer
}

// NewDWebP creates new WebP instance
func NewDWebP(optionFuncs ...OptionFunc) *DWebP {
	bin := &DWebP{
		BinWrapper: createBinWrapper(optionFuncs...),
	}
	bin.ExecPath("dwebp")

	return bin
}

// InputFile sets webp file to convert.
// Input or InputImage called before will be ignored.
func (c *DWebP) InputFile(file string) *DWebP {
	c.input = nil
	c.inputFile = file
	return c
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *DWebP) Input(reader io.Reader) *DWebP {
	c.inputFile = ""
	c.input = reader
	return c
}

// OutputFile specify the name of the output image file.
// Output called before will be ignored.
func (c *DWebP) OutputFile(file string) *DWebP {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specify writer to write image file content.
// OutputFile called before will be ignored.
func (c *DWebP) Output(writer io.Writer) *DWebP {
	c.outputFile = ""
	c.output = writer
	return c
}

// Version returns dwebp version.
func (c *DWebP) Version() (string, error) {
	return version(c.BinWrapper)
}

// Run starts dwebp with specified parameters.
func (c *DWebP) Run() (image.Image, error) {
	defer c.BinWrapper.Reset()

	output, err := c.getOutput()

	if err != nil {
		return nil, err
	}

	c.Arg("-o", output)

	err = c.setInput()

	if err != nil {
		return nil, err
	}

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return nil, errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	if c.output == nil && c.outputFile == "" {
		return png.Decode(bytes.NewReader(c.BinWrapper.StdOut()))
	}

	return nil, nil
}

func (c *DWebP) setInput() error {
	if c.input != nil {
		c.Arg("--").Arg("-")
		c.StdIn(c.input)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

func (c *DWebP) getOutput() (string, error) {
	if c.outputFile != "" {
		return c.outputFile, nil
	}

	return "-", nil
}
