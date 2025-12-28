package pattern

import (
	"image"
	"image/color"
)

// NewAsphalt creates a procedural asphalt concrete texture (Micro view).
// It combines fine grain noise, larger aggregate stones, and variations.
func NewAsphalt() image.Image {
	// 1. Fine Grain (Bitumen/Sand)
	// High frequency noise
	grain := NewNoise(
		NoiseSeed(101),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 1.5, Octaves: 3, Persistence: 0.7}),
	)

	// Map grain to dark grey/black asphalt colors
	base := NewColorMap(grain,
		ColorStop{0.0, color.RGBA{20, 20, 20, 255}},
		ColorStop{1.0, color.RGBA{60, 60, 60, 255}},
	)

	// 2. Aggregate (Stones)
	// Medium/High frequency noise, thresholded to create "spots"
	stonesNoise := NewNoise(
		NoiseSeed(202),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 2.5, Octaves: 1}),
	)

	// White/Grey stones, sparse
	stones := NewColorMap(stonesNoise,
		ColorStop{0.0, color.RGBA{0, 0, 0, 0}},   // Transparent
		ColorStop{0.65, color.RGBA{0, 0, 0, 0}},
		ColorStop{0.70, color.RGBA{100, 100, 100, 255}}, // Dark stone
		ColorStop{0.85, color.RGBA{180, 180, 180, 255}}, // Light stone
		ColorStop{1.0, color.RGBA{200, 200, 200, 255}},
	)

	// 3. Surface Variation (Patches/Wear)
	// Low frequency noise
	wearNoise := NewNoise(
		NoiseSeed(303),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.1}),
	)

	wear := NewColorMap(wearNoise,
		ColorStop{0.0, color.RGBA{0, 0, 0, 40}}, // Darker patches (oil/tar)
		ColorStop{0.5, color.RGBA{0, 0, 0, 0}},
		ColorStop{1.0, color.RGBA{255, 255, 255, 10}}, // Lighter patches (dry)
	)

	// Composite
	// Base + Stones + Wear

	// Add stones to base
	step1 := NewBlend(base, stones, BlendNormal)

	// Apply wear
	step2 := NewBlend(step1, wear, BlendNormal)

	return step2
}
