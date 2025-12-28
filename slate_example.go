package pattern

import (
	"image"
	"image/color"
)

var (
	SlateOutputFilename = "slate.png"
)

func ExampleNewSlate() image.Image {
	// Slate: Layered noise, laminar structure.
	// Dark grey, slight blue/green tint.

	base := NewNoise(
		NoiseSeed(808),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:      808,
			Frequency: 0.05,
			Octaves:   5,
		}),
	)

	// Map to slate colors
	slateColor := NewColorMap(base,
		ColorStop{Position: 0.0, Color: color.RGBA{40, 45, 50, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{60, 65, 70, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{80, 85, 90, 255}},
	)

	// Laminar effect: Scale Y slightly to stretch horizontally? Or vertically?
	// Slate cleaves. Usually fine layers.
	laminar := NewScale(slateColor, ScaleX(2.0), ScaleY(1.0))

	// Surface bumpiness
	bump := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.2}))

	return NewBlend(laminar, bump, BlendOverlay)
}

func GenerateSlate(rect image.Rectangle) image.Image {
	return ExampleNewSlate()
}

func init() {
	GlobalGenerators["Slate"] = GenerateSlate
}
