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

// ExampleNewWood demonstrates a procedural wood texture using ConcentricRings, Warp, and Noise.
func ExampleNewWood() image.Image {
	// 1. Create a wood color palette (gradient from light to dark brown)
	woodLight := color.RGBA{210, 180, 140, 255} // Tan
	woodDark := color.RGBA{139, 69, 19, 255}    // SaddleBrown

	colors := []color.Color{}
	steps := 12
	// Gradient cycle
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		colors = append(colors, lerpColor(woodLight, woodDark, t))
	}
	for i := steps - 1; i >= 0; i-- {
		t := float64(i) / float64(steps-1)
		colors = append(colors, lerpColor(woodLight, woodDark, t))
	}

	// 2. Base Rings: "Plank" style
	// Center far to the left to create vertical arcs.
	rings := NewConcentricRings(colors,
		SetCenter(-300, 128),
		SetFrequency(0.12),
	)

	// 3. Grain Distortion: Stretched noise for vertical grain.
	// We generate a small height noise and scale it up vertically.
	grainNoiseBase := NewNoise(NoiseSeed(42), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.1, // High frequency in the base (horizontal detail)
		Octaves: 2,
	}))
	// Set bounds for the base noise to be wide but short (stretched vertically when upscaled)
	if n, ok := grainNoiseBase.(interface{ SetBounds(image.Rectangle) }); ok {
		n.SetBounds(image.Rect(0, 0, 256, 32))
	}

	// Stretch it vertically by scaling to 256x256
	grainDistortion := NewScale(grainNoiseBase, ScaleToSize(256, 256))

	warpedRings := NewWarp(rings,
		WarpDistortion(grainDistortion),
		WarpScale(12.0), // Distortion magnitude
	)

	// 4. Pores: Small dark dashes.
	// Use Worley Noise, generated at short height and scaled up to elongate dots into dashes.
	poreBase := NewWorleyNoise(
		SetFrequency(0.2), // Higher frequency
		NoiseSeed(101),
	)
	if n, ok := poreBase.(interface{ SetBounds(image.Rectangle) }); ok {
		n.SetBounds(image.Rect(0, 0, 256, 64)) // 4x stretch later
	}

	poreStretched := NewScale(poreBase, ScaleToSize(256, 256))

	// Worley returns distance (0 at center, 1 at edge).
	// We want dark dashes at centers.
	poreLayer := NewColorMap(poreStretched,
		ColorStop{0.0, color.RGBA{60, 30, 0, 255}}, // Center: Dark pore
		ColorStop{0.25, color.RGBA{60, 30, 0, 255}},
		ColorStop{0.35, color.White},                 // Background: White (Transparent in Multiply)
		ColorStop{1.0, color.White},
	)

	// Blend pores onto wood using Multiply.
	final := NewBlend(warpedRings, poreLayer, BlendMultiply)

	return final
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
