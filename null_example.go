package pattern

import (
	"image"
	"image/png"
	"os"
)

var NullOutputFilename = "null.png"
var NullZoomLevels = []int{}

const NullOrder = 0

// Null Pattern
// Returns a transparent color for all pixels.
func ExampleNewNull() {
	i := NewNull()
	f, err := os.Create(NullOutputFilename)
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

func GenerateNull(b image.Rectangle) image.Image {
	return NewDemoNull(SetBounds(b))
}

func init() {
	RegisterGenerator("Null", GenerateNull)
}
