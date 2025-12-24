package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var SimpleZoomOutputFilename = "simplezoom.png"
var SimpleZoomZoomLevels = []int{2, 4}

const SimpleZoomOrder = 30

// SimpleZoom Pattern
// Scales an input pattern by a factor.
func ExampleNewSimpleZoom() {
	i := NewSimpleZoom(NewChecker(color.Black, color.White), 2)
	f, err := os.Create(SimpleZoomOutputFilename)
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

func GenerateSimpleZoom(b image.Rectangle) image.Image {
	// SimpleZoom needs an input image.
	// We can use a Checker pattern as the input for the demo.
	return NewDemoSimpleZoom(NewChecker(color.Black, color.White), SetBounds(b))
}

func init() {
	RegisterGenerator("SimpleZoom", GenerateSimpleZoom)
}
