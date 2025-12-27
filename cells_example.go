package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var CellsOutputFilename = "cells.png"

const CellsBaseLabel = "Cells"

// Cells Example (Biological)
// Demonstrates using Worley Noise to create a biological cell structure (e.g., plant cells).
func ExampleNewCells() {
	// F1 Euclidean gives distance to center of cell.
	// We want irregular organic cells.
	noise := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(777),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyJitter(0.8), // High jitter for organic look
	)

	// ColorMap:
	// 0.0 - 0.2: Nucleus (Dark Green)
	// 0.2 - 0.25: Nucleus Membrane (Lighter)
	// 0.25 - 0.7: Cytoplasm (Light Green, Translucent look)
	// 0.7 - 0.9: Cell Wall Inner (Darker Green)
	// 0.9 - 1.0: Cell Wall (Thick Dark Border)

	cells := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{20, 80, 20, 255}},    // Nucleus Center
		ColorStop{Position: 0.18, Color: color.RGBA{40, 100, 40, 255}},  // Nucleus
		ColorStop{Position: 0.20, Color: color.RGBA{100, 180, 100, 255}},// Membrane
		ColorStop{Position: 0.25, Color: color.RGBA{150, 220, 150, 255}},// Cytoplasm Start
		ColorStop{Position: 0.70, Color: color.RGBA{140, 210, 140, 255}},// Cytoplasm End
		ColorStop{Position: 0.85, Color: color.RGBA{50, 120, 50, 255}},  // Wall Inner
		ColorStop{Position: 0.95, Color: color.RGBA{10, 40, 10, 255}},   // Wall Outer
		ColorStop{Position: 1.0, Color: color.RGBA{0, 20, 0, 255}},      // Gap
	)

	f, err := os.Create(CellsOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, cells); err != nil {
		panic(err)
	}
}

func GenerateCells(b image.Rectangle) image.Image {
	noise := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.04),
		SetSeed(777),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyJitter(0.8),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{20, 80, 20, 255}},
		ColorStop{Position: 0.18, Color: color.RGBA{40, 100, 40, 255}},
		ColorStop{Position: 0.20, Color: color.RGBA{100, 180, 100, 255}},
		ColorStop{Position: 0.25, Color: color.RGBA{150, 220, 150, 255}},
		ColorStop{Position: 0.70, Color: color.RGBA{140, 210, 140, 255}},
		ColorStop{Position: 0.85, Color: color.RGBA{50, 120, 50, 255}},
		ColorStop{Position: 0.95, Color: color.RGBA{10, 40, 10, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{0, 20, 0, 255}},
	)
}

func init() {
	RegisterGenerator(CellsBaseLabel, GenerateCells)
}
