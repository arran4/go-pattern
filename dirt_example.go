package pattern

import (
	"image"
	"image/color"
)

var (
	DirtOutputFilename = "dirt.png"
	Dirt_mudOutputFilename = "dirt_mud.png"
)

func ExampleNewDirt() image.Image {
	// 1. Base Dirt: Brown, grainy noise
	base := NewNoise(
		NoiseSeed(101),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        101,
			Frequency:   0.1,
			Octaves:     4,
			Persistence: 0.6,
		}),
	)

	dirtColor := NewColorMap(base,
		ColorStop{Position: 0.0, Color: color.RGBA{40, 30, 20, 255}}, // Dark Brown
		ColorStop{Position: 0.5, Color: color.RGBA{80, 60, 40, 255}}, // Brown
		ColorStop{Position: 0.8, Color: color.RGBA{100, 80, 60, 255}}, // Light Brown
		ColorStop{Position: 1.0, Color: color.RGBA{120, 100, 80, 255}}, // Pebbles
	)

	// 2. Grain: High freq noise overlay
	grain := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.5}))
	detailed := NewBlend(dirtColor, grain, BlendOverlay)

	return detailed
}

func ExampleNewDirt_mud() image.Image {
	// Base dirt
	dirt := ExampleNewDirt()

	// 3. Wetness Mask: Puddles
	// Low frequency noise thresholded
	puddleNoise := NewNoise(
		NoiseSeed(202),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:      202,
			Frequency: 0.02,
		}),
	)

	// Mask: White where puddles are, Black where dirt is
	// Threshold at 0.6
	mask := NewColorMap(puddleNoise,
		ColorStop{Position: 0.0, Color: color.Black},
		ColorStop{Position: 0.55, Color: color.Black},
		ColorStop{Position: 0.6, Color: color.White},
		ColorStop{Position: 1.0, Color: color.White},
	)

	// Puddles: Darker, smoother, reflective (mocked by color)
	// Or use NormalMap to make them flat vs rough dirt.
	// Let's make puddles dark brown/black and subtract detail.

	puddleColor := NewRect(SetFillColor(color.RGBA{20, 15, 10, 255}))

	// Blend puddle color based on mask?
	// We don't have a "BlendMask" pattern yet that takes a mask image.
	// But we can use boolean ops or just Blend?
	// Or we can use the mask as alpha for the puddle layer and overlay it.
	// But our patterns usually return opaque images unless alpha is handled.

	// Let's assume we want to composite Puddle over Dirt using Mask.
	// This usually requires a MaskedComposite pattern.
	// I don't see one.

	// Workaround:
	// 1. Create Puddle Layer (Dark)
	// 2. Create Dirt Layer
	// 3. Blend them? No, we want distinct areas.
	// If I use `NewBlend` with a mode? No standard mode does masking.

	// I can use `NewBoolean` (BitwiseAnd) if mask is binary?
	// Dirt AND (NOT Mask) + Puddle AND Mask.

	// Invert mask for dirt
	invMask := NewBitwiseNot(mask)

	dirtPart := NewBitwiseAnd([]image.Image{dirt, invMask})
	puddlePart := NewBitwiseAnd([]image.Image{puddleColor, mask})

	return NewBitwiseOr([]image.Image{dirtPart, puddlePart})
}

func GenerateDirt(rect image.Rectangle) image.Image {
	return ExampleNewDirt()
}

func GenerateDirt_mud(rect image.Rectangle) image.Image {
	return ExampleNewDirt_mud()
}

func init() {
	GlobalGenerators["Dirt"] = GenerateDirt
	GlobalGenerators["Dirt_mud"] = GenerateDirt_mud
}
