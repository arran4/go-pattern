package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var CheckerOutputFilename = "checker.png"
var CheckerZoomLevels = []int{2, 4}

const CheckerOrder = 1

// Checker Pattern
// Alternates between two colors in a checkerboard fashion.
func ExampleNewChecker() {
	i := NewChecker(color.Black, color.White)
	f, err := os.Create(CheckerOutputFilename)
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

func GenerateChecker(b image.Rectangle) image.Image {
	return NewDemoChecker(SetBounds(b))
}

func init() {
	RegisterGenerator("Checker", GenerateChecker)
}
