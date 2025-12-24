package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var Bayer2x2DitherOutputFilename = "bayer2x2.png"
var Bayer2x2DitherZoomLevels = []int{}

const Bayer2x2DitherOrder = 31
const Bayer2x2DitherBaseLabel = "Bayer2x2Dither"

// Bayer2x2Dither Pattern
// Example of applying a 2x2 Bayer ordered dither.
func ExampleNewBayer2x2Dither() {
	// Black and White Palette
	palette := []color.Color{color.Black, color.White}
	i := NewBayer2x2Dither(NewGopher(), palette)

	f, err := os.Create(Bayer2x2DitherOutputFilename)
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

func GenerateBayer2x2Dither(b image.Rectangle) image.Image {
	palette := []color.Color{color.Black, color.White}
	return NewBayer2x2Dither(NewGopher(), palette)
}

func GenerateBayer2x2DitherReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	// Grayscale palette (4 levels)
	grayPalette := []color.Color{
		color.Black,
		color.Gray{Y: 85},
		color.Gray{Y: 170},
		color.White,
	}

	return map[string]func(image.Rectangle) image.Image{
		"Bayer2x2 (B&W)": func(b image.Rectangle) image.Image {
			return NewBayer2x2Dither(NewGopher(), []color.Color{color.Black, color.White})
		},
		"Bayer2x2 (4 Grays)": func(b image.Rectangle) image.Image {
			return NewBayer2x2Dither(NewGopher(), grayPalette)
		},
	}, []string{"Bayer2x2 (B&W)", "Bayer2x2 (4 Grays)"}
}

func init() {
	RegisterGenerator("Bayer2x2Dither", GenerateBayer2x2Dither)
	RegisterReferences("Bayer2x2Dither", GenerateBayer2x2DitherReferences)
}
