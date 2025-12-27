package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var GrassOutputFilename = "grass.png"

const GrassBaseLabel = "Grass"

// Grass Example
// Demonstrates using Perlin Noise with ColorMap to create a simple grass texture.
func ExampleNewGrass() {
	// 1. Create a base noise layer for general color variation.
	baseNoise := NewNoise(
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        500,
			Frequency:   0.02,
			Octaves:     4,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)

	// 2. Create a high-frequency noise layer for "blades" or detail.
	detailNoise := NewNoise(
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        600,
			Frequency:   0.2, // High frequency for grass blades
			Octaves:     2,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)

	// 3. Blend them. We want the detail to be prominent but influenced by the base.
	// Multiply might darken too much, let's use Overlay or just simple addition/average.
	// Actually, let's just use the detail noise warped by base noise for a wind-blown look?
	// Or simply blend them.

	// Let's try blending: Base * 0.5 + Detail * 0.5
	// Using BlendAverage is simple.
	blended := NewBlend(baseNoise, detailNoise, BlendAverage)

	// 4. Map to Grass Colors.
	grass := NewColorMap(blended,
		ColorStop{Position: 0.0, Color: color.RGBA{10, 40, 10, 255}},    // Deep shadow/dirt
		ColorStop{Position: 0.3, Color: color.RGBA{30, 80, 30, 255}},    // Dark Grass
		ColorStop{Position: 0.6, Color: color.RGBA{60, 140, 40, 255}},   // Mid Grass
		ColorStop{Position: 0.8, Color: color.RGBA{100, 180, 60, 255}},  // Light Grass
		ColorStop{Position: 1.0, Color: color.RGBA{140, 220, 100, 255}}, // Tips/Highlights
	)

	f, err := os.Create(GrassOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, grass); err != nil {
		panic(err)
	}
}

func GenerateGrass(b image.Rectangle) image.Image {
	baseNoise := NewNoise(
		SetBounds(b),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        500,
			Frequency:   0.02,
			Octaves:     4,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)

	detailNoise := NewNoise(
		SetBounds(b),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        600,
			Frequency:   0.2,
			Octaves:     2,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}),
	)

	blended := NewBlend(baseNoise, detailNoise, BlendAverage)

	return NewColorMap(blended,
		ColorStop{Position: 0.0, Color: color.RGBA{10, 40, 10, 255}},
		ColorStop{Position: 0.3, Color: color.RGBA{30, 80, 30, 255}},
		ColorStop{Position: 0.6, Color: color.RGBA{60, 140, 40, 255}},
		ColorStop{Position: 0.8, Color: color.RGBA{100, 180, 60, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{140, 220, 100, 255}},
	)
}

func init() {
	RegisterGenerator(GrassBaseLabel, GenerateGrass)
}
