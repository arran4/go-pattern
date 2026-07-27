package pattern

import (
	"image"
	"image/color"
)

// Wood Example

var (
	WoodOutputFilename = "wood.png"
	WoodZoomLevels     = []int{} // No zoom needed for texture
	WoodBaseLabel      = "Wood"
)

// ExampleNewWood demonstrates a procedural wood texture using domain warping on a distance field.
func ExampleNewWood() image.Image {
	// 1. Wood Palette
	// Dark brown (Late wood / Rings) -> Light Tan (Early wood) -> Dark
	woodPalette := []ColorStop{
		{0.0, color.RGBA{101, 67, 33, 255}},   // Dark Brown (Ring Edge)
		{0.15, color.RGBA{160, 120, 80, 255}}, // Transition
		{0.5, color.RGBA{222, 184, 135, 255}}, // Light Tan (Center - Burlywood)
		{0.85, color.RGBA{160, 120, 80, 255}}, // Transition
		{1.0, color.RGBA{101, 67, 33, 255}},   // Back to Edge
	}

	// 2. Base "Heightmap" Generator
	// We create a grayscale gradient for rings (0-255).
	grayScale := make([]color.Color, 256)
	for i := range grayScale {
		grayScale[i] = color.Gray{Y: uint8(i)}
	}

	// Use ConcentricRings to generate the base distance field.
	// We want ~10 rings across the 256px width.
	// 256 colors in palette.
	// To get 1 cycle every 25 pixels: Freq = 256/25 ≈ 10.
	// To get elongated vertical rings, FreqY should be lower (slower change).
	ringsBase := NewConcentricRings(grayScale,
		SetCenter(128, -100), // Off-center top
		SetFrequencyX(8.0),   // ~30px width per ring
		SetFrequencyY(0.8),   // Stretched vertically (10x elongation)
	)

	// 3. Main Distortion (Growth Wobble)
	// Low frequency noise to warp the rings.
	// Noise values are 0..1 (from NewNoise/Perlin).
	// Warp maps intensity to offset.
	wobbleNoise := NewNoise(NoiseSeed(101), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.015,
		Octaves:   2,
	}))

	// Apply warp.
	// Scale 20.0 means max offset is +/- 20 pixels.
	// Since rings are ~30px wide, this distorts them significantly but keeps structure.
	warpedRings := NewWarp(ringsBase,
		WarpDistortion(wobbleNoise),
		WarpScale(20.0),
	)

	// 4. Fiber Grain (Fine Detail)
	// Add "Turbulence" to the warp using higher frequency noise.
	// This simulates the jagged edges of the grain.
	fiberDistortion := NewNoise(NoiseSeed(303), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.1, // Higher freq
		Octaves:   3,   // More detail
	}))

	// Chain Warps: WarpedRings -> Warp again with fiber distortion
	doubleWarped := NewWarp(warpedRings,
		WarpDistortion(fiberDistortion),
		WarpScale(2.0), // Small jaggedness (2 pixels)
	)

	// 5. Color Mapping
	// Map the grayscale intensity (warped distance) to the wood palette.
	finalWood := NewColorMap(doubleWarped, woodPalette...)

	return finalWood
}

func init() {
	GlobalGenerators[WoodBaseLabel] = GenerateWood
	GlobalReferences[WoodBaseLabel] = GenerateWoodReferences
}

func GenerateWood(rect image.Rectangle) image.Image {
	return ExampleNewWood()
}

func GenerateWoodReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}
