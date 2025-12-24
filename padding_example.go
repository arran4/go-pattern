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
			// Ignore bounds passed from bootstrap to prevent overwriting centering logic
			return NewDemoPaddingCentered()
		},
		"TopLeft": func(bounds image.Rectangle) image.Image {
			return NewDemoPaddingTopLeft()
		},
		"BottomRight": func(bounds image.Rectangle) image.Image {
			return NewDemoPaddingBottomRight()
		},
		"Right": func(bounds image.Rectangle) image.Image {
			return NewDemoPaddingRight()
		},
		"AlignedWithOffset": func(bounds image.Rectangle) image.Image {
			return NewDemoPaddingAlignedWithOffset()
		},
	}, []string{"Color", "Bounded", "Centered", "TopLeft", "BottomRight", "Right", "AlignedWithOffset"}
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
	gopher := NewScale(NewGopher(), ScaleToRatio(0.2)) // Larger scale

	// Center the gopher in a 150x150 box (fills the available space)
	return NewCenter(gopher, 150, 150, image.NewUniform(color.RGBA{240, 240, 240, 255}))
}

func NewDemoPaddingTopLeft(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.2))
	return NewAligned(gopher, 150, 150, 0.0, 0.0, image.NewUniform(color.RGBA{220, 220, 220, 255}))
}

func NewDemoPaddingBottomRight(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.2))
	return NewAligned(gopher, 150, 150, 1.0, 1.0, image.NewUniform(color.RGBA{220, 220, 220, 255}))
}

func NewDemoPaddingRight(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.2))
	return NewAligned(gopher, 150, 150, 1.0, 0.5, image.NewUniform(color.RGBA{220, 220, 220, 255}))
}

func NewDemoPaddingAlignedWithOffset(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.2))
	// Align Top-Left (0,0) but add 20px padding on Top/Left to offset it.
	// We use variadic padding args: 20 (all sides).
	// Actually user asked for "padding on one or two of the aligned sides".
	// Let's do Top/Left alignment but with extra padding on Left/Top.
	// Using 4 args: Top=20, Right=0, Bottom=0, Left=20.
	return NewAligned(gopher, 150, 150, 0.0, 0.0, image.NewUniform(color.RGBA{200, 220, 200, 255}), 20, 0, 0, 20)
}
