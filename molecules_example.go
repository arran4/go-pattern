package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var MoleculesOutputFilename = "molecules.png"

const MoleculesBaseLabel = "Molecules"

// Molecules Example (formerly Stones)
// Demonstrates using Worley Noise to create an atomic/molecular structure.
func ExampleNewMolecules() {
	// Base Worley Noise (F1) provides the cellular structure
	noise := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(42),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)

	// ColorMap:
	// Center (distance 0) -> Light
	// Edge (distance ~0.5) -> Dark
	// Gaps -> Black

	molecules := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{180, 180, 190, 255}}, // Center
		ColorStop{Position: 0.4, Color: color.RGBA{100, 100, 110, 255}}, // Edge
		ColorStop{Position: 0.45, Color: color.RGBA{50, 50, 55, 255}},   // Darker edge
		ColorStop{Position: 0.5, Color: color.RGBA{10, 10, 10, 255}},    // Gap
		ColorStop{Position: 1.0, Color: color.RGBA{0, 0, 0, 255}},       // Deep gap
	)

	f, err := os.Create(MoleculesOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, molecules); err != nil {
		panic(err)
	}
}

func GenerateMolecules(b image.Rectangle) image.Image {
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
	RegisterGenerator(MoleculesBaseLabel, GenerateMolecules)
}
