package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var HexGridOutputFilename = "hex_grid.png"

const HexGridBaseLabel = "HexGrid"

// HexGrid example: alternating palette across axial coordinates with a subtle bevel.
func ExampleNewHexGrid() {
	img := GenerateHexGrid(image.Rect(0, 0, 255, 255))
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
	palette := color.Palette{
		color.NRGBA{R: 58, G: 90, B: 101, A: 255},
		color.NRGBA{R: 173, G: 216, B: 230, A: 255},
	}
	return NewHexGrid(
		SetBounds(b),
		SetRadius(28),
		SetHexPalette(palette),
		SetHexBevelDepth(9),
	)
}

func init() {
	RegisterGenerator(HexGridBaseLabel, GenerateHexGrid)
}
