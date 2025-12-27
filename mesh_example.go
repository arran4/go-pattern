package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var MeshOutputFilename = "mesh.png"

const MeshBaseLabel = "Mesh"

// Mesh Example (formerly Cells)
// Demonstrates using Worley Noise (CellID) to create a 3D polygon mesh surface look.
func ExampleNewMesh() {
	// CellID output gives a random grayscale value per cell
	noise := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(777),
		SetWorleyOutput(OutputCellID),
	)

	// Use ColorMap to map random IDs to a palette (e.g., biological cells or tech mesh)
	mesh := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{255, 100, 100, 255}}, // Red
		ColorStop{Position: 0.2, Color: color.RGBA{255, 150, 150, 255}},
		ColorStop{Position: 0.4, Color: color.RGBA{200, 50, 50, 255}},
		ColorStop{Position: 0.6, Color: color.RGBA{220, 80, 80, 255}},
		ColorStop{Position: 0.8, Color: color.RGBA{180, 20, 20, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{255, 120, 120, 255}},
	)

	f, err := os.Create(MeshOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, mesh); err != nil {
		panic(err)
	}
}

func GenerateMesh(b image.Rectangle) image.Image {
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
	RegisterGenerator(MeshBaseLabel, GenerateMesh)
}
