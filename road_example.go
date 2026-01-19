package pattern

import (
	"image"
	"image/color"
)

var (
	RoadOutputFilename         = "road.png"
	Road_markedOutputFilename  = "road_marked.png"
	Road_terrainOutputFilename = "road_terrain.png"
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

	// Cracks: Voronoi edges
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
	// Yellow center line (dashed?)
	// Let's do a solid double yellow or single yellow.
	// VerticalLine pattern repeats.
	// Image width is usually 255.
	// We want one line in the center.
	// LineSize 10. SpaceSize big enough to push next line off screen.

	lines := NewVerticalLine(
		SetLineSize(8),
		SetSpaceSize(300),
		SetLineColor(color.RGBA{255, 200, 0, 255}), // Paint
		SetSpaceColor(color.Transparent),
		SetPhase(123), // Center: ~127 minus half line width (4) = 123.
	)

	// Composite lines over road using Normal blend (Paint on top)
	return NewBlend(road, lines, BlendNormal)
}

func ExampleNewRoad_terrain() image.Image {
	// Winding road on grass

	// 1. Terrain (Grass)
	grass := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.1}))
	grassColor := NewColorMap(grass,
		ColorStop{0.0, color.RGBA{30, 100, 30, 255}},
		ColorStop{1.0, color.RGBA{50, 150, 50, 255}},
	)

	// 2. Road Mask (Winding curve)
	// We can use a low freq noise thresholded to a thin band?
	// Or use `ModuloStripe` or `Sine` warped.
	// Let's use a warped VerticalLine.

	roadPath := NewVerticalLine(
		SetLineSize(40), // Road width
		SetSpaceSize(300),
		SetLineColor(color.White),  // Mask: White = Road
		SetSpaceColor(color.Black), // Mask: Black = Grass
		SetPhase(105),
	)

	// Warp the road path to make it winding
	warpNoise := NewNoise(NoiseSeed(999), SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.02}))
	windingRoadMask := NewWarp(roadPath, WarpDistortionX(warpNoise), WarpScale(50.0))

	// 3. Road Texture
	// Use asphalt from ExampleNewRoad, but we need to map it to the winding path?
	// A simple tiled asphalt is fine.
	roadTex := ExampleNewRoad()

	// 4. Composite
	// We have Grass (Bg), RoadTex (Fg), Mask (windingRoadMask).
	// We don't have a MaskedBlend.
	// Workaround:
	// GrassPart = Grass * (NOT Mask)
	// RoadPart = RoadTex * Mask
	// Result = GrassPart + RoadPart

	// Invert mask
	invMask := NewBitwiseNot(windingRoadMask)

	// Masking requires BitwiseAnd?
	// But BitwiseAnd operates on colors bits.
	// If Mask is pure Black/White, it works like a stencil for RGB.

	grassPart := NewBitwiseAnd([]image.Image{grassColor, invMask})
	roadPart := NewBitwiseAnd([]image.Image{roadTex, windingRoadMask})

	return NewBitwiseOr([]image.Image{grassPart, roadPart})
}

func GenerateRoad(rect image.Rectangle) image.Image {
	return ExampleNewRoad()
}

func GenerateRoad_marked(rect image.Rectangle) image.Image {
	return ExampleNewRoad_marked()
}

func GenerateRoad_terrain(rect image.Rectangle) image.Image {
	return ExampleNewRoad_terrain()
}

func init() {
	GlobalGenerators["Road"] = GenerateRoad
	GlobalGenerators["Road_marked"] = GenerateRoad_marked
	GlobalGenerators["Road_terrain"] = GenerateRoad_terrain
}
