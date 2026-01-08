package pattern

import (
	"image"
	"image/png"
	"os"
)

var SubpixelLinesOutputFilename = "subpixel_lines.png"
var SubpixelLinesZoomLevels = []int{}

const SubpixelLinesOrder = 41

// Subpixel lines with per-channel offset and vignette.
func ExampleNewSubpixelLines() {
	i := NewSubpixelLines(
		SetLineThickness(2),
		SetOffsetStrength(0.65),
		SetVignetteRadius(0.82),
	)
	f, err := os.Create(SubpixelLinesOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, i)
	if err != nil {
		panic(err)
	}
}

func GenerateSubpixelLines(b image.Rectangle) image.Image {
	return NewSubpixelLines(
		SetBounds(b),
		SetLineThickness(3),
		SetOffsetStrength(0.8),
		SetVignetteRadius(0.85),
	)
}

func init() {
	RegisterGenerator("SubpixelLines", GenerateSubpixelLines)
}
