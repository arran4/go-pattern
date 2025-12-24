package pattern

import (
	"image"
	"image/color"
)

func ExampleNewEdgeDetect() {
	// Generate the main demo image
	// This function is detected by the bootstrap tool
}

func GenerateEdgeDetect(args ...any) image.Image {
	return NewDemoEdgeDetect()
}

func GenerateEdgeDetectReferences() (map[string]func() image.Image, []string) {
	// We want to show the original image vs the edge detected one

	// Recreate the source used in NewDemoEdgeDetect
	sourceGen := func() image.Image {
		chk := NewChecker(color.Black, color.White)
		return NewSimpleZoom(chk, 20)
	}

	return map[string]func() image.Image{
		"Source": sourceGen,
	}, []string{"Source"}
}
