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
		"Bounded": func(bounds image.Rectangle) image.Image {
			return NewDemoPaddingBounded(SetBounds(bounds))
		},
		"Centered": func(bounds image.Rectangle) image.Image {
			return NewDemoPaddingCentered(SetBounds(bounds))
		},
	}, []string{"Color", "Bounded", "Centered"}
}

func ExampleNewPadding(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.5))
	// Padding with transparent background (nil)
	return NewPadding(gopher, PaddingMargin(20))
}

func NewDemoPaddingColor(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.5))
	// Padding with checker background
	// We use a checker pattern as the background
	bg := NewChecker(color.Black, color.RGBA{200, 0, 0, 255})
	return NewPadding(gopher, PaddingMargin(20), PaddingBackground(bg))
}

func NewDemoPaddingBounded(ops ...func(any)) image.Image {
	// Padding bounding an Unbounded image
	// We use an unbounded checker pattern (default NewChecker)
	checker := NewChecker(color.Black, color.White)

	// Bound it to 100x100 with 10px padding (content 80x80)
	return NewPadding(checker, PaddingMargin(10), PaddingBoundary(image.Rect(0, 0, 100, 100)))
}

func NewDemoPaddingCentered(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.25))

	// Center the gopher in a 200x200 box with a light background
	return NewCenter(gopher, 200, 200, image.NewUniform(color.RGBA{240, 240, 240, 255}))
}
