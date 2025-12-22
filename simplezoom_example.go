package pattern

import (
	"image/color"
	"image/png"
	"os"
)

var SimpleZoomOutputFilename = "simplezoom.png"

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
