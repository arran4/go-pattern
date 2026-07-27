package pattern

import (
	"image"
	"image/color"
)

var (
	PaintedPlanksOutputFilename = "painted_planks.png"
	PaintedPlanksZoomLevels     = []int{}
	PaintedPlanksBaseLabel      = "PaintedPlanks"
)

// ExampleNewPaintedPlanks demonstrates segmented planks with grain noise per board
// and a chipped paint overlay.
func ExampleNewPaintedPlanks() image.Image {
	return NewPaintedPlanks(
		SetPlankBaseWidth(72),
		SetPlankWidthVariance(0.32),
		SetGrainIntensity(0.75),
		SetPaintWear(0.42),
		SetPaintColor(color.RGBA{177, 202, 214, 255}),
	)
}

func init() {
	GlobalGenerators[PaintedPlanksBaseLabel] = GeneratePaintedPlanks
	GlobalReferences[PaintedPlanksBaseLabel] = GeneratePaintedPlanksReferences
}

func GeneratePaintedPlanks(rect image.Rectangle) image.Image {
	return NewPaintedPlanks(SetBounds(rect))
}

func GeneratePaintedPlanksReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}
