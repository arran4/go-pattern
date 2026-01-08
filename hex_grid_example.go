package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var HexGridOutputFilename = "hex_grid.png"

// HexGrid example: alternating palette across axial coordinates with a subtle bevel.
func ExampleNewHexGrid() {
	palette := color.Palette{
		color.NRGBA{R: 58, G: 90, B: 101, A: 255},
		color.NRGBA{R: 173, G: 216, B: 230, A: 255},
	}
	img := NewHexGrid(
		SetRadius(28),
		SetHexPalette(palette),
		SetHexBevelDepth(9),
	)

	f, err := os.Create(HexGridOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func GenerateHexGrid(b image.Rectangle) image.Image {
	return NewHexGrid(
		SetBounds(b),
		SetRadius(26),
		SetHexBevelDepth(8),
	)
}

func init() {
	RegisterGenerator("HexGridMask", GenerateHexGrid)
}
