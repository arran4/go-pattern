package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var ScalesOutputFilename = "scales.png"

const ScalesBaseLabel = "Scales"

// Scales Example
// Demonstrates using Worley Noise to create reptile/dragon scales.
func ExampleNewScales() {
	// Reptile scales are often irregular polygons (Voronoi) but closely packed.
	// F1 Euclidean gives distance to center.
	// F2-F1 gives distance to edge.
	// Let's use F1 for the "dome" shape of the scale.
	noise := NewWorleyNoise(
		SetFrequency(0.025),
		SetSeed(303), // Different seed for variety
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyJitter(1.0), // Fully organic/random
	)

	// ColorMap:
	// Create a "convex" look with lighting.
	// Center (0.0): Highlight (Shiny)
	// Body: Scale color (e.g., Red/Orange for dragon)
	// Edge (~0.6-0.8): Darker (Shadow/Curve)
	// Gap (0.9-1.0): Black/Dark Brown (Skin between scales)

	scales := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{255, 200, 150, 255}}, // Specular Highlight
		ColorStop{Position: 0.2, Color: color.RGBA{200, 80, 20, 255}},   // Main Body (Orange)
		ColorStop{Position: 0.5, Color: color.RGBA{160, 40, 10, 255}},   // Darker Body
		ColorStop{Position: 0.7, Color: color.RGBA{100, 20, 5, 255}},    // Shadow/Curve
		ColorStop{Position: 0.85, Color: color.RGBA{40, 10, 0, 255}},    // Deep Shadow
		ColorStop{Position: 1.0, Color: color.RGBA{10, 5, 0, 255}},      // Gap
	)

	f, err := os.Create(ScalesOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, scales); err != nil {
		panic(err)
	}
}

func GenerateScales(b image.Rectangle) image.Image {
	noise := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.05),
		SetSeed(303),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyJitter(1.0),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{255, 200, 150, 255}},
		ColorStop{Position: 0.2, Color: color.RGBA{200, 80, 20, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{160, 40, 10, 255}},
		ColorStop{Position: 0.7, Color: color.RGBA{100, 20, 5, 255}},
		ColorStop{Position: 0.85, Color: color.RGBA{40, 10, 0, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{10, 5, 0, 255}},
	)
}

func init() {
	RegisterGenerator(ScalesBaseLabel, GenerateScales)
}
