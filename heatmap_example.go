package pattern

import (
	"image"
	"image/color"
	"math"
)

var HeatmapOutputFilename = "heatmap.png"
var HeatmapZoomLevels = []int{}

const HeatmapOrder = 24

// Heatmap
// Generates a heatmap for the function z = sin(x) * cos(y).
func ExampleNewHeatmap() {
	// See GenerateHeatmap for implementation details
}

func GenerateHeatmap(b image.Rectangle) image.Image {
	f := func(x, y float64) float64 {
		return math.Sin(x) * math.Cos(y)
	}

	return NewHeatmap(f,
		SetBounds(b),
		SetXRange(-math.Pi, math.Pi),
		SetYRange(-math.Pi, math.Pi),
		SetZRange(-1.0, 1.0),
		SetStartColor(color.RGBA{0, 0, 255, 255}), // Blue
		SetEndColor(color.RGBA{255, 0, 0, 255}),   // Red
	)
}

func init() {
	RegisterGenerator("Heatmap", GenerateHeatmap)
}
