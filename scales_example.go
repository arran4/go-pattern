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
// Demonstrates using the Scales pattern to create Amazonian fish scales.
func ExampleNewScales() {
	// Use the explicit Scales pattern for proper overlapping geometry.
	// Radius 40, SpacingX 40 (touching horizontally), SpacingY 20 (half-overlap vertically).
	pattern := NewScales(
		SetScaleRadius(40),
		SetScaleXSpacing(40),
		SetScaleYSpacing(25),
	)

	// The Scales pattern returns a heightmap (0 edge, 1 center).
	// We want to map this to look like a tough fish scale.
	// Center: Shiny/Metallic
	// Gradient towards edge.
	// Edge: Dark border.

	scales := NewColorMap(pattern,
		ColorStop{Position: 0.0, Color: color.RGBA{10, 10, 10, 255}},    // Deep edge (overlap shadow)
		ColorStop{Position: 0.2, Color: color.RGBA{40, 40, 30, 255}},    // Rim
		ColorStop{Position: 0.5, Color: color.RGBA{100, 100, 80, 255}},  // Body
		ColorStop{Position: 0.8, Color: color.RGBA{160, 150, 120, 255}}, // Highlight start
		ColorStop{Position: 1.0, Color: color.RGBA{200, 190, 160, 255}}, // Peak Highlight
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
	pattern := NewScales(
		SetBounds(b),
		SetScaleRadius(40),
		SetScaleXSpacing(40),
		SetScaleYSpacing(25),
	)
	return NewColorMap(pattern,
		ColorStop{Position: 0.0, Color: color.RGBA{10, 10, 10, 255}},
		ColorStop{Position: 0.2, Color: color.RGBA{40, 40, 30, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{100, 100, 80, 255}},
		ColorStop{Position: 0.8, Color: color.RGBA{160, 150, 120, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{200, 190, 160, 255}},
	)
}

func init() {
	RegisterGenerator(ScalesBaseLabel, GenerateScales)
}
