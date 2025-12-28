package pattern

import (
	"image"
	"image/color"
)

var (
	RoadOutputFilename = "road.png"
	Road_markedOutputFilename = "road_marked.png"
)

func ExampleNewRoad() image.Image {
	// Asphalt: Aggregate noise (grey with black/white speckles)
	base := NewNoise(
		NoiseSeed(707),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.8}),
	)

	asphalt := NewColorMap(base,
		ColorStop{Position: 0.0, Color: color.RGBA{40, 40, 40, 255}},
		ColorStop{Position: 0.2, Color: color.RGBA{60, 60, 60, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{50, 50, 50, 255}},
		ColorStop{Position: 0.8, Color: color.RGBA{70, 70, 70, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{90, 90, 90, 255}},
	)

	// Cracks: Voronoi edges or Lightning-like noise
	// Let's use Voronoi edges.
	// We don't have direct edge access in Voronoi pattern, but we can do difference of Voronoi outputs?
	// Or use `EdgeDetect` on a Voronoi pattern.

	// Voronoi with colors
	v2 := NewVoronoi(
		[]image.Point{
			{10, 10}, {50, 200}, {200, 50}, {220, 220},
			{100, 100}, {150, 150}, {80, 20}, {20, 80},
		},
		[]color.Color{color.Black, color.White},
	)

	edges := NewEdgeDetect(v2)
	// Edges are white on black.

	// Invert to get Black cracks on White background
	cracks := NewBitwiseNot(edges)

	// Multiply cracks onto asphalt
	return NewBlend(asphalt, cracks, BlendMultiply)
}

func ExampleNewRoad_marked() image.Image {
	road := ExampleNewRoad()

	// Painted lines
	// Yellow center line
	lines := NewVerticalLine(
		SetLineSize(10),
		SetSpaceSize(255), // Only one line in middle?
		SetLineColor(color.RGBA{255, 200, 0, 255}),
		SetSpaceColor(color.Transparent),
	)

	// Shift line to center?
	// VerticalLine starts at x=0.
	// We can use Padding or Translation? No translate pattern.
	// But `VerticalLine` repeats.
	// If we want it centered, we need to adjust phase or sizes.

	// Or use Rect with bounds?
	// Let's use a Rect for the line.
	// We don't have easy positioning for Rect.

	// Let's stick to VerticalLine, it will appear at left.
	// We can offset it if we had offset.
	// Maybe just accept it repeats.

	// Composite lines over road
	return NewBlend(road, lines, BlendOverlay) // Overlay blends it. Normal would be better if we had Alpha composite.
}

func GenerateRoad(rect image.Rectangle) image.Image {
	return ExampleNewRoad()
}

func GenerateRoad_marked(rect image.Rectangle) image.Image {
	return ExampleNewRoad_marked()
}

func init() {
	GlobalGenerators["Road"] = GenerateRoad
	GlobalGenerators["Road_marked"] = GenerateRoad_marked
}
