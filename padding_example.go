package pattern

import (
	"image"
	"image/color"
)

var (
	PaddingOutputFilename = "padding.png"
	PaddingZoomLevels     = []int{2}
)

const (
	PaddingOrder          = 8
	PaddingBaseLabel      = "Padding"
)

func init() {
	RegisterGenerator("Padding", func(bounds image.Rectangle) image.Image {
		return ExampleNewPadding(SetBounds(bounds))
	})
	RegisterReferences("Padding", BootstrapPaddingReferences)
}

func BootstrapPaddingReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Color": func(bounds image.Rectangle) image.Image {
			return NewDemoPaddingColor(SetBounds(bounds))
		},
	}, []string{"Color"}
}

func ExampleNewPadding(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleFactor(0.5))
	// Padding with transparent background (nil)
	return NewPadding(gopher, 20, nil)
}

func NewDemoPaddingColor(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleFactor(0.5))
	// Padding with checker background
	// We use a checker pattern as the background
	bg := NewChecker(color.Black, color.RGBA{200, 0, 0, 255})
	return NewPadding(gopher, 20, bg)
}
