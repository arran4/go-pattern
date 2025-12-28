package pattern

import (
	"image"
	"image/color"
)

// NewAsphalt creates a procedural asphalt texture.
// It uses noise for grain and Voronoi for potential cracks or aggregate.
func NewAsphalt() image.Image {
	// 1. Base Grain (High frequency noise)
	// Simulates the bitumen and small stones.
	grain := NewNoise(
		NoiseSeed(100),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.8, Octaves: 4, Persistence: 0.6}),
	)

	// Map to dark greys
	base := NewColorMap(grain,
		ColorStop{0.0, color.RGBA{30, 30, 30, 255}},
		ColorStop{0.5, color.RGBA{50, 50, 50, 255}},
		ColorStop{1.0, color.RGBA{70, 70, 70, 255}},
	)

	// 2. Larger Aggregate (Specks)
	// We can use a Scatter pattern or high threshold noise.
	specksNoise := NewNoise(
		NoiseSeed(200),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 2.0}),
	)
	// Threshold to get sparse dots
	specks := NewColorMap(specksNoise,
		ColorStop{0.0, color.RGBA{0, 0, 0, 0}},   // Transparent
		ColorStop{0.7, color.RGBA{0, 0, 0, 0}},
		ColorStop{0.75, color.RGBA{180, 180, 180, 255}}, // Light stones
		ColorStop{1.0, color.RGBA{200, 200, 200, 255}},
	)

	// 3. Tar/Dark patches (Low freq)
	patchesNoise := NewNoise(
		NoiseSeed(300),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.05}),
	)
	patches := NewColorMap(patchesNoise,
		ColorStop{0.0, color.RGBA{0, 0, 0, 50}}, // Darken
		ColorStop{0.5, color.RGBA{0, 0, 0, 0}},  // No change
		ColorStop{1.0, color.RGBA{255, 255, 255, 10}}, // Slight lighten
	)

	// Composite
	// Base + Specks + Patches

	// Blend specks over base
	layer1 := NewBlend(base, specks, BlendNormal)

	// Blend patches (Overlay or Multiply)
	// Simple alpha blending handles the patches logic defined above.
	final := NewBlend(layer1, patches, BlendNormal)

	return final
}
