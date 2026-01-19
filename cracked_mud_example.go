package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var CrackedMudOutputFilename = "cracked_mud.png"

const CrackedMudBaseLabel = "CrackedMud"

// Cracked Mud Example
// Demonstrates using Worley Noise (F2-F1) to create cracked earth.
func ExampleNewCrackedMud() {
	// F2-F1 gives thick lines at cell boundaries (where distance to 1st and 2nd closest points are similar)
	noise := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(123),
		SetWorleyOutput(OutputF2MinusF1),
		SetWorleyMetric(MetricEuclidean),
	)

	// Map distance to mud colors.
	// Low value (close to 0) means F1 ~= F2, i.e., boundary/crack.
	// High value means center of cell.

	mud := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{30, 20, 10, 255}},   // Crack (Dark brown/black)
		ColorStop{Position: 0.1, Color: color.RGBA{60, 40, 20, 255}},   // Crack edge
		ColorStop{Position: 0.2, Color: color.RGBA{130, 100, 70, 255}}, // Mud surface
		ColorStop{Position: 1.0, Color: color.RGBA{160, 120, 80, 255}}, // Center of mud chunk
	)

	f, err := os.Create(CrackedMudOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, mud); err != nil {
		panic(err)
	}
}

func GenerateCrackedMud(b image.Rectangle) image.Image {
	noise := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.04),
		SetSeed(123),
		SetWorleyOutput(OutputF2MinusF1),
		SetWorleyMetric(MetricEuclidean),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{30, 20, 10, 255}},
		ColorStop{Position: 0.1, Color: color.RGBA{60, 40, 20, 255}},
		ColorStop{Position: 0.2, Color: color.RGBA{130, 100, 70, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{160, 120, 80, 255}},
	)
}

func init() {
	RegisterGenerator(CrackedMudBaseLabel, GenerateCrackedMud)
}
