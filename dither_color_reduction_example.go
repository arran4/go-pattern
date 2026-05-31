package pattern

import (
	"image"
	"image/color"
)

var DitherColorReductionOutputFilename = "dither_color_reduction.png"
var DitherColorReductionZoomLevels = []int{}
const DitherColorReductionOrder = 104

// ExampleNewDitherColorReduction demonstrates color reduction capabilities using various palettes.
func ExampleNewDitherColorReduction() image.Image {
	return NewBayer8x8Dither(NewGopher(), PaletteCGA)
}

func GenerateDitherColorReduction(b image.Rectangle) image.Image {
	return NewBayer8x8Dither(NewGopher(), PaletteCGA)
}

// GenerateDitherColorReductionReferences generates examples of color reduction.
func GenerateDitherColorReductionReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	gopher := NewGopher()

	// Palettes

	// CGA Palette 1 High Intensity
	cga := PaletteCGA

	// Gameboy (4 shades of green)
	gameboy := color.Palette{
		color.RGBA{0x0f, 0x38, 0x0f, 0xff}, // Darkest
		color.RGBA{0x30, 0x62, 0x30, 0xff},
		color.RGBA{0x8b, 0xac, 0x0f, 0xff},
		color.RGBA{0x9b, 0xbc, 0x0f, 0xff}, // Lightest
	}

	// 1-bit Black and White
	bw := color.Palette{color.Black, color.White}

	// WebSafe (216 colors) - simplified subset or full generation?
	// Let's generate a simple WebSafe-ish 6x6x6 palette
	webSafe := make(color.Palette, 0, 216)
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				webSafe = append(webSafe, color.RGBA{uint8(r * 51), uint8(g * 51), uint8(b * 51), 0xff})
			}
		}
	}

	return map[string]func(image.Rectangle) image.Image{
		"CGA_Bayer8x8": func(b image.Rectangle) image.Image {
			return NewBayer8x8Dither(gopher, cga)
		},
		"CGA_Yliluoma1": func(b image.Rectangle) image.Image {
			return NewYliluoma1Dither(gopher, cga, 8)
		},
		"Gameboy_Bayer4x4": func(b image.Rectangle) image.Image {
			return NewBayer4x4Dither(gopher, gameboy)
		},
		"Gameboy_Knoll": func(b image.Rectangle) image.Image {
			return NewKnollDither(gopher, gameboy, 8)
		},
		"BW_Halftone": func(b image.Rectangle) image.Image {
			return NewHalftoneDither(gopher, 6, bw)
		},
		"WebSafe_Bayer8x8": func(b image.Rectangle) image.Image {
			return NewBayer8x8Dither(gopher, webSafe)
		},
		"WebSafe_Yliluoma2": func(b image.Rectangle) image.Image {
			return NewYliluoma2Dither(gopher, webSafe, 8)
		},
	}, []string{
		"CGA_Bayer8x8",
		"CGA_Yliluoma1",
		"Gameboy_Bayer4x4",
		"Gameboy_Knoll",
		"BW_Halftone",
		"WebSafe_Bayer8x8",
		"WebSafe_Yliluoma2",
	}
}

// PaletteCGA is the CGA Palette 1 High Intensity
var PaletteCGA = color.Palette{
	color.RGBA{0x00, 0x00, 0x00, 0xff}, // Black
	color.RGBA{0xff, 0x55, 0xff, 0xff}, // Magenta
	color.RGBA{0x55, 0xff, 0xff, 0xff}, // Cyan
	color.RGBA{0xff, 0xff, 0xff, 0xff}, // White
}

func init() {
	RegisterGenerator("DitherColorReduction", GenerateDitherColorReduction)
	RegisterReferences("DitherColorReduction", GenerateDitherColorReductionReferences)
}
