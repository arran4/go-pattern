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
// Demonstrates composing patterns using Blend to create a realistic island terrain.
func ExampleNewIslands() {
	// Layer 1: Base Shape (Worley F1 Euclidean) - Large distinct landmasses
	baseShape := NewWorleyNoise(
		SetFrequency(0.01),
		SetSeed(555),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)

	// Layer 2: Detail (Perlin Noise) - Adds coastline complexity and terrain roughness
	detail := NewNoise(
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        123,
			Frequency:   0.05,
			Octaves:     4,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)

	// Blend: Subtract detail from base shape? Or Overlay?
	// Worley F1 is 0 at center (Peak), 1 at edge (Deep Water).
	// We want Peaks to be high (1.0). So let's Invert Worley first?
	// Or just use ColorMap on the result.
	// If we Add detail to Worley, the values increase.
	// Let's use BlendOverlay to mix the gradients.

	mixed := NewBlend(baseShape, detail, BlendOverlay)

	// ColorMap:
	// Worley: 0 (Peak) -> 1 (Edge)
	// Overlay tends to push contrast.
	// Let's define:
	// 0.0 - 0.2: Snow (Peak)
	// 0.2 - 0.4: Mountain/Rock
	// 0.4 - 0.5: Forest
	// 0.5 - 0.6: Sand
	// 0.6 - 1.0: Water

	islands := NewColorMap(mixed,
		ColorStop{Position: 0.0, Color: color.RGBA{250, 250, 250, 255}}, // Snow
		ColorStop{Position: 0.15, Color: color.RGBA{120, 120, 120, 255}}, // Rock
		ColorStop{Position: 0.30, Color: color.RGBA{34, 139, 34, 255}},  // Forest
		ColorStop{Position: 0.50, Color: color.RGBA{210, 180, 140, 255}}, // Sand
		ColorStop{Position: 0.55, Color: color.RGBA{64, 164, 223, 255}}, // Water
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
	baseShape := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.02),
		SetSeed(555),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)

	detail := NewNoise(
		SetBounds(b),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        123,
			Frequency:   0.1,
			Octaves:     4,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)

	mixed := NewBlend(baseShape, detail, BlendOverlay)

	return NewColorMap(mixed,
		ColorStop{Position: 0.0, Color: color.RGBA{250, 250, 250, 255}},
		ColorStop{Position: 0.15, Color: color.RGBA{120, 120, 120, 255}},
		ColorStop{Position: 0.30, Color: color.RGBA{34, 139, 34, 255}},
		ColorStop{Position: 0.50, Color: color.RGBA{210, 180, 140, 255}},
		ColorStop{Position: 0.55, Color: color.RGBA{64, 164, 223, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{0, 0, 128, 255}},
	)
}

func init() {
	RegisterGenerator(IslandsBaseLabel, GenerateIslands)
}
