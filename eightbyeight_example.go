package pattern

import (
	"image"
	"image/png"
	"os"
)

var EightByEightOutputFilename = "eightbyeight.png"
var EightByEightZoomLevels = []int{4}

func ExampleNewEightByEight() {
	i := NewEightByEight(1)
	f, err := os.Create(EightByEightOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err := png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateEightByEight(b image.Rectangle) image.Image {
	return NewEightByEight(1, SetBounds(b))
}

func init() {
	RegisterGenerator("EightByEight", GenerateEightByEight)
}
