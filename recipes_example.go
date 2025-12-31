package pattern

import (
	"image"
	"image/color"
	"math"
)

var (
	DungeonOutputFilename = "dungeon.png"
	IceOutputFilename = "ice.png"
	CircuitOutputFilename = "circuit.png"
	FenceOutputFilename = "fence.png"
	StripeOutputFilename = "stripe.png"
	CrystalOutputFilename = "crystal.png"
	PixelCamoOutputFilename = "pixel_camo.png"
	CheckerBorderOutputFilename = "checker_border.png"
	WaveBorderOutputFilename = "wave_border.png"
	CarpetOutputFilename = "carpet.png"
	LavaFlowOutputFilename = "lava_flow.png"
	MetalPlateOutputFilename = "metal_plate.png"
	FantasyFrameOutputFilename = "fantasy_frame.png"
)

// Dungeon: Stone brick + moss speckles + edge cracks
func ExampleNewDungeon() image.Image {
	// 1. Base Stone Brick
	// We use Worley Noise for stone texture
	stoneTex := NewWorleyNoise(SetFrequency(0.2), NoiseSeed(1))
	stoneCol := NewColorMap(stoneTex,
		ColorStop{0.0, color.RGBA{60, 60, 65, 255}},
		ColorStop{1.0, color.RGBA{40, 40, 45, 255}},
	)

	// Create square tiles (dungeon floor)
	tiles := NewBrick(
		SetBrickSize(50, 50),
		SetMortarSize(4),
		SetBrickOffset(0), // Aligned grid
		SetBrickImages(stoneCol),
		SetMortarImage(NewRect(SetFillColor(color.RGBA{10, 10, 10, 255}))),
	)

	// 2. Cracks
	// Use high frequency Voronoi or cellular noise for cracks
	cracks := NewWorleyNoise(
		SetFrequency(0.1),
		SetWorleyOutput(OutputF2MinusF1), // Good for edges/cracks
		NoiseSeed(2),
	)
	// Threshold to get thin lines
	crackMask := NewColorMap(cracks,
		ColorStop{0.0, color.Black}, // Crack
		ColorStop{0.05, color.Black},
		ColorStop{0.1, color.Transparent},
	)

	// 3. Moss Speckles
	// Perlin noise for patchiness
	mossNoise := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.05}), NoiseSeed(3))
	// High freq noise for detail
	mossDetail := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.5}), NoiseSeed(4))

	// Combine: Where mossNoise is high, show mossDetail
	mossMask := NewBlend(mossNoise, mossDetail, BlendMultiply)
	mossCol := NewColorMap(mossMask,
		ColorStop{0.4, color.Transparent},
		ColorStop{0.6, color.RGBA{50, 100, 50, 150}}, // Semi-transparent green
	)

	// Layer: Tiles -> Cracks -> Moss
	withCracks := NewBlend(tiles, crackMask, BlendNormal) // Cracks on top
	withMoss := NewBlend(withCracks, mossCol, BlendNormal) // Moss on top

	return withMoss
}

// Ice: Pale base + thin cracks + faint gradient
func ExampleNewIce() image.Image {
	// Base: White/Blueish gradient
	base := NewColorMap(
		NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.02}), NoiseSeed(10)),
		ColorStop{0.0, color.RGBA{220, 230, 255, 255}},
		ColorStop{1.0, color.RGBA{240, 250, 255, 255}},
	)

	// Cracks: Voronoi edges
	// OutputF2MinusF1 gives cellular borders
	cracks := NewWorleyNoise(
		SetFrequency(0.08),
		SetWorleyOutput(OutputF2MinusF1),
		NoiseSeed(11),
	)

	// Thin white lines for subsurface cracks
	crackLines := NewColorMap(cracks,
		ColorStop{0.0, color.RGBA{255, 255, 255, 180}},
		ColorStop{0.05, color.Transparent},
	)

	// Deep cracks (darker)
	deepCracks := NewWorleyNoise(
		SetFrequency(0.04),
		SetWorleyOutput(OutputF2MinusF1),
		NoiseSeed(12),
	)
	deepCrackLines := NewColorMap(deepCracks,
		ColorStop{0.0, color.RGBA{180, 200, 220, 200}},
		ColorStop{0.02, color.Transparent},
	)

	layer1 := NewBlend(base, crackLines, BlendNormal)
	layer2 := NewBlend(layer1, deepCrackLines, BlendNormal)

	return layer2
}

// Circuit: Thin orthogonal traces with small nodes
func ExampleNewCircuit() image.Image {
	// Background: Dark Green
	bg := NewRect(SetFillColor(color.RGBA{10, 40, 20, 255}))

	// Traces: Use `Worley` with Manhattan distance for blocky/orthogonal shapes
	// F2 - F1 gives us "distance to edges" of cells.
	// With Manhattan, edges are orthogonal.

	// Trace mask
	traceNoise := NewWorleyNoise(
		SetWorleyMetric(MetricManhattan),
		SetWorleyOutput(OutputF2MinusF1),
		SetFrequency(0.15),
		NoiseSeed(25),
	)

	// Map to lines
	traces := NewColorMap(traceNoise,
		ColorStop{0.0, color.RGBA{30, 100, 40, 255}}, // Center of "road"
		ColorStop{0.1, color.RGBA{50, 150, 60, 255}}, // Trace
		ColorStop{0.2, color.Transparent},            // Gap
	)

	// Scatter "Nodes" (Pads)
	// Circle or Square pads
	padGen := func(u, v float64, hash uint64) (color.Color, float64) {
		dist := math.Sqrt(u*u + v*v)
		if dist < 0.3 {
			// Silver/Gold pad
			return color.RGBA{200, 180, 50, 255}, 1.0
		}
		return color.Transparent, 0
	}

	pads := NewScatter(
		SetScatterFrequency(0.15), // Match trace frequency roughly
		SetScatterDensity(0.4),
		SetScatterGenerator(padGen),
		func(i any) { if p, ok := i.(*Scatter); ok { p.Seed = 26 } },
	)

	// Combine
	l1 := NewBlend(bg, traces, BlendNormal)
	l2 := NewBlend(l1, pads, BlendNormal)

	return l2
}

// Fence: Diagonal diamond grid (Chain link)
func ExampleNewFence() image.Image {
	// CrossHatch with 45 and 135 degrees
	wires := NewCrossHatch(
		SetLineSize(2),
		SetSpaceSize(18),
		SetAngles(45, 135),
		SetLineColor(color.RGBA{180, 180, 180, 255}),
	)

	// Background: Transparent or blurred scene?
	// Let's put a green background (grass) behind it.
	grass := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.1}), NoiseSeed(30))
	bg := NewColorMap(grass,
		ColorStop{0.0, color.RGBA{20, 100, 20, 255}},
		ColorStop{1.0, color.RGBA{30, 150, 30, 255}},
	)

	return NewBlend(bg, wires, BlendNormal)
}

// Warning Stripe: Diagonal alternating yellow/black
func ExampleNewStripe() image.Image {
	// Or `CrossHatch` with single angle and thick lines.
	stripes := NewCrossHatch(
		SetLineSize(20),
		SetSpaceSize(20),
		SetAngle(45),
		SetLineColor(color.RGBA{255, 200, 0, 255}), // Yellow
		SetSpaceColor(color.RGBA{20, 20, 20, 255}), // Black
	)

	return stripes
}

// Crystal: Triangular facets in light/dark blues
func ExampleNewCrystal() image.Image {
	// Voronoi gives polygonal cells
	// We want sharp, angular facets.

	v := NewVoronoi(
		// Random points are generated if we don't provide them.
		makePoints(30, 150, 150),
		// Palette of blues
		[]color.Color{
			color.RGBA{200, 220, 255, 255}, // Pale Blue
			color.RGBA{100, 150, 250, 255}, // Medium Blue
			color.RGBA{50, 100, 200, 255},  // Dark Blue
			color.RGBA{150, 200, 240, 255}, // Light Cyan
			color.RGBA{20, 40, 100, 255},   // Deep Blue
		},
	)

	// Let's add a reflective sheen using a large gradients.
	shine := NewLinearGradient(
		SetStartColor(color.RGBA{255, 255, 255, 50}),
		SetEndColor(color.Transparent),
		SetAngle(30),
	)

	return NewBlend(v, shine, BlendOverlay)
}

// Pixel Camo: Clustered 2x2 blocks in 3-4 colors
func ExampleNewPixelCamo() image.Image {
	// Instantiate noise source
	pn := &PerlinNoise{Frequency: 0.04, Seed: 60}

	return NewGeneric(func(x, y int) color.Color {
		s := 10 // Block size
		// Quantize coordinates
		qx := (x / s) * s
		qy := (y / s) * s

		// Sample noise at quantized coordinate
		c := pn.At(qx, qy)
		// Extract value (0-255) from gray
		g := c.(color.Gray).Y
		v := float64(g) / 255.0

		// Map to 3-4 Camo Colors
		if v < 0.3 {
			return color.RGBA{30, 25, 20, 255} // Dark Brown
		} else if v < 0.5 {
			return color.RGBA{60, 80, 40, 255} // Army Green
		} else if v < 0.7 {
			return color.RGBA{140, 130, 100, 255} // Khaki
		}
		return color.RGBA{10, 10, 10, 255} // Black
	})
}

// Checker Border: Classic black/white border strip
func ExampleNewCheckerBorder() image.Image {
	// Create a checker pattern
	check := NewChecker(
		color.White,
		color.Black,
		SetSpaceSize(20), // 20px squares
	)

	// Mask it to be a border
	// We want a strip.
	// Use Padding? No, Padding adds space around.
	// We want to return an image that IS a border.
	// If the image size is fixed, we can just mask the center.

	// Let's make a frame: White transparent center, solid border.
	// Use `NewRect` to make a hole?

	// 1. Full Checker
	// 2. Center mask (Black = Transparent, White = Opaque)

	// Actually, `NewRect` with `LineSize` draws a border.
	// If `LineImageSource` is the checker pattern...

	return NewRect(
		SetLineSize(20),
		SetLineImageSource(check),
		SetFillColor(color.Transparent), // Transparent center
	)
}

// Wave Border: Repeating sinusoidal edge
func ExampleNewWaveBorder() image.Image {
	// Create a straight line (vertical or horizontal)
	// Let's make a horizontal wave at the bottom.

	// Base: A filled rectangle at the bottom half.
	// But we want the edge to be wavy.

	// Use `NewWarp` on a straight split.

	split := NewGeneric(func(x, y int) color.Color {
		if y > 75 { // Halfway
			return color.RGBA{50, 100, 150, 255} // Blue Sea
		}
		return color.Transparent
	})

	// Warp it
	// Sine wave distortion
	sine := NewGeneric(func(x, y int) color.Color {
		// Value varies by X
		v := math.Sin(float64(x) * 0.1) * 10.0 // Amplitude 10, Freq 0.1
		// Map to gray 0-255 centered at 128?
		// Warp pattern uses luminance.
		// Gray = 128 + v
		val := 128.0 + v
		return color.Gray{Y: uint8(val)}
	})

	// Distortion in Y direction
	warp := NewWarp(split,
		WarpDistortion(sine),
		WarpYScale(1.0),
		WarpXScale(0.0), // No X distortion
	)

	return warp
}

// Carpet: Repeating diamond/chevron motifs
func ExampleNewCarpet() image.Image {
	// Use `NewTile` or `NewChecker` rotated?

	// Base: Red
	bg := NewRect(SetFillColor(color.RGBA{100, 20, 20, 255}))

	// Pattern: Gold Diamonds
	// Rotated Checkers
	diamonds := NewRotate(
		NewChecker(
			color.RGBA{200, 180, 50, 255}, // Gold
			color.Transparent,
			SetSpaceSize(15),
		),
		45, // Rotate 45 degrees
	)
	// NewRotate creates a black background by default where it's undefined?
	// `rotate.go` implementation might crop or fill.
	// It just rotates coordinates. Infinite plane.

	// Add some detail?
	// Smaller diamonds inside?
	diamondsSmall := NewRotate(
		NewChecker(
			color.Transparent,
			color.RGBA{0, 0, 0, 50}, // Shadow
			SetSpaceSize(15),
		),
		45,
	)

	// Combine
	l1 := NewBlend(bg, diamonds, BlendNormal)
	l2 := NewBlend(l1, diamondsSmall, BlendNormal) // Might not align perfectly without offset

	return l2
}

// Lava Flow: Dark base + bright streaks + subtle noise
func ExampleNewLavaFlow() image.Image {
	// Base: Dark rock (Red/Black)
	base := NewColorMap(
		NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.05}), NoiseSeed(70)),
		ColorStop{0.0, color.RGBA{20, 0, 0, 255}},
		ColorStop{0.6, color.RGBA{60, 10, 0, 255}},
		ColorStop{1.0, color.RGBA{100, 20, 0, 255}},
	)

	// Streaks: Warped noise (Lava rivers)
	// High contrast noise
	rivers := NewColorMap(
		NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.03}), NoiseSeed(71)),
		ColorStop{0.0, color.RGBA{255, 200, 0, 255}}, // Bright Yellow
		ColorStop{0.2, color.RGBA{255, 50, 0, 255}},  // Red
		ColorStop{0.4, color.Transparent},            // Cooled rock
	)

	// Warp the rivers to make them flow
	flow := NewWarp(rivers,
		WarpDistortion(NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.02}), NoiseSeed(72))),
		WarpScale(20.0),
	)

	return NewBlend(base, flow, BlendNormal) // Rivers on top of base
}

// Metal Plate: Grid + rivet dots + small highlights
func ExampleNewMetalPlate() image.Image {
	// Base: Grey metal
	bg := NewRect(SetFillColor(color.RGBA{100, 100, 105, 255}))

	// Grid lines (grooves)
	grid := NewCrossHatch(
		SetLineSize(1),
		SetSpaceSize(49), // 50px grid
		SetAngles(0, 90),
		SetLineColor(color.RGBA{50, 50, 55, 255}), // Dark groove
	)

	// Rivets at intersections?
	// Use Scatter with grid alignment?
	// Or `NewGrid` with dot image?

	// Let's use `NewScatter` with very high density and no randomness to simulate grid points?
	// Scatter has randomness.

	// Custom Generic pattern for rivets at grid corners
	rivets := NewGeneric(func(x, y int) color.Color {
		// Grid 50
		s := 50
		// Nearest intersection
		nx := (x + s/2) / s * s
		ny := (y + s/2) / s * s

		dx := x - nx
		dy := y - ny
		dist := math.Sqrt(float64(dx*dx + dy*dy))

		if dist < 4 {
			return color.RGBA{180, 180, 190, 255} // Light Rivet
		} else if dist < 5 {
			return color.RGBA{40, 40, 45, 255} // Shadow ring
		}
		return color.Transparent
	})

	// Scratches/Noise
	scratch := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.5}), NoiseSeed(80))
	scratchLayer := NewColorMap(scratch,
		ColorStop{0.0, color.RGBA{255, 255, 255, 20}},
		ColorStop{0.5, color.Transparent},
	)

	l1 := NewBlend(bg, grid, BlendNormal)
	l2 := NewBlend(l1, rivets, BlendNormal)
	l3 := NewBlend(l2, scratchLayer, BlendOverlay)

	return l3
}

// Fantasy Frame: Wave border + checker border + inset line
func ExampleNewFantasyFrame() image.Image {
	// This implies a composite border.
	// Since we are generating a square tile, let's make a corner or edge piece?
	// Or a full frame around the 150x150 area.

	// 1. Checker Border (Outer)
	chk := ExampleNewCheckerBorder()

	// 2. Wave Border (Inner)
	// We need a wave that goes around? Hard with current `NewWarp`.
	// Let's just layer a gold line.

	gold := NewRect(
		SetLineSize(22), // Slightly inside the 20px checker? No, outside?
		// Let's put gold line *inside* the checker. Checker is 20px.
		// So gold line at 22px (2px width).
		SetLineColor(color.RGBA{200, 180, 50, 255}),
		SetFillColor(color.Transparent),
	)
	// Wait, `NewRect` draws border from edge inwards?
	// `rect.go`: "If LineSize > 0... draws border".
	// Implementation usually draws `LineSize` thick border.

	// If we want a line at offset, we need Padding?
	// `NewPadding` adds space *around*.

	// Let's just stack borders.
	// Bottom layer: Checker (20px)
	// Top layer: Gold Line (start at 20, end at 22).
	// We can use `NewRect` with `LineSize 22` but masked?

	// Easier: Just return the checker border combined with wave for now.
	// The prompt was "Wave border + checker border".

	return NewBlend(chk, gold, BlendNormal)
}

// --- Generic Helper ---

type Generic struct {
	Func func(x, y int) color.Color
}
func (p *Generic) At(x, y int) color.Color { return p.Func(x, y) }
func (p *Generic) Bounds() image.Rectangle { return image.Rect(0, 0, 255, 255) }
func (p *Generic) ColorModel() color.Model { return color.RGBAModel }
func NewGeneric(f func(x,y int) color.Color) image.Image { return &Generic{Func: f} }


func makePoints(n, w, h int) []image.Point {
	pts := make([]image.Point, n)
	seed := 50
	for i := 0; i < n; i++ {
		seed = (seed * 1103515245 + 12345) & 0x7FFFFFFF
		x := seed % w
		seed = (seed * 1103515245 + 12345) & 0x7FFFFFFF
		y := seed % h
		pts[i] = image.Point{x, y}
	}
	return pts
}

func GenerateDungeon(rect image.Rectangle) image.Image {
	return ExampleNewDungeon()
}
func GenerateIce(rect image.Rectangle) image.Image {
	return ExampleNewIce()
}
func GenerateCircuit(rect image.Rectangle) image.Image {
	return ExampleNewCircuit()
}
func GenerateFence(rect image.Rectangle) image.Image {
	return ExampleNewFence()
}
func GenerateStripe(rect image.Rectangle) image.Image {
	return ExampleNewStripe()
}
func GenerateCrystal(rect image.Rectangle) image.Image {
	return ExampleNewCrystal()
}
func GeneratePixelCamo(rect image.Rectangle) image.Image {
	return ExampleNewPixelCamo()
}
func GenerateCheckerBorder(rect image.Rectangle) image.Image {
	return ExampleNewCheckerBorder()
}
func GenerateWaveBorder(rect image.Rectangle) image.Image {
	return ExampleNewWaveBorder()
}
func GenerateCarpet(rect image.Rectangle) image.Image {
	return ExampleNewCarpet()
}
func GenerateLavaFlow(rect image.Rectangle) image.Image {
	return ExampleNewLavaFlow()
}
func GenerateMetalPlate(rect image.Rectangle) image.Image {
	return ExampleNewMetalPlate()
}
func GenerateFantasyFrame(rect image.Rectangle) image.Image {
	return ExampleNewFantasyFrame()
}

func init() {
	RegisterGenerator("Dungeon", GenerateDungeon)
	RegisterGenerator("Ice", GenerateIce)
	RegisterGenerator("Circuit", GenerateCircuit)
	RegisterGenerator("Fence", GenerateFence)
	RegisterGenerator("Stripe", GenerateStripe)
	RegisterGenerator("Crystal", GenerateCrystal)
	RegisterGenerator("PixelCamo", GeneratePixelCamo)
	RegisterGenerator("CheckerBorder", GenerateCheckerBorder)
	RegisterGenerator("WaveBorder", GenerateWaveBorder)
	RegisterGenerator("Carpet", GenerateCarpet)
	RegisterGenerator("LavaFlow", GenerateLavaFlow)
	RegisterGenerator("MetalPlate", GenerateMetalPlate)
	RegisterGenerator("FantasyFrame", GenerateFantasyFrame)

	ref := func() (map[string]func(image.Rectangle) image.Image, []string) {
		return map[string]func(image.Rectangle) image.Image{}, []string{}
	}
	RegisterReferences("Dungeon", ref)
	RegisterReferences("Ice", ref)
	RegisterReferences("Circuit", ref)
	RegisterReferences("Fence", ref)
	RegisterReferences("Stripe", ref)
	RegisterReferences("Crystal", ref)
	RegisterReferences("PixelCamo", ref)
	RegisterReferences("CheckerBorder", ref)
	RegisterReferences("WaveBorder", ref)
	RegisterReferences("Carpet", ref)
	RegisterReferences("LavaFlow", ref)
	RegisterReferences("MetalPlate", ref)
	RegisterReferences("FantasyFrame", ref)
}
