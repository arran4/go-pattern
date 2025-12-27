package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var PebblesOutputFilename = "pebbles.png"

const PebblesBaseLabel = "Pebbles"

// Pebbles Example (Chipped Stone / Gravel)
// Demonstrates using Worley Noise combined with Perlin Noise (via Blend) to create irregular, chipped stones.
func ExampleNewPebbles() {
	// 1. Create the base Worley noise.
	// F2-F1 is often used for cells, but F1 gives us the "distance from center" which allows us to control the stone shape/size.
	base := NewWorleyNoise(
		SetFrequency(0.05),
		SetSeed(200),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyJitter(1.0), // Full jitter for organic placement
	)

	// 2. Create Perlin noise for "chips" and surface texture.
	noise := NewNoise(
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        300,
			Frequency:   0.15, // Higher frequency for detail
			Octaves:     3,
			Persistence: 0.6,
			Lacunarity:  2.0,
		}),
	)

	// 3. Scale down the noise intensity.
	// We only want the noise to slightly perturb the Worley distance field.
	// Mapping 0-255 to roughly 0-50 range.
	scaledNoise := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{0, 0, 0, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{60, 60, 60, 255}},
	)

	// 4. Combine Worley Base + Scaled Noise.
	// The noise adds to the distance, effectively bringing the "edge" threshold closer in random spots (chipping).
	blended := NewBlend(base, scaledNoise, BlendAdd)

	// 5. Map to Stone Colors.
	// The blended value represents "Distance from center + Noise".
	// Low values = Center of stone (High).
	// Medium values = Edge of stone (Sloping down).
	// High values = Gap/Mortar.
	pebbles := NewColorMap(blended,
		ColorStop{Position: 0.0, Color: color.RGBA{160, 160, 165, 255}}, // Highlight
		ColorStop{Position: 0.2, Color: color.RGBA{120, 120, 125, 255}}, // Body
		ColorStop{Position: 0.45, Color: color.RGBA{80, 80, 85, 255}},   // Darker Body
		ColorStop{Position: 0.5, Color: color.RGBA{50, 50, 55, 255}},    // Edge/Rim (Sharp transition)
		ColorStop{Position: 0.52, Color: color.RGBA{20, 15, 10, 255}},   // Gap Start
		ColorStop{Position: 1.0, Color: color.RGBA{10, 5, 0, 255}},      // Gap Deep
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
		SetFrequency(0.05),
		SetSeed(200),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyJitter(1.0),
	)

	noise := NewNoise(
		SetBounds(b),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        300,
			Frequency:   0.15,
			Octaves:     3,
			Persistence: 0.6,
			Lacunarity:  2.0,
		}),
	)

	scaledNoise := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{0, 0, 0, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{60, 60, 60, 255}},
	)

	blended := NewBlend(base, scaledNoise, BlendAdd)

	return NewColorMap(blended,
		ColorStop{Position: 0.0, Color: color.RGBA{160, 160, 165, 255}},
		ColorStop{Position: 0.2, Color: color.RGBA{120, 120, 125, 255}},
		ColorStop{Position: 0.45, Color: color.RGBA{80, 80, 85, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{50, 50, 55, 255}},
		ColorStop{Position: 0.52, Color: color.RGBA{20, 15, 10, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{10, 5, 0, 255}},
	)
}

func init() {
	RegisterGenerator(PebblesBaseLabel, GeneratePebbles)
}
