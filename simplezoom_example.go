package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var SimpleZoomOutputFilename = "simplezoom.png"
var SimpleZoomZoomLevels = []int{2, 4}

const SimpleZoomOrder = 2

// Simple Zoom Pattern
// Zooms in on an underlying image.
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

func BootstrapSimpleZoom(b image.Rectangle) image.Image {
	return NewDemoChecker(SetBounds(b))
}

func init() {
	RegisterGenerator("SimpleZoom", BootstrapSimpleZoom)
}
