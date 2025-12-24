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

func ExampleRect_stroked() {
	// A blue rectangle with a 10px green border
	r := NewRect(
		SetFillColor(color.RGBA{0, 0, 255, 255}),
		SetLineColor(color.RGBA{0, 255, 0, 255}),
		SetLineSize(10),
		SetBounds(image.Rect(0, 0, 100, 100)),
	)
	fmt.Printf("%v", r.Bounds())
	// Output: (0,0)-(100,100)
}

func ExampleRect_image_stroke() {
	// A yellow rectangle with a 20px border made of the Gopher image
	gopher := NewGopher()
	r := NewRect(
		SetFillColor(color.RGBA{255, 255, 0, 255}),
		SetLineImageSource(gopher),
		SetLineSize(20),
		SetBounds(image.Rect(0, 0, 100, 100)),
	)
	fmt.Printf("%v", r.Bounds())
	// Output: (0,0)-(100,100)
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
	// 20px margin around the 100x100 inner image.
	// Resulting size: 140x140.
	p := NewPadding(
		inner,
		PaddingMargin(20),
		PaddingBackground(bg),
	)

	fmt.Printf("%v", p.Bounds())
	// Output: (0,0)-(140,140)
}

func GenerateRect(b image.Rectangle) image.Image {
	return NewDemoRect(SetBounds(b))
}

func GenerateRectReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Red": func(b image.Rectangle) image.Image {
			return NewDemoRect(SetBounds(b), SetFillColor(color.RGBA{255, 0, 0, 255}))
		},
		"BorderPad": func(b image.Rectangle) image.Image {
			inner := NewHorizontalLine(
				SetBounds(image.Rect(0, 0, b.Dx()-40, b.Dy()-40)),
				SetLineSize(5),
				SetSpaceSize(5),
				SetLineColor(color.Black),
				SetSpaceColor(color.White),
			)
			borderBg := NewRect(SetFillColor(color.RGBA{100, 0, 0, 255}))
			return NewPadding(
				inner,
				PaddingMargin(20),
				PaddingBackground(borderBg),
				PaddingBoundary(b),
			)
		},
		"Stroked": func(b image.Rectangle) image.Image {
			return NewDemoRect(
				SetBounds(b),
				SetFillColor(color.RGBA{0, 0, 255, 255}),
				SetLineColor(color.RGBA{0, 255, 0, 255}),
				SetLineSize(10),
			)
		},
		"ImgStroke": func(b image.Rectangle) image.Image {
			gopher := NewGopher()
			return NewDemoRect(
				SetBounds(b),
				SetFillColor(color.RGBA{255, 255, 0, 255}),
				SetLineImageSource(gopher),
				SetLineSize(20),
			)
		},
		"Transparent": func(b image.Rectangle) image.Image {
			return NewDemoRect(
				SetBounds(b),
				SetFillColor(color.Transparent),
				SetLineColor(color.RGBA{255, 0, 255, 255}),
				SetLineSize(15),
			)
		},
		"CheckFrame": func(b image.Rectangle) image.Image {
			checker := NewChecker(color.Black, color.White)
			return NewDemoRect(
				SetBounds(b),
				SetFillColor(color.Transparent),
				SetLineImageSource(checker),
				SetLineSize(25),
			)
		},
	}, []string{"Red", "BorderPad", "Stroked", "ImgStroke", "Transparent", "CheckFrame"}
}

func init() {
	RegisterGenerator("Rect", GenerateRect)
	RegisterReferences("Rect", GenerateRectReferences)
}
