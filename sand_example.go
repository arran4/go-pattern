package pattern

import (
	"image"
	"image/color"
)

var (
	SandOutputFilename = "sand.png"
	Sand_dunesOutputFilename = "sand_dunes.png"
	Sand_zoomedOutputFilename = "sand_zoomed.png"
)

func ExampleNewSand() image.Image {
	// 1. Fine grain noise
	grain := NewNoise(
		NoiseSeed(303),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        303,
			Frequency:   0.5,
			Octaves:     2,
		}),
	)

	sandColor := NewColorMap(grain,
		ColorStop{Position: 0.0, Color: color.RGBA{194, 178, 128, 255}}, // Sand
		ColorStop{Position: 1.0, Color: color.RGBA{225, 205, 150, 255}}, // Light Sand
	)

	return sandColor
}

func ExampleNewSand_zoomed() image.Image {
	// Zoomed in sand to show grains
	// Use Scatter or just low freq noise mapped to dots?
	// Let's use noise with thresholding to make "grains".

	noise := NewNoise(
		NoiseSeed(304),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.2}),
	)

	// Map to distinct grains
	grains := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{160, 140, 100, 255}}, // Dark grain
		ColorStop{Position: 0.4, Color: color.RGBA{210, 190, 150, 255}}, // Main sand
		ColorStop{Position: 0.7, Color: color.RGBA{230, 210, 170, 255}}, // Light grain
		ColorStop{Position: 0.9, Color: color.RGBA{255, 255, 255, 255}}, // Quartz sparkle
	)

	return grains
}

func ExampleNewSand_dunes() image.Image {
	// Base sand
	sand := ExampleNewSand()

	// 2. Ripples
	ripples := NewNoise(
		NoiseSeed(404),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:      404,
			Frequency: 0.05,
		}),
	)
	// Stretch to make lines
	ripplesStretched := NewScale(ripples, ScaleX(1.0), ScaleY(10.0))

	// Darken troughs using Multiply
	shadows := NewColorMap(ripplesStretched,
		ColorStop{Position: 0.0, Color: color.RGBA{180, 180, 180, 255}}, // Darker (Grey for Multiply)
		ColorStop{Position: 0.5, Color: color.White}, // No change
		ColorStop{Position: 1.0, Color: color.White}, // No change
	)

	// Rotate ripples
	rotatedShadows := NewRotate(shadows, 90)

	return NewBlend(sand, rotatedShadows, BlendMultiply)
}

func GenerateSand(rect image.Rectangle) image.Image {
	return ExampleNewSand()
}

func GenerateSand_dunes(rect image.Rectangle) image.Image {
	return ExampleNewSand_dunes()
}

func GenerateSand_zoomed(rect image.Rectangle) image.Image {
	return ExampleNewSand_zoomed()
}

func init() {
	GlobalGenerators["Sand"] = GenerateSand
	GlobalGenerators["Sand_dunes"] = GenerateSand_dunes
	GlobalGenerators["Sand_zoomed"] = GenerateSand_zoomed
}
