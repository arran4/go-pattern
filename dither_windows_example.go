package pattern

import (
	"image"
)

var WindowsDitherOutputFilename = "dither_windows.png"
var WindowsDitherZoomLevels = []int{}

// ExampleNewWindowsDither demonstrates the classic Windows 16-color dithering
// using standard Bayer ordered dithering (comparable to what the user requested).
// This uses a 4x4 matrix which was common, or 8x8.
// The user linked article discusses standard ordered dithering with Bayer matrix.
func ExampleNewWindowsDither() image.Image {
	img := NewGopher()
	// Spread 0 = auto calculate, or we can fine tune.
	// Standard Windows dithering often just used the nearest color after thresholding.
	// We use NewBayer8x8Dither for "Standard Ordered Dithering".
	return NewBayer8x8Dither(img, Windows16)
}

func GenerateWindowsDither(b image.Rectangle) image.Image {
	img := NewGopher()
	return NewBayer8x8Dither(img, Windows16)
}

var WindowsDither4x4OutputFilename = "dither_windows_4x4.png"
var WindowsDither4x4ZoomLevels = []int{}

// ExampleNewWindowsDither4x4 demonstrates 4x4 variant.
func ExampleNewWindowsDither4x4() image.Image {
	img := NewGopher()
	return NewBayer4x4Dither(img, Windows16)
}

func GenerateWindowsDither4x4(b image.Rectangle) image.Image {
	img := NewGopher()
	return NewBayer4x4Dither(img, Windows16)
}

var WindowsDitherHalftoneOutputFilename = "dither_windows_halftone.png"
var WindowsDitherHalftoneZoomLevels = []int{}

// ExampleNewWindowsDitherHalftone uses a halftone pattern.
func ExampleNewWindowsDitherHalftone() image.Image {
	img := NewGopher()
	return NewHalftoneDither(img, 8, Windows16)
}

func GenerateWindowsDitherHalftone(b image.Rectangle) image.Image {
	img := NewGopher()
	return NewHalftoneDither(img, 8, Windows16)
}

func GenerateWindowsDitherReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	gopher := NewGopher()
	return map[string]func(image.Rectangle) image.Image{
			"Bayer2x2": func(b image.Rectangle) image.Image {
				return NewBayer2x2Dither(gopher, Windows16)
			},
			"Bayer4x4": func(b image.Rectangle) image.Image {
				return NewBayer4x4Dither(gopher, Windows16)
			},
			"Bayer8x8": func(b image.Rectangle) image.Image {
				return NewBayer8x8Dither(gopher, Windows16)
			},
			"Halftone4": func(b image.Rectangle) image.Image {
				return NewHalftoneDither(gopher, 4, Windows16)
			},
			"Halftone8": func(b image.Rectangle) image.Image {
				return NewHalftoneDither(gopher, 8, Windows16)
			},
			"Random": func(b image.Rectangle) image.Image {
				return NewRandomDither(gopher, Windows16, 12345)
			},
			"Yliluoma1": func(b image.Rectangle) image.Image {
				return NewYliluoma1Dither(gopher, Windows16, 8)
			},
			"Yliluoma2": func(b image.Rectangle) image.Image {
				return NewYliluoma2Dither(gopher, Windows16, 8)
			},
			"Knoll": func(b image.Rectangle) image.Image {
				return NewKnollDither(gopher, Windows16, 8)
			},
			"FloydSteinberg": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, FloydSteinberg, Windows16)
			},
			"SierraLite": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, SierraLite, Windows16)
			},
		}, []string{
			"Bayer2x2", "Bayer4x4", "Bayer8x8",
			"Halftone4", "Halftone8",
			"Random",
			"Yliluoma1", "Yliluoma2", "Knoll",
			"FloydSteinberg", "SierraLite",
		}
}

// Helper to use Windows16 palette
func init() {
	RegisterGenerator("WindowsDither", GenerateWindowsDither)
	RegisterReferences("WindowsDither", GenerateWindowsDitherReferences)
	RegisterGenerator("WindowsDither4x4", GenerateWindowsDither4x4)
	RegisterGenerator("WindowsDitherHalftone", GenerateWindowsDitherHalftone)
}
