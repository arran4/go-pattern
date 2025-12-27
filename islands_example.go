package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var IslandsOutputFilename = "islands.png"

const IslandsBaseLabel = "Islands"

// Islands Example
// Demonstrates using Worley Noise as a heightmap for islands/biomes.
func ExampleNewIslands() {
	// F1 Euclidean gives us cone-like shapes growing from the points.
	// Inverting it (or just mapping correctly) gives islands.
	// Distance 0 (at point) = Peak (Mountain) or Center (Deep Water)?
	// Let's assume points are peaks of islands. So Distance 0 = High, Distance 1 = Low.
	// But Worley returns distance increasing from center.
	// So 0.0 = Peak, 0.5 = Coast, 1.0 = Deep Ocean.

	noise := NewWorleyNoise(
		SetFrequency(0.015),
		SetSeed(555),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)

	// ColorMap:
	// 0.0 - 0.2: Snow/Mountain (Top of the cone/cell center)
	// 0.2 - 0.4: Green/Forest
	// 0.4 - 0.5: Sand/Beach
	// 0.5 - 0.6: Shallow Water
	// 0.6 - 1.0: Deep Water

	islands := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{255, 255, 255, 255}}, // Snow
		ColorStop{Position: 0.15, Color: color.RGBA{100, 100, 100, 255}}, // Rock
		ColorStop{Position: 0.20, Color: color.RGBA{34, 139, 34, 255}},  // Forest
		ColorStop{Position: 0.40, Color: color.RGBA{210, 180, 140, 255}}, // Sand
		ColorStop{Position: 0.45, Color: color.RGBA{64, 164, 223, 255}}, // Water
		ColorStop{Position: 1.0, Color: color.RGBA{0, 0, 128, 255}},     // Deep Water
	)

	f, err := os.Create(IslandsOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, islands); err != nil {
		panic(err)
	}
}

func GenerateIslands(b image.Rectangle) image.Image {
	noise := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.03), // Zoom out a bit for demo box
		SetSeed(555),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{255, 255, 255, 255}},
		ColorStop{Position: 0.15, Color: color.RGBA{100, 100, 100, 255}},
		ColorStop{Position: 0.20, Color: color.RGBA{34, 139, 34, 255}},
		ColorStop{Position: 0.40, Color: color.RGBA{210, 180, 140, 255}},
		ColorStop{Position: 0.45, Color: color.RGBA{64, 164, 223, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{0, 0, 128, 255}},
	)
}

func init() {
	RegisterGenerator(IslandsBaseLabel, GenerateIslands)
}
