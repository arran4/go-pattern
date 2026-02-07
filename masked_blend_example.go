package pattern

import (
	"image"
	"image/color"
)

var MaskedBlendOutputFilename = "masked_blend.png"

func ExampleNewMaskedBlend() image.Image {
	// Background: A dark blue plasma
	bg := NewPlasma(
		SetStartColor(color.RGBA{0, 0, 50, 255}),
		SetEndColor(color.RGBA{0, 0, 150, 255}),
	)

	// Foreground: A bright orange noise
	fgNoise := NewNoise(
		NoiseSeed(42),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.1}),
	)
	fg := NewColorMap(fgNoise,
		ColorStop{Position: 0.0, Color: color.RGBA{255, 100, 0, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{255, 200, 0, 255}},
	)

	// Mask: A white circle on black background
	mask := NewCircle(
		SetRadius(80),
		SetFillColor(color.White),
		SetSpaceColor(color.Black),
	)

	return NewMaskedBlend(bg, fg, mask)
}

func GenerateMaskedBlend(rect image.Rectangle) image.Image {
	return ExampleNewMaskedBlend()
}

func init() {
	RegisterGenerator("MaskedBlend", GenerateMaskedBlend)
}
