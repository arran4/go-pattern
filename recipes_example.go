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
	AbstractArtOutputFilename = "abstract_art.png" // Renamed from Crystal
	PixelCamoOutputFilename = "pixel_camo.png"
	CheckerBorderOutputFilename = "checker_border.png"
	WaveBorderOutputFilename = "wave_border.png"
	CarpetOutputFilename = "carpet.png"
	PersianRugOutputFilename = "persian_rug.png"
	LavaFlowOutputFilename = "lava_flow.png"
	MetalPlateOutputFilename = "metal_plate.png"
	FantasyFrameOutputFilename = "fantasy_frame.png"
)

// Dungeon: Stone brick + moss speckles + edge cracks
func ExampleNewDungeon() image.Image {
	stoneTex := NewWorleyNoise(SetFrequency(0.2), NoiseSeed(1))
	stoneCol := NewColorMap(stoneTex,
		ColorStop{0.0, color.RGBA{60, 60, 65, 255}},
		ColorStop{1.0, color.RGBA{40, 40, 45, 255}},
	)
	tiles := NewBrick(
		SetBrickSize(50, 50),
		SetMortarSize(4),
		SetBrickOffset(0),
		SetBrickImages(stoneCol),
		SetMortarImage(NewRect(SetFillColor(color.RGBA{10, 10, 10, 255}))),
	)
	cracks := NewWorleyNoise(
		SetFrequency(0.1),
		SetWorleyOutput(OutputF2MinusF1),
		NoiseSeed(2),
	)
	crackMask := NewColorMap(cracks,
		ColorStop{0.0, color.Black},
		ColorStop{0.05, color.Black},
		ColorStop{0.1, color.Transparent},
	)
	mossNoise := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.05}), NoiseSeed(3))
	mossDetail := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.5}), NoiseSeed(4))
	mossMask := NewBlend(mossNoise, mossDetail, BlendMultiply)
	mossCol := NewColorMap(mossMask,
		ColorStop{0.4, color.Transparent},
		ColorStop{0.6, color.RGBA{50, 100, 50, 150}},
	)
	withCracks := NewBlend(tiles, crackMask, BlendNormal)
	withMoss := NewBlend(withCracks, mossCol, BlendNormal)
	return withMoss
}

// Ice: Pale base + thin cracks + faint gradient
func ExampleNewIce() image.Image {
	base := NewColorMap(
		NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.02}), NoiseSeed(10)),
		ColorStop{0.0, color.RGBA{220, 230, 255, 255}},
		ColorStop{1.0, color.RGBA{240, 250, 255, 255}},
	)
	cracks := NewWorleyNoise(
		SetFrequency(0.08),
		SetWorleyOutput(OutputF2MinusF1),
		NoiseSeed(11),
	)
	crackLines := NewColorMap(cracks,
		ColorStop{0.0, color.RGBA{255, 255, 255, 180}},
		ColorStop{0.05, color.Transparent},
	)
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

// Circuit: Thin orthogonal traces with small nodes (Refined)
func ExampleNewCircuit() image.Image {
	// Base: Dark Green PCB
	bg := NewRect(SetFillColor(color.RGBA{0, 60, 20, 255}))

	// Create orthogonal traces using a specialized approach.
	// Since we don't have a maze generator, we can use Worley F1-F2 with Manhattan distance
	// but thresholded strictly to create "roads".
	// Or use `NewGrid` with random connections.

	// Let's try Worley Manhattan again but with cleaner thresholding.
	// Low frequency for main buses.
	traces := NewWorleyNoise(
		SetWorleyMetric(MetricManhattan),
		SetWorleyOutput(OutputF2MinusF1),
		SetFrequency(0.1),
		NoiseSeed(100),
	)

	// Traces are where F2-F1 is small (centers of edges).
	// Let's make lines.
	traceLayer := NewColorMap(traces,
		ColorStop{0.0, color.RGBA{100, 200, 100, 255}}, // Trace color
		ColorStop{0.05, color.RGBA{100, 200, 100, 255}},
		ColorStop{0.06, color.Transparent},
	)

	// Pads/Vias at cell centers (F1 small)
	pads := NewWorleyNoise(
		SetWorleyMetric(MetricManhattan),
		SetWorleyOutput(OutputF1),
		SetFrequency(0.1),
		NoiseSeed(100), // Same seed to align with traces
	)
	padLayer := NewColorMap(pads,
		ColorStop{0.0, color.RGBA{200, 200, 50, 255}}, // Gold Pad
		ColorStop{0.15, color.RGBA{200, 200, 50, 255}},
		ColorStop{0.16, color.Transparent},
	)

	// Add chips (black rectangles)
	// We use Brick with large spacing? Or Scatter?
	// Let's use Scatter with a rect generator.

	chipGen := func(u, v float64, hash uint64) (color.Color, float64) {
		// Draw Rect
		if math.Abs(u) < 0.6 && math.Abs(v) < 0.6 {
			return color.RGBA{20, 20, 20, 255}, 1.0
		}
		return color.Transparent, 0
	}

	chips := NewScatter(
		SetScatterFrequency(0.04), // Sparse
		SetScatterDensity(0.4),
		SetScatterGenerator(chipGen),
		func(i any) { if p, ok := i.(*Scatter); ok { p.Seed = 101 } },
	)

	l1 := NewBlend(bg, traceLayer, BlendNormal)
	l2 := NewBlend(l1, padLayer, BlendNormal)
	l3 := NewBlend(l2, chips, BlendNormal)

	return l3
}

// Fence: Diagonal diamond grid (Chain link)
func ExampleNewFence() image.Image {
	wires := NewCrossHatch(
		SetLineSize(2),
		SetSpaceSize(18),
		SetAngles(45, 135),
		SetLineColor(color.RGBA{180, 180, 180, 255}),
	)
	grass := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.1}), NoiseSeed(30))
	bg := NewColorMap(grass,
		ColorStop{0.0, color.RGBA{20, 100, 20, 255}},
		ColorStop{1.0, color.RGBA{30, 150, 30, 255}},
	)
	return NewBlend(bg, wires, BlendNormal)
}

// Warning Stripe: Diagonal alternating yellow/black
func ExampleNewStripe() image.Image {
	stripes := NewCrossHatch(
		SetLineSize(20),
		SetSpaceSize(20),
		SetAngle(45),
		SetLineColor(color.RGBA{255, 200, 0, 255}), // Yellow
		SetSpaceColor(color.RGBA{20, 20, 20, 255}), // Black
	)
	return stripes
}

// Abstract Art: Renamed from Crystal (Original implementation)
func ExampleNewAbstractArt() image.Image {
	v := NewVoronoi(
		makePoints(30, 150, 150),
		[]color.Color{
			color.RGBA{200, 220, 255, 255},
			color.RGBA{100, 150, 250, 255},
			color.RGBA{50, 100, 200, 255},
			color.RGBA{150, 200, 240, 255},
			color.RGBA{20, 40, 100, 255},
		},
	)
	shine := NewLinearGradient(
		SetStartColor(color.RGBA{255, 255, 255, 50}),
		SetEndColor(color.Transparent),
		SetAngle(30),
	)
	return NewBlend(v, shine, BlendOverlay)
}

// Pixel Camo: Clustered 2x2 blocks in 3-4 colors
func ExampleNewPixelCamo() image.Image {
	pn := &PerlinNoise{Frequency: 0.04, Seed: 60}
	return NewGeneric(func(x, y int) color.Color {
		s := 10
		qx := (x / s) * s
		qy := (y / s) * s
		c := pn.At(qx, qy)
		g := c.(color.Gray).Y
		v := float64(g) / 255.0
		if v < 0.3 {
			return color.RGBA{30, 25, 20, 255}
		} else if v < 0.5 {
			return color.RGBA{60, 80, 40, 255}
		} else if v < 0.7 {
			return color.RGBA{140, 130, 100, 255}
		}
		return color.RGBA{10, 10, 10, 255}
	})
}

// Checker Border: Classic black/white border strip
func ExampleNewCheckerBorder() image.Image {
	check := NewChecker(
		color.White,
		color.Black,
		SetSpaceSize(20),
	)
	return NewRect(
		SetLineSize(20),
		SetLineImageSource(check),
		SetFillColor(color.Transparent),
	)
}

// Wave Border: Repeating sinusoidal edge
func ExampleNewWaveBorder() image.Image {
	split := NewGeneric(func(x, y int) color.Color {
		if y > 75 {
			return color.RGBA{50, 100, 150, 255}
		}
		return color.Transparent
	})
	sine := NewGeneric(func(x, y int) color.Color {
		v := math.Sin(float64(x) * 0.1) * 10.0
		val := 128.0 + v
		return color.Gray{Y: uint8(val)}
	})
	warp := NewWarp(split,
		WarpDistortion(sine),
		WarpYScale(1.0),
		WarpXScale(0.0),
	)
	return warp
}

// Carpet: Visual interest increased
func ExampleNewCarpet() image.Image {
	bg := NewRect(SetFillColor(color.RGBA{80, 0, 0, 255})) // Dark Red Base

	// Main Pattern: Large Diamonds
	d1 := NewRotate(
		NewChecker(
			color.RGBA{160, 120, 40, 255}, // Gold
			color.Transparent,
			SetSpaceSize(30),
		),
		45,
	)

	// Secondary Pattern: Smaller overlay diamonds
	d2 := NewRotate(
		NewChecker(
			color.Transparent,
			color.RGBA{0, 0, 0, 60}, // Shadow
			SetSpaceSize(10),
		),
		45,
	)

	// Border Elements?
	// Striped background for texture
	stripes := NewCrossHatch(SetLineSize(1), SetSpaceSize(3), SetAngle(0), SetLineColor(color.RGBA{0,0,0,30}))

	l1 := NewBlend(bg, stripes, BlendNormal)
	l2 := NewBlend(l1, d1, BlendNormal)
	l3 := NewBlend(l2, d2, BlendNormal)

	return l3
}

// Persian Rug: Ornate patterns
func ExampleNewPersianRug() image.Image {
	// Base: Deep Blue
	bg := NewRect(SetFillColor(color.RGBA{10, 10, 60, 255}))

	// Medallion (Center)
	// Use Concentric Rings with color map
	// NewConcentricRings takes []color.Color as first arg
	center := NewConcentricRings(
		[]color.Color{color.Black, color.White}, // Placeholder, using ColorMap over it
		SetCenter(75, 75),
		SetFrequency(0.5),
	)
	medallion := NewColorMap(center,
		ColorStop{0.0, color.RGBA{200, 50, 50, 255}}, // Red Center
		ColorStop{0.2, color.RGBA{200, 180, 100, 255}}, // Gold Ring
		ColorStop{0.4, color.RGBA{10, 10, 60, 255}}, // Blue Gap
		ColorStop{0.5, color.RGBA{200, 50, 50, 255}}, // Red Ring
		ColorStop{1.0, color.Transparent},
	)

	// Field Pattern: Small intricate details (Worley)
	fieldNoise := NewWorleyNoise(SetFrequency(0.3), NoiseSeed(200))
	field := NewColorMap(fieldNoise,
		ColorStop{0.0, color.RGBA{100, 100, 200, 100}},
		ColorStop{0.5, color.Transparent},
	)

	// Borders
	border := NewRect(
		SetLineSize(15),
		SetLineColor(color.RGBA{150, 30, 30, 255}), // Red Border
		SetFillColor(color.Transparent),
	)

	// Combine
	l1 := NewBlend(bg, field, BlendNormal)
	l2 := NewBlend(l1, medallion, BlendNormal)
	l3 := NewBlend(l2, border, BlendNormal)

	return l3
}

// Lava Flow: Dark base + bright streaks + subtle noise
func ExampleNewLavaFlow() image.Image {
	base := NewColorMap(
		NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.05}), NoiseSeed(70)),
		ColorStop{0.0, color.RGBA{20, 0, 0, 255}},
		ColorStop{0.6, color.RGBA{60, 10, 0, 255}},
		ColorStop{1.0, color.RGBA{100, 20, 0, 255}},
	)
	rivers := NewColorMap(
		NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.03}), NoiseSeed(71)),
		ColorStop{0.0, color.RGBA{255, 200, 0, 255}}, // Bright Yellow
		ColorStop{0.2, color.RGBA{255, 50, 0, 255}},  // Red
		ColorStop{0.4, color.Transparent},            // Cooled rock
	)
	flow := NewWarp(rivers,
		WarpDistortion(NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.02}), NoiseSeed(72))),
		WarpScale(20.0),
	)
	return NewBlend(base, flow, BlendNormal)
}

// Metal Plate: Improved texture
func ExampleNewMetalPlate() image.Image {
	// Base: Brushed Metal
	// Use directional noise
	noise := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.2, Octaves: 5}))
	// Scale X/Y differently to make it "brushed"
	// We don't have scale options on PerlinNoise directly.
	// We can use NewScale pattern.
	brushed := NewScale(noise, ScaleX(1.0), ScaleY(0.1))

	metalBase := NewColorMap(brushed,
		ColorStop{0.0, color.RGBA{100, 100, 105, 255}},
		ColorStop{1.0, color.RGBA{180, 180, 190, 255}},
	)

	// Rivets
	rivets := NewGeneric(func(x, y int) color.Color {
		s := 50
		nx := (x + s/2) / s * s
		ny := (y + s/2) / s * s
		dx := x - nx
		dy := y - ny
		dist := math.Sqrt(float64(dx*dx + dy*dy))

		if dist < 6 {
			// Rivet shading (simple gradient)
			v := uint8(255 - dist*20)
			return color.RGBA{v, v, v, 255}
		}
		return color.Transparent
	})

	return NewBlend(metalBase, rivets, BlendNormal)
}

// Fantasy Frame: Ornate
func ExampleNewFantasyFrame() image.Image {
	// Background: Dark Wood or Stone
	bg := NewRect(SetFillColor(color.RGBA{40, 20, 10, 255}))

	// 1. Ornate Border Pattern (Gold Scrolls)
	// Use a pattern that looks scroll-like? Maybe Worley?
	scrolls := NewWorleyNoise(SetFrequency(0.15), NoiseSeed(300))
	goldScrolls := NewColorMap(scrolls,
		ColorStop{0.0, color.RGBA{220, 200, 50, 255}}, // Gold
		ColorStop{0.3, color.Transparent},
	)

	// Mask for border area (e.g. outer 20px)
	borderMask := NewRect(
		SetLineSize(25),
		SetLineImageSource(goldScrolls),
		SetFillColor(color.Transparent), // Transparent center (to show image... or just black here)
	)

	// Corner Gems?
	gems := NewGeneric(func(x, y int) color.Color {
		// Corners of 150x150: (0,0), (150,0), ...
		// But generic doesn't know bounds.
		// Assuming 150x150 for demo.
		w, h := 150, 150
		distTL := math.Sqrt(float64(x*x + y*y))
		distTR := math.Sqrt(float64((x-w)*(x-w) + y*y))
		distBL := math.Sqrt(float64(x*x + (y-h)*(y-h)))
		distBR := math.Sqrt(float64((x-w)*(x-w) + (y-h)*(y-h)))

		if distTL < 20 || distTR < 20 || distBL < 20 || distBR < 20 {
			return color.RGBA{255, 50, 50, 255} // Ruby
		}
		return color.Transparent
	})

	l1 := NewBlend(bg, borderMask, BlendNormal)
	l2 := NewBlend(l1, gems, BlendNormal)

	return l2
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
func GenerateAbstractArt(rect image.Rectangle) image.Image {
	return ExampleNewAbstractArt()
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
func GeneratePersianRug(rect image.Rectangle) image.Image {
	return ExampleNewPersianRug()
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
	RegisterGenerator("AbstractArt", GenerateAbstractArt)
	RegisterGenerator("PixelCamo", GeneratePixelCamo)
	RegisterGenerator("CheckerBorder", GenerateCheckerBorder)
	RegisterGenerator("WaveBorder", GenerateWaveBorder)
	RegisterGenerator("Carpet", GenerateCarpet)
	RegisterGenerator("PersianRug", GeneratePersianRug)
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
	RegisterReferences("AbstractArt", ref)
	RegisterReferences("PixelCamo", ref)
	RegisterReferences("CheckerBorder", ref)
	RegisterReferences("WaveBorder", ref)
	RegisterReferences("Carpet", ref)
	RegisterReferences("PersianRug", ref)
	RegisterReferences("LavaFlow", ref)
	RegisterReferences("MetalPlate", ref)
	RegisterReferences("FantasyFrame", ref)
}
