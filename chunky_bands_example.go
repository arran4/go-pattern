package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var ChunkyBandsOutputFilename = "chunky_bands.png"

// Chunky pixel bands example.
func ExampleChunkyBands() {
	i := NewChunkyBands(
		SetBlockSize(12),
		SetAngle(32),
		SetPalette(
			color.RGBA{20, 24, 44, 255},
			color.RGBA{90, 50, 120, 255},
			color.RGBA{220, 140, 60, 255},
		),
	)
	f, err := os.Create(ChunkyBandsOutputFilename)
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

func GenerateChunkyBands(b image.Rectangle) image.Image {
	return NewDemoChunkyBands(SetBounds(b))
}

func init() {
	RegisterGenerator("ChunkyBands", GenerateChunkyBands)
}
