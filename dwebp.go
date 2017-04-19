package webpbin

import (
	"github.com/nickalie/go-binwrapper"
	"io"
	"image"
	"errors"
	"bytes"
	"image/png"
)

//decompresses WebP files into PNG
type DWebP struct {
	*binwrapper.BinWrapper
	inputFile  string
	input      io.Reader
	outputFile string
	output     io.Writer
}

func NewDWebP() *DWebP {
	bin := &DWebP{
		BinWrapper: createBinWrapper(),
	}
	bin.ExecPath("dwebp")

	return bin
}

func (c *DWebP) InputFile(file string) *DWebP {
	c.input = nil
	c.inputFile = file
	return c
}

func (c *DWebP) Input(reader io.Reader) *DWebP {
	c.inputFile = ""
	c.input = reader
	return c
}

func (c *DWebP) OutputFile(file string) *DWebP {
	c.output = nil
	c.outputFile = file
	return c
}

func (c *DWebP) Output(writer io.Writer) *DWebP {
	c.outputFile = ""
	c.output = writer
	return c
}

func (c *DWebP) Version() (string, error) {
	return version(c.BinWrapper)
}

func (c *DWebP) Run() (image.Image, error) {
	defer c.BinWrapper.Reset()

	input, err := c.getInput()

	if err != nil {
		return nil, err
	}

	output, err := c.getOutput()

	if err != nil {
		return nil, err
	}

	err = c.Arg(input).Arg("-o", output).Run()

	if err != nil {
		return nil, errors.New(string(c.StdErr))
	}

	if c.output != nil {
		_, err = io.Copy(c.output, bytes.NewReader(c.BinWrapper.StdOut))
		return nil, err
	} else if c.outputFile == "" {
		return png.Decode(bytes.NewReader(c.BinWrapper.StdOut))
	}

	return nil, nil
}

func (c *DWebP) getInput() (string, error) {
	if c.input != nil {
		return createFileFromReader(c.input)

	} else if c.inputFile != "" {
		return c.inputFile, nil
	} else {
		return "", errors.New("Undefined input")
	}
}

func (c *DWebP) getOutput() (string, error) {
	if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "-", nil
	}
}