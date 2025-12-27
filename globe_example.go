package pattern

import (
	"image"
	"image/color"
)

var (
	GlobeOutputFilename = "globe.png"
	Globe_wireframeOutputFilename = "globe_wireframe.png"
	Globe_texturedOutputFilename = "globe_textured.png"
	Globe_terrain_maskedOutputFilename = "globe_terrain_masked.png"
)

// ExampleNewGlobe generates a globe pattern.
func ExampleNewGlobe() image.Image {
	return NewGlobe(
		SetLatitudeLines(5),
		SetLongitudeLines(12),
		SetLineSize(2),
		SetLineColor(color.RGBA{200, 200, 200, 255}),
		SetFillColor(color.RGBA{0, 0, 50, 255}),
		SetSpaceColor(color.Transparent),
		SetAngle(15),
		SetTilt(23.5),
	)
}

// ExampleNewGlobe_wireframe generates a wireframe globe demonstrating transparency and back-face visibility.
func ExampleNewGlobe_wireframe() image.Image {
	return NewGlobe(
		SetLatitudeLines(10),
		SetLongitudeLines(20),
		SetLineSize(1),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
		SetAngle(15),
		SetTilt(23.5),
		// No FillColor set -> Transparent
	)
}

// ExampleNewGlobe_textured demonstrates UV mapping a texture onto the globe.
func ExampleNewGlobe_textured() image.Image {
	return NewGlobe(
		SetLatitudeLines(8),
		SetLongitudeLines(16),
		SetLineSize(1),
		SetLineColor(color.White),
		SetFillImageSource(NewWarp(
			NewColorMap(NewNoise(
				SetNoiseAlgorithm(&PerlinNoise{Seed: 42, Frequency: 0.1})),
				ColorStop{0.0, color.RGBA{0, 0, 100, 255}},
				ColorStop{0.5, color.RGBA{0, 100, 0, 255}},
				ColorStop{1.0, color.White},
			),
			WarpXScale(0.05),
			WarpYScale(0.05),
			WarpDistortion(NewNoise(
				SetNoiseAlgorithm(&PerlinNoise{Seed: 1}),
			)),
		)),
		SetAngle(45),
		SetTilt(10),
	)
}

// ExampleNewGlobe_terrain_masked demonstrates a "basic globe" using 2D masking of the terrain pattern.
func ExampleNewGlobe_terrain_masked() image.Image {
	// Create terrain texture using Warp + Perlin + ColorMap (similar to Warp_terrain example)
	terrain := NewWarp(
		NewColorMap(NewNoise(SetNoiseAlgorithm(&PerlinNoise{Seed: 42, Frequency: 0.02})),
			ColorStop{0.0, color.RGBA{0, 0, 50, 255}},   // Deep Water
			ColorStop{0.4, color.RGBA{0, 100, 200, 255}}, // Water
			ColorStop{0.42, color.RGBA{200, 200, 150, 255}}, // Sand
			ColorStop{0.5, color.RGBA{0, 150, 0, 255}},   // Grass
			ColorStop{0.7, color.RGBA{100, 100, 100, 255}}, // Rock
			ColorStop{0.9, color.White},                  // Snow
		),
		WarpXScale(0.02),
		WarpYScale(0.02),
		WarpDistortion(NewNoise(SetNoiseAlgorithm(&PerlinNoise{Seed: 42}))),
	)

	// Use Circle to mask it
	return NewCircle(
		SetFillImageSource(terrain),
		SetLineSize(2),
		SetLineColor(color.RGBA{50, 50, 50, 255}), // Atmosphere/Border
		SetSpaceColor(color.Transparent),
	)
}
