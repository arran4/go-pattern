package pattern

import (
	"image/color"
	"image/png"
	"os"
)

var CheckerOutputFilename = "checker.png"

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

var CheckerZoomLevels = []int{2, 4}

const CheckerOrder = 1
