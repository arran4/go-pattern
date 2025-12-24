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
	return NewDemoRect(SetBounds(b), SetFillColor(color.RGBA{255, 0, 0, 255}))
}

func init() {
	RegisterGenerator("Rect", GenerateRect)
}
