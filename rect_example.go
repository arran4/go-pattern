package pattern

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var RectOutputFilename = "rect.png"
var RectZoomLevels = []int{}

const RectOrder = 20

// Rect Pattern
// A pattern that draws a filled rectangle.
func ExampleNewRect() {
	// A simple black rectangle (default)
	i := NewRect()
	// Output:

	// Create the file for the example
	f, err := os.Create(RectOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func ExampleRect_blue() {
	// A blue rectangle
	_ = NewRect(
		SetFillColor(color.RGBA{0, 0, 255, 255}),
	)
	// Output:
}

func ExampleRect_bounded() {
	// A red rectangle with specific bounds
	r := NewRect(
		SetFillColor(color.RGBA{255, 0, 0, 255}),
		SetBounds(image.Rect(0, 0, 100, 50)),
	)
	fmt.Printf("%v", r.Bounds())
	// Output: (0,0)-(100,50)
}

func GenerateRect(b image.Rectangle) image.Image {
	v1 := NewDemoRect(SetBounds(b), SetFillColor(color.RGBA{255, 0, 0, 255}))

	// Border Demo
	// Inner image: Horizontal Lines
	// Background: Red Rect
	// Padding makes the border
	inner := NewHorizontalLine(
		SetBounds(image.Rect(0, 0, b.Dx()-20, b.Dy()-20)),
		SetLineSize(5),
		SetSpaceSize(5),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
	)

	border := NewRect(SetFillColor(color.RGBA{255, 0, 0, 255}))

	v2 := NewPadding(
		inner,
		PaddingMargin(10),
		PaddingBackground(border),
		PaddingBoundary(b),
	)

	return stitchImagesForDemo(v1, v2)
}

func ExampleRect_border_demo() {
	// Creating a border using Padding and a Rect background.
	// Inner pattern: Horizontal lines
	inner := NewHorizontalLine(
		SetLineSize(5),
		SetSpaceSize(5),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
		SetBounds(image.Rect(0, 0, 100, 100)),
	)

	// Background pattern: Red Rectangle
	bg := NewRect(SetFillColor(color.RGBA{255, 0, 0, 255}))

	// Apply Padding with background
	// 10px margin around the 100x100 inner image.
	// Resulting size: 120x120.
	p := NewPadding(
		inner,
		PaddingMargin(10),
		PaddingBackground(bg),
	)

	fmt.Printf("%v", p.Bounds())
	// Output: (0,0)-(120,120)
}

func init() {
	RegisterGenerator("Rect", GenerateRect)
}
