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
// Demonstrates using Worley Noise (Euclidean) to create organic scales (fish/reptile).
func ExampleNewScales() {
	// Euclidean distance creates circular/hexagonal cells (honeycomb).
	// Using F1 we get the distance from center.
	noise := NewWorleyNoise(
		SetFrequency(0.025),
		SetSeed(202),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyJitter(0.6), // Slightly regular, but organic
	)

	// Map distance to scale look.
	// Center (0.0) -> Highlight (Top of convexity)
	// Edge (~0.5-0.6) -> Dark (Overlap/Shadow)

	scales := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{100, 220, 100, 255}}, // Center Highlight
		ColorStop{Position: 0.4, Color: color.RGBA{50, 160, 50, 255}},   // Body
		ColorStop{Position: 0.55, Color: color.RGBA{20, 80, 20, 255}},   // Shadow
		ColorStop{Position: 0.60, Color: color.RGBA{10, 30, 10, 255}},   // Edge
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
		SetSeed(202),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyJitter(0.6),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{100, 220, 100, 255}},
		ColorStop{Position: 0.4, Color: color.RGBA{50, 160, 50, 255}},
		ColorStop{Position: 0.55, Color: color.RGBA{20, 80, 20, 255}},
		ColorStop{Position: 0.60, Color: color.RGBA{10, 30, 10, 255}},
	)
}

func init() {
	RegisterGenerator(ScalesBaseLabel, GenerateScales)
}
