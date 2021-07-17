package webpbin

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/nickalie/go-binwrapper"
)

var skipDownload bool
var dest = ".bin/webp"
var libwebpVersion = "1.2.0"

type OptionFunc func(binWrapper *binwrapper.BinWrapper) error

func SetSkipDownload(isSkipDownload bool) OptionFunc {
	return func(binWrapper *binwrapper.BinWrapper) error {
		skipDownload = isSkipDownload
		return nil
	}
}

func SetVendorPath(path string) OptionFunc {
	return func(binWrapper *binwrapper.BinWrapper) error {
		dest = path
		return nil
	}
}

func loadDefaultFromENV(binWrapper *binwrapper.BinWrapper) error {
	if os.Getenv("SKIP_DOWNLOAD") == "true" {
		skipDownload = true
	}

	if path := os.Getenv("VENDOR_PATH"); path != "" {
		dest = path
	}

	if version := os.Getenv("LIBWEBP_VERSION"); version != "" {
		libwebpVersion = version
	}

	return nil
}

// DetectUnsupportedPlatforms detects platforms without prebuilt binaries (alpine and arm).
// For this platforms libwebp tools should be built manually.
// See https://github.com/nickalie/go-webpbin/blob/master/docker/Dockerfile and https://github.com/nickalie/go-webpbin/blob/master/docker/Dockerfile.arm for details
func DetectUnsupportedPlatforms() {
	if runtime.GOARCH == "arm" {
		skipDownload = true
	} else if runtime.GOOS == "linux" {
		output, err := ioutil.ReadFile("/etc/issue")

		if err == nil && bytes.Contains(bytes.ToLower(output), []byte("alpine")) {
			skipDownload = true
		}
	}
}

func createBinWrapper(optionFuncs ...OptionFunc) *binwrapper.BinWrapper {
	macVersionMap := map[string]string{
		"0.4.1":     "10.8-2",
		"0.4.1-rc1": "10.8",
		"0.4.2":     "10.8",
		"0.4.2-rc2": "10.8",
		"0.4.3":     "10.9",
		"0.4.3-rc1": "10.9",
		"0.4.4":     "10.9",
		"0.4.4-rc2": "10.9",
		"0.5.0":     "10.9",
		"0.5.0-rc1": "10.9",
		"0.5.1":     "10.9",
		"0.5.1-rc5": "10.9",
		"0.5.2":     "10.9",
		"0.5.2-rc2": "10.9",
		"0.6.0":     "10.12",
		"0.6.0-rc2": "10.12",
		"0.6.0-rc3": "10.12",
		"0.6.1":     "10.12",
		"0.6.1-rc2": "10.12",
		"1.0.0":     "10.13",
		"1.0.0-rc1": "10.13",
		"1.0.0-rc2": "10.13",
		"1.0.0-rc3": "10.13",
		"1.0.1":     "10.13",
		"1.0.1-rc2": "10.13",
		"1.0.2":     "10.14",
		"1.0.2-rc1": "10.14",
		"1.0.3":     "10.14",
		"1.0.3-rc1": "10.14",
		"1.1.0":     "10.15",
		"1.1.0-rc2": "10.15",
		"1.2.0":     "10.15",
		"1.2.0-rc3": "10.15",
	}
	base := "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/"

	b := binwrapper.NewBinWrapper().AutoExe()

	loadDefaultFromENV(b)

	for _, optionFunc := range optionFuncs {
		optionFunc(b)
	}

	if !skipDownload {
		b.Src(
			binwrapper.NewSrc().
				URL(base + "libwebp-" + libwebpVersion + "-mac-" + macVersionMap[libwebpVersion] + ".tar.gz").
				Os("darwin")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-linux-x86-32.tar.gz").
					Os("linux").
					Arch("x86")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-linux-x86-64.tar.gz").
					Os("linux").
					Arch("x64")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-windows-x64.zip").
					Os("win32").
					Arch("x64")).
			Src(
				binwrapper.NewSrc().
					URL(base + "libwebp-" + libwebpVersion + "-windows-x86.zip").
					Os("win32").
					Arch("x86"))
	}

	return b.Strip(2).Dest(dest)
}

func createReaderFromImage(img image.Image) (io.Reader, error) {
	enc := &png.Encoder{
		CompressionLevel: png.NoCompression,
	}

	var buffer bytes.Buffer
	err := enc.Encode(&buffer, img)

	if err != nil {
		return nil, err
	}

	return &buffer, nil
}

func version(b *binwrapper.BinWrapper) (string, error) {
	b.Reset()
	err := b.Run("-version")

	if err != nil {
		return "", err
	}

	version := string(b.StdOut())
	version = strings.Replace(version, "\n", "", -1)
	version = strings.Replace(version, "\r", "", -1)
	return version, nil
}
