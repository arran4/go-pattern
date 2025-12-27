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

var WindowsDither4x4OutputFilename = "dither_windows_4x4.png"
var WindowsDither4x4ZoomLevels = []int{}

// ExampleNewWindowsDither4x4 demonstrates 4x4 variant.
func ExampleNewWindowsDither4x4() image.Image {
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

// Helper to use Windows16 palette
func init() {
	// Register Windows16 palette if we had a registry, but we don't need to for examples.
}
