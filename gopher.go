package pattern

import (
	_ "embed"
	"image"
	"image/png"
	"strings"
)

//go:embed assets/gopher.png
var gopherPng []byte

//go:embed assets/go.png
var goPng []byte

// NewGopher returns an image of the Go Gopher.
// If the embedded asset cannot be decoded, it panics (should not happen in production).
func NewGopher() image.Image {
	img, err := png.Decode(strings.NewReader(string(gopherPng)))
	if err != nil {
		panic("failed to decode embedded gopher image: " + err.Error())
	}
	return img
}

// NewGoLogo returns an image of the Go Logo (or a Gopher related image).
func NewGoLogo() image.Image {
	img, err := png.Decode(strings.NewReader(string(goPng)))
	if err != nil {
		panic("failed to decode embedded go logo image: " + err.Error())
	}
	return img
}
