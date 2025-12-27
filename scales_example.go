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
// Demonstrates using Worley Noise (Manhattan distance) to create scales.
func ExampleNewScales() {
	// Manhattan distance creates diamond/square shapes resembling scales.
	noise := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(99),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricManhattan),
	)

	// Map distance to scale colors.
	// Center (low distance) -> Highlight
	// Edge (high distance) -> Shadow/Border
	scales := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{100, 200, 100, 255}}, // Center (Green)
		ColorStop{Position: 0.6, Color: color.RGBA{50, 150, 50, 255}},   // Body
		ColorStop{Position: 0.9, Color: color.RGBA{20, 80, 20, 255}},    // Edge
		ColorStop{Position: 1.0, Color: color.RGBA{10, 40, 10, 255}},    // Border
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
		SetFrequency(0.04),
		SetSeed(99),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricManhattan),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{100, 200, 100, 255}},
		ColorStop{Position: 0.6, Color: color.RGBA{50, 150, 50, 255}},
		ColorStop{Position: 0.9, Color: color.RGBA{20, 80, 20, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{10, 40, 10, 255}},
	)
}

func init() {
	RegisterGenerator(ScalesBaseLabel, GenerateScales)
}
