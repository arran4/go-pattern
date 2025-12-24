package pattern

import (
	"fmt"
	"image"
	"image/color"
)

func ExampleNewRect() {
	// A simple black rectangle (default)
	_ = NewRect()
	// Output:
}

func ExampleNewRect_blue() {
	// A blue rectangle
	_ = NewRect(
		SetFillColor(color.RGBA{0, 0, 255, 255}),
	)
	// Output:
}

func ExampleNewRect_bounded() {
	// A red rectangle with specific bounds
	r := NewRect(
		SetFillColor(color.RGBA{255, 0, 0, 255}),
		SetBounds(image.Rect(0, 0, 100, 50)),
	)
	fmt.Printf("%v", r.Bounds())
	// Output: (0,0)-(100,50)
}
