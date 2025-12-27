package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var PebblesOutputFilename = "pebbles.png"

const PebblesBaseLabel = "Pebbles"

// Pebbles Example (Chipped Stone)
// Demonstrates using Worley Noise with Domain Warping (via Perlin Noise) to create chipped stones.
func ExampleNewPebbles() {
	// 1. Create the base Worley noise.
	// We use F1 to get the distance from the center of the cell.
	// Inverting this (or properly mapping it) gives us distinct stones.
	base := NewWorleyNoise(
		SetFrequency(0.04),
		SetSeed(200),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)

	// 2. Create Perlin noise for distortion (chipped edges).
	// High frequency noise for fine details.
	distortionX := NewNoise(
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        300,
			Frequency:   0.1,
			Octaves:     3,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)
	distortionY := NewNoise(
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        400,
			Frequency:   0.1,
			Octaves:     3,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)

	// 3. Warp the Worley noise using the Perlin noise.
	// This makes the smooth cellular boundaries jagged.
	warped := NewWarp(base,
		WarpDistortionX(distortionX),
		WarpDistortionY(distortionY),
		WarpXScale(5.0), // Magnitude of distortion
		WarpYScale(5.0),
		WarpDistortionScale(1.0), // Scale of noise sampling
	)

	// 4. Map to colors.
	// Worley F1 is 0 at center, higher at edges.
	// We map low values to stone color, high values to gap/mortar.
	// Because of warping, the transition will be noisy.
	pebbles := NewColorMap(warped,
		ColorStop{Position: 0.0, Color: color.RGBA{180, 180, 190, 255}}, // Center Highlight
		ColorStop{Position: 0.2, Color: color.RGBA{120, 120, 130, 255}}, // Stone Body
		ColorStop{Position: 0.5, Color: color.RGBA{80, 80, 90, 255}},    // Stone Edge
		ColorStop{Position: 0.6, Color: color.RGBA{40, 35, 30, 255}},    // Gap/Mortar start
		ColorStop{Position: 1.0, Color: color.RGBA{20, 15, 10, 255}},    // Gap Center
	)

	f, err := os.Create(PebblesOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, pebbles); err != nil {
		panic(err)
	}
}

func GeneratePebbles(b image.Rectangle) image.Image {
	base := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.04),
		SetSeed(200),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)

	distortionX := NewNoise(
		SetBounds(b),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        300,
			Frequency:   0.1,
			Octaves:     3,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)
	distortionY := NewNoise(
		SetBounds(b),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        400,
			Frequency:   0.1,
			Octaves:     3,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)

	warped := NewWarp(base,
		WarpDistortionX(distortionX),
		WarpDistortionY(distortionY),
		WarpXScale(5.0),
		WarpYScale(5.0),
	)

	return NewColorMap(warped,
		ColorStop{Position: 0.0, Color: color.RGBA{180, 180, 190, 255}},
		ColorStop{Position: 0.2, Color: color.RGBA{120, 120, 130, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{80, 80, 90, 255}},
		ColorStop{Position: 0.6, Color: color.RGBA{40, 35, 30, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{20, 15, 10, 255}},
	)
}

func init() {
	RegisterGenerator(PebblesBaseLabel, GeneratePebbles)
}
