package webpbin

import (
	"image"
	"image/png"
	"github.com/nickalie/go-binwrapper"
	"io/ioutil"
	"io"
	"strings"
	"runtime"
	"os/exec"
	"bytes"
)

var skipDownload bool
var dest string = "vendor/webp"

func DetectUnsupportedPlatforms()  {
	if runtime.GOARCH == "arm" {
		SkipDownload()
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("cat", "/etc/issue")
		output, err := cmd.CombinedOutput()

		if err == nil && bytes.Contains(bytes.ToLower(output), []byte("alpine")) {
			SkipDownload()
		}
	}
}

func SkipDownload() {
	skipDownload = true
	dest = ""
}

func Dest(value string) {
	dest = value
}

func createBinWrapper() *binwrapper.BinWrapper {
	base := "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/"

	b := binwrapper.NewBinWrapper().AutoExe()

	if !skipDownload {
		b.Src(
			binwrapper.NewSrc().
				Url(base + "libwebp-0.6.0-mac-10.12.tar.gz").
				Os("darwin")).
			Src(
			binwrapper.NewSrc().
				Url(base + "libwebp-0.6.0-linux-x86-32.tar.gz").
				Os("linux").
				Arch("x86")).
			Src(
			binwrapper.NewSrc().
				Url(base + "libwebp-0.6.0-linux-x86-64.tar.gz").
				Os("linux").
				Arch("x64")).
			Src(
			binwrapper.NewSrc().
				Url(base + "libwebp-0.6.0-windows-x64.zip").
				Os("win32").
				Arch("x64")).
			Src(
			binwrapper.NewSrc().
				Url(base + "libwebp-0.6.0-windows-x86.zip").
				Os("win32").
				Arch("x86"))
	}

	return b.Strip(2).Dest(dest)
}

func createFileFromReader(reader io.Reader) (string, error) {
	file, err := ioutil.TempFile("", "gocwebpbin")

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func createFileFromImage(img image.Image) (string, error) {
	file, err := ioutil.TempFile("", "gocwebpbin")

	if err != nil {
		return "", err
	}

	defer file.Close()

	enc := &png.Encoder{
		CompressionLevel: png.NoCompression,
	}
	err = enc.Encode(file, img)

	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func version(b *binwrapper.BinWrapper) (string, error) {
	b.Reset()
	err := b.Run("-version")

	if err != nil {
		return "", err
	}

	version := string(b.StdOut)
	version = strings.Replace(version, "\n", "", -1)
	version = strings.Replace(version, "\r", "", -1)
	return version, nil
}
