package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func ExampleNewCrossHatch() {
	i := NewCrossHatch(
		SetAngles(45, -45),
		SetLineSize(2),
		SetSpaceSize(8),
		SetLineColor(color.RGBA{0, 0, 0, 255}),
		SetSpaceColor(color.RGBA{255, 255, 255, 255}),
	)
	f, err := os.Create(CrossHatchOutputFilename)
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
