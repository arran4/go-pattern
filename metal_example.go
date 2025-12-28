package pattern

import (
	"image"
	"image/color"
)

var (
	MetalOutputFilename = "metal.png"
	Metal_scratchedOutputFilename = "metal_scratched.png"
)

func ExampleNewMetal() image.Image {
	// Brushed Metal
	// 1. High frequency noise
	noise := NewNoise(
		NoiseSeed(333),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        333,
			Frequency:   0.1, // Lower frequency to avoid aliasing when scaled
			Octaves:     3,
			Persistence: 0.5,
		}),
	)

	// Map to grey gradients before scaling
	metalBase := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{50, 50, 50, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{150, 150, 150, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{200, 200, 200, 255}},
	)

	// 2. Anisotropy: Scale heavily
	// Scale X large, Y small? Or X 1, Y large?
	// Vertical streaks: Scale Y > 1.
	// But `NewScale` interpolates.
	// Try Scale X=1, Y=10.
	brushed := NewScale(metalBase, ScaleX(1.0), ScaleY(10.0))

	return brushed
}

func ExampleNewMetal_scratched() image.Image {
	// Base brushed metal
	base := ExampleNewMetal()

	// Scratches using CrossHatch
	hatchMultiply := NewCrossHatch(
		SetLineColor(color.Gray{100}), // Dark scratches
		SetSpaceColor(color.White),    // No change
		SetLineSize(1),
		SetSpaceSize(40),
		SetAngles(10, 80, 170),
	)

	// Distort scratches slightly
	distort := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.1}))
	hatchWarped := NewWarp(hatchMultiply, WarpDistortion(distort), WarpScale(2.0))

	return NewBlend(base, hatchWarped, BlendMultiply)
}

func GenerateMetal(rect image.Rectangle) image.Image {
	return ExampleNewMetal()
}

func GenerateMetal_scratched(rect image.Rectangle) image.Image {
	return ExampleNewMetal_scratched()
}

func init() {
	GlobalGenerators["Metal"] = GenerateMetal
	GlobalGenerators["Metal_scratched"] = GenerateMetal_scratched
}
