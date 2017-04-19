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

type CWebP struct {
	*binwrapper.BinWrapper
	inputFile  string
	inputImage image.Image
	input      io.Reader
	outputFile string
	output     io.Writer
	quality    int
}

func NewCWebP() *CWebP {
	bin := &CWebP{
		BinWrapper: createBinWrapper(),
		quality:    -1,
	}
	bin.ExecPath("cwebp")

	return bin
}

func (c *CWebP) Version() (string, error) {
	return version(c.BinWrapper)
}

func (c *CWebP) InputFile(file string) *CWebP {
	c.input = nil
	c.inputImage = nil
	c.inputFile = file
	return c
}

func (c *CWebP) Input(reader io.Reader) *CWebP {
	c.inputFile = ""
	c.inputImage = nil
	c.input = reader
	return c
}

func (c *CWebP) InputImage(img image.Image) *CWebP {
	c.inputFile = ""
	c.input = nil
	c.inputImage = img
	return c
}

func (c *CWebP) OutputFile(file string) *CWebP {
	c.output = nil
	c.outputFile = file
	return c
}

func (c *CWebP) Output(writer io.Writer) *CWebP {
	c.outputFile = ""
	c.output = writer
	return c
}

func (c *CWebP) Quality(quality uint) *CWebP {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

func (c *CWebP) Run() error {
	defer c.BinWrapper.Reset()
	if c.quality > -1 {
		c.Arg("-q", fmt.Sprintf("%d", c.quality))
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