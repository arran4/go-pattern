package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var StonesOutputFilename = "stones.png"

const StonesBaseLabel = "Stones"

// Stones Example
// Demonstrates using Worley Noise to create a stone texture.
func ExampleNewStones() {
	// Base Worley Noise (F1) provides the cellular structure (stones)
	noise := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(42),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)

	// Invert the noise to make stones white and edges dark (optional, depending on look)
	// Or use ColorMap to map distance to stone colors.

	// Let's create a ColorMap to define the stone look.
	// Center of stone (distance 0) -> Light Grey
	// Edge of stone (distance ~0.5) -> Dark Grey
	// Gaps -> Black

	stones := NewColorMap(noise,
		// Color ramp
		ColorStop{Position: 0.0, Color: color.RGBA{180, 180, 190, 255}}, // Center
		ColorStop{Position: 0.4, Color: color.RGBA{100, 100, 110, 255}}, // Edge of stone
		ColorStop{Position: 0.45, Color: color.RGBA{50, 50, 55, 255}},   // Darker edge
		ColorStop{Position: 0.5, Color: color.RGBA{10, 10, 10, 255}},    // Gap
		ColorStop{Position: 1.0, Color: color.RGBA{0, 0, 0, 255}},       // Deep gap
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
		SetFrequency(0.04), // Higher freq for demo box
		SetSeed(42),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{180, 180, 190, 255}},
		ColorStop{Position: 0.4, Color: color.RGBA{100, 100, 110, 255}},
		ColorStop{Position: 0.45, Color: color.RGBA{50, 50, 55, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{10, 10, 10, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{0, 0, 0, 255}},
	)
}

func init() {
	RegisterGenerator(StonesBaseLabel, GenerateStones)
}
