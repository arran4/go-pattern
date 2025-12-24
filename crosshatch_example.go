package pattern

import (
	"image"
	"image/color"
)

func ExampleNewCrossHatch() {
	// This function body is empty because the bootstrap tool uses the function signature
	// and the following variable to generate the documentation and image.
}

var (
	// CrossHatchZoomLevels defines the zoom levels for the CrossHatch pattern documentation.
	CrossHatchZoomLevels = []int{}

	// CrossHatchOutputFilename defines the output filename for the CrossHatch pattern image.
	CrossHatchOutputFilename = "crosshatch.png"

	// CrossHatchBaseLabel defines the base label for the CrossHatch pattern.
	CrossHatchBaseLabel = "CrossHatch"
)

func init() {
	GlobalGenerators[CrossHatchBaseLabel] = GenerateCrossHatch
	GlobalReferences[CrossHatchBaseLabel] = GenerateCrossHatchReferences
}

func GenerateCrossHatch(rect image.Rectangle) image.Image {
	return NewCrossHatch(
		SetAngles(45, -45),
		SetLineSize(2),
		SetSpaceSize(8),
		SetLineColor(color.RGBA{0, 0, 0, 255}),
		SetSpaceColor(color.RGBA{255, 255, 255, 255}),
	)
}

func GenerateCrossHatchReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Single Hatch": func(rect image.Rectangle) image.Image {
			return NewCrossHatch(
				SetAngles(45),
				SetLineSize(2),
				SetSpaceSize(8),
				SetLineColor(color.RGBA{0, 0, 0, 255}),
				SetSpaceColor(color.RGBA{255, 255, 255, 255}),
			)
		},
		"Triple Hatch": func(rect image.Rectangle) image.Image {
			return NewCrossHatch(
				SetAngles(0, 60, 120),
				SetLineSize(1),
				SetSpaceSize(10),
				SetLineColor(color.RGBA{0, 0, 0, 255}),
				SetSpaceColor(color.RGBA{255, 255, 255, 255}),
			)
		},
	}, []string{"Single Hatch", "Triple Hatch"}
}
