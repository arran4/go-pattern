package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var StonesOutputFilename = "stones.png"

const StonesBaseLabel = "Stones"

// Stones Example (Riverbed / Cobblestones)
// Demonstrates using Worley Noise (F2-F1) to create smooth stones with mortar.
func ExampleNewStones() {
	// F2-F1 gives distance to the border.
	// Border is 0. Center is High.
	noise := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(100),
		SetWorleyOutput(OutputF2MinusF1),
		SetWorleyMetric(MetricEuclidean),
	)

	// Map:
	// 0.0 - 0.1: Mortar (Dark)
	// 0.1 - 0.3: Edge of stone (Darker Grey)
	// 0.3 - 1.0: Stone Body (Grey/Blueish with gradient)

	stones := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{20, 15, 10, 255}},    // Mortar
		ColorStop{Position: 0.15, Color: color.RGBA{40, 40, 45, 255}},   // Stone Edge
		ColorStop{Position: 0.3, Color: color.RGBA{80, 80, 90, 255}},    // Stone Body
		ColorStop{Position: 0.8, Color: color.RGBA{150, 150, 160, 255}}, // Highlight
	)

	f, err := os.Create(StonesOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, stones); err != nil {
		panic(err)
	}
}

func GenerateStones(b image.Rectangle) image.Image {
	noise := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.04),
		SetSeed(100),
		SetWorleyOutput(OutputF2MinusF1),
		SetWorleyMetric(MetricEuclidean),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{20, 15, 10, 255}},
		ColorStop{Position: 0.15, Color: color.RGBA{40, 40, 45, 255}},
		ColorStop{Position: 0.3, Color: color.RGBA{80, 80, 90, 255}},
		ColorStop{Position: 0.8, Color: color.RGBA{150, 150, 160, 255}},
	)
}

func init() {
	RegisterGenerator(StonesBaseLabel, GenerateStones)
}
