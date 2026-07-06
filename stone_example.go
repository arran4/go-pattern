package pattern

import (
	"image"
	"image/color"
)

var (
	StoneOutputFilename        = "stone.png"
	Stone_cobbleOutputFilename = "stone_cobble.png"
)

func ExampleNewStone() image.Image {
	// Voronoi base for cells (cobblestones)
	// We want cells to be somewhat irregular.
	voronoi := NewVoronoi(
		// Points
		[]image.Point{
			{50, 50}, {150, 40}, {230, 60},
			{40, 140}, {130, 130}, {240, 150},
			{60, 230}, {160, 240}, {220, 220},
			{100, 100}, {200, 200}, {30, 30},
			{180, 80}, {80, 180},
		},
		// Colors: Using Greyscale for heightmap initially, or color for texture
		// Let's make a texture.
		[]color.Color{
			color.RGBA{100, 100, 100, 255},
			color.RGBA{120, 115, 110, 255},
			color.RGBA{90, 90, 95, 255},
			color.RGBA{110, 110, 110, 255},
			color.RGBA{130, 125, 120, 255},
		},
	)

	// 1. Edge Wear: Distort the Voronoi
	distort := NewNoise(
		NoiseSeed(77),
		SetNoiseAlgorithm(&PerlinNoise{Seed: 77, Frequency: 0.1}),
	)
	worn := NewWarp(voronoi, WarpDistortion(distort), WarpScale(5.0))

	// 2. Surface Detail: Grain
	grain := NewNoise(
		NoiseSeed(88),
		SetNoiseAlgorithm(&PerlinNoise{Seed: 88, Frequency: 0.5}),
	)

	// Blend grain onto stones (Overlay)
	textured := NewBlend(worn, grain, BlendOverlay)

	// We return the Albedo texture, not the normal map.
	// Normal map can be a separate pass or derived.
	return textured
}

func ExampleNewStone_cobble() image.Image {
	// Cellular noise (Worley) for cobblestones heightmap
	worley := NewWorleyNoise(
		SetWorleyMetric(MetricEuclidean),
		SetWorleyOutput(OutputF1), // Distance to closest point
		NoiseSeed(123),
		SetFrequency(0.06),
	)

	// Map Worley (0-1 distance) to Stone Colors
	// Worley: 0 is center, 1 is edge.
	// Cobbles: Center is high/bright, Edge is low/dark (mortar).
	// We want to map the distance to a color gradient.

	cobbleColor := NewColorMap(worley,
		ColorStop{Position: 0.0, Color: color.RGBA{180, 175, 170, 255}}, // Center (Light Stone)
		ColorStop{Position: 0.4, Color: color.RGBA{140, 135, 130, 255}}, // Mid Stone
		ColorStop{Position: 0.7, Color: color.RGBA{100, 95, 90, 255}},   // Dark Stone edge
		ColorStop{Position: 0.85, Color: color.RGBA{60, 55, 50, 255}},   // Mortar start
		ColorStop{Position: 1.0, Color: color.RGBA{40, 35, 30, 255}},    // Deep Mortar
	)

	// Add noise for texture
	noise := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.2}))
	textured := NewBlend(cobbleColor, noise, BlendOverlay)

	return textured
}

func GenerateStone(rect image.Rectangle) image.Image {
	return ExampleNewStone()
}

func GenerateStone_cobble(rect image.Rectangle) image.Image {
	return ExampleNewStone_cobble()
}

func init() {
	GlobalGenerators["Stone"] = GenerateStone
	GlobalGenerators["Stone_cobble"] = GenerateStone_cobble
}
