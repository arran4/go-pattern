package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var CellsOutputFilename = "cells.png"

const CellsBaseLabel = "Cells"

// Cells Example
// Demonstrates using Worley Noise (CellID) to color cells randomly.
func ExampleNewCells() {
	// CellID output gives a random grayscale value per cell
	noise := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(777),
		SetWorleyOutput(OutputCellID),
	)

	// Use ColorMap to map random IDs to a palette (e.g., biological cells)
	cells := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{255, 100, 100, 255}}, // Red
		ColorStop{Position: 0.2, Color: color.RGBA{255, 150, 150, 255}},
		ColorStop{Position: 0.4, Color: color.RGBA{200, 50, 50, 255}},
		ColorStop{Position: 0.6, Color: color.RGBA{220, 80, 80, 255}},
		ColorStop{Position: 0.8, Color: color.RGBA{180, 20, 20, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{255, 120, 120, 255}},
	)

	// Can combine with F1 to add borders or gradients inside cells?
	// For now, just the solid colored cells.

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
		SetWorleyOutput(OutputCellID),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{255, 100, 100, 255}},
		ColorStop{Position: 0.2, Color: color.RGBA{255, 150, 150, 255}},
		ColorStop{Position: 0.4, Color: color.RGBA{200, 50, 50, 255}},
		ColorStop{Position: 0.6, Color: color.RGBA{220, 80, 80, 255}},
		ColorStop{Position: 0.8, Color: color.RGBA{180, 20, 20, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{255, 120, 120, 255}},
	)
}

func init() {
	RegisterGenerator(CellsBaseLabel, GenerateCells)
}
