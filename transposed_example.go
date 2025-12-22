package pattern

import (
	"image"
	"image/png"
	"os"
)

var TransposedOutputFilename = "transposed.png"
var TransposedZoomLevels = []int{2, 4}

const TransposedOrder = 31

// Transposed Pattern
// Transposes the coordinates of an input pattern.
func ExampleNewTransposed() {
	i := NewTransposed(NewDemoNull(), 10, 10)
	f, err := os.Create(TransposedOutputFilename)
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

func GenerateTransposed(b image.Rectangle) image.Image {
	return NewDemoTransposed(SetBounds(b))
}

func init() {
	RegisterGenerator("Transposed", GenerateTransposed)
}
