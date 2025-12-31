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

// Circuit: Thin orthogonal traces with small nodes (Redesigned from foundations)
func GenerateCircuitImpl(rect image.Rectangle) image.Image {
	// Base: Dark Green PCB
	bg := NewRect(SetFillColor(color.RGBA{0, 60, 20, 255}), SetBounds(rect))

	// Use a grid-based approach.
	// Divide into cells. Each cell contains a trace segment.

	traces := NewGeneric(func(x, y int) color.Color {
		cellSize := 20
		cx, cy := x/cellSize, y/cellSize

		// Stable hash for cell type
		h := StableHash(cx, cy, 123)

		// Type of trace:
		// 0: Empty
		// 1: Horizontal Line
		// 2: Vertical Line
		// 3: Corner (Top-Left)
		// 4: Corner (Top-Right)
		// 5: Corner (Bottom-Left)
		// 6: Corner (Bottom-Right)
		// 7: Cross
		// 8: Pad

		// Bias towards connections
		t := (h % 100)
		typeCode := 0
		if t < 10 {
			typeCode = 0 // Empty
		} else if t < 40 {
			typeCode = 1 // Horizontal
		} else if t < 70 {
			typeCode = 2 // Vertical
		} else if t < 75 {
			typeCode = 3
		} else if t < 80 {
			typeCode = 4
		} else if t < 85 {
			typeCode = 5
		} else if t < 90 {
			typeCode = 6
		} else if t < 95 {
			typeCode = 7
		} else {
			typeCode = 8 // Pad
		}

		// Local coordinates centered in cell
		// u, v range roughly -10 to 10
		ux := float64(x%cellSize) - float64(cellSize)/2.0
		uy := float64(y%cellSize) - float64(cellSize)/2.0

		thickness := 2.0

		// Helper for drawing lines
		drawLine := func(x1, y1, x2, y2 float64) bool {
			// Distance from point to line segment
			l2 := (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
			if l2 == 0 { return math.Sqrt((ux-x1)*(ux-x1) + (uy-y1)*(uy-y1)) < thickness }
			tVal := ((ux-x1)*(x2-x1) + (uy-y1)*(y2-y1)) / l2
			if tVal < 0 { tVal = 0 }
			if tVal > 1 { tVal = 1 }
			px, py := x1 + tVal*(x2-x1), y1 + tVal*(y2-y1)
			dist := math.Sqrt((ux-px)*(ux-px) + (uy-py)*(uy-py))
			return dist < thickness
		}

		drawPad := func() bool {
			dist := math.Sqrt(ux*ux + uy*uy)
			return dist < 5.0
		}

		isTrace := false
		isPad := false

		half := float64(cellSize)/2.0

		switch typeCode {
		case 1: // Horiz
			isTrace = drawLine(-half, 0, half, 0)
		case 2: // Vert
			isTrace = drawLine(0, -half, 0, half)
		case 3: // Corner TL (Left to Top)
			isTrace = drawLine(-half, 0, 0, 0) || drawLine(0, -half, 0, 0)
		case 4: // Corner TR (Right to Top)
			isTrace = drawLine(half, 0, 0, 0) || drawLine(0, -half, 0, 0)
		case 5: // Corner BL (Left to Bottom)
			isTrace = drawLine(-half, 0, 0, 0) || drawLine(0, half, 0, 0)
		case 6: // Corner BR (Right to Bottom)
			isTrace = drawLine(half, 0, 0, 0) || drawLine(0, half, 0, 0)
		case 7: // Cross
			isTrace = drawLine(-half, 0, half, 0) || drawLine(0, -half, 0, half)
		case 8: // Pad
			isPad = drawPad()
			// And connect to a neighbor? Let's just be a dot for now.
		}

		if isPad {
			return color.RGBA{200, 180, 50, 255} // Gold
		}
		if isTrace {
			return color.RGBA{100, 200, 100, 255} // Light Green
		}

		return color.Transparent
	})

	// Add Chips (Black Rectangles) on top
	chipGen := func(u, v float64, hash uint64) (color.Color, float64) {
		if (hash & 255) < 30 {
			w, h := 0.6, 0.4
			if (hash & 1) == 1 { w, h = 0.4, 0.6 }
			if math.Abs(u) < w/2 && math.Abs(v) < h/2 {
				return color.RGBA{20, 20, 20, 255}, 1.0
			}
		}
		return color.Transparent, 0
	}

	chips := NewScatter(
		SetScatterFrequency(0.04), // Larger grid for chips
		SetScatterDensity(1.0),
		SetScatterGenerator(chipGen),
		SetSpaceColor(color.Transparent), // Important: Transparent background!
		func(i any) { if p, ok := i.(*Scatter); ok { p.Seed = 200 } },
	)

	l1 := NewBlend(bg, traces, BlendNormal)
	l2 := NewBlend(l1, chips, BlendNormal)

	return l2
}

func ExampleNewCircuit() image.Image {
	return GenerateCircuitImpl(image.Rect(0, 0, 150, 150))
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

// Persian Rug: Ornate patterns (Redesigned + Internals Improved)
func GeneratePersianRugImpl(rect image.Rectangle) image.Image {
	// Base: Deep Red
	bg := NewRect(SetFillColor(color.RGBA{60, 10, 10, 255}), SetBounds(rect))

	// Central Medallion
	cx, cy := rect.Dx()/2, rect.Dy()/2
	medBase := NewConcentricRings(
		[]color.Color{color.Transparent},
		SetCenter(cx, cy),
		SetFrequency(0.8),
	)
	medColor := NewColorMap(medBase,
		ColorStop{0.0, color.RGBA{20, 20, 80, 255}}, // Center Blue
		ColorStop{0.3, color.RGBA{150, 120, 50, 255}}, // Gold
		ColorStop{0.35, color.RGBA{20, 20, 80, 255}}, // Blue
		ColorStop{0.6, color.Transparent},
	)

	// Field Pattern: Intricate Internals
	// Rotate 45 Checker + Floral motifs
	lattice := NewRotate(
		NewChecker(
			color.RGBA{100, 30, 30, 255},
			color.Transparent,
			SetSpaceSize(20),
		), 45)

	// Textured Background for Field (Worley)
	fieldTexture := NewWorleyNoise(SetFrequency(0.5))
	fieldOverlay := NewColorMap(fieldTexture,
		ColorStop{0.0, color.RGBA{0, 0, 0, 30}},
		ColorStop{0.5, color.Transparent},
	)

	// Small flowers in grid
	flowers := NewGeneric(func(x, y int) color.Color {
		// Grid 20x20 offset
		tx, ty := (x+10)%20, (y+10)%20
		dx, dy := float64(tx)-10.0, float64(ty)-10.0
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 5 {
			// Flower petals logic?
			// Just a soft dot for now
			return color.RGBA{200, 180, 150, 100}
		}
		return color.Transparent
	})

	// Complex Border System
	borderGen := NewGeneric(func(x, y int) color.Color {
		w, h := rect.Dx(), rect.Dy()
		d := x
		if w-1-x < d { d = w-1-x }
		if y < d { d = y }
		if h-1-y < d { d = h-1-y }

		// Outer Band (Blue)
		if d < 10 { return color.RGBA{20, 20, 60, 255} }
		// Gold Line
		if d < 12 { return color.RGBA{180, 160, 100, 255} }
		// Main Border Band (Red with pattern)
		if d < 25 {
			// Detailed border pattern
			// Alternating geometric shapes
			pat := (x + y) / 5
			if pat%2 == 0 {
				return color.RGBA{100, 30, 30, 255}
			}
			return color.RGBA{120, 40, 40, 255}
		}
		// Gold Line
		if d < 27 { return color.RGBA{180, 160, 100, 255} }

		return color.Transparent
	})

	l1 := NewBlend(bg, fieldOverlay, BlendOverlay) // Add texture
	l2 := NewBlend(l1, lattice, BlendNormal)
	l3 := NewBlend(l2, flowers, BlendNormal) // Internal Details
	l4 := NewBlend(l3, medColor, BlendNormal)
	l5 := NewBlend(l4, borderGen, BlendNormal)

	return l5
}

func ExampleNewPersianRug() image.Image {
	return GeneratePersianRugImpl(image.Rect(0, 0, 150, 150))
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

// Metal Plate: Improved texture (Brushed)
func ExampleNewMetalPlate() image.Image {
	// Base: Brushed Metal
	// Use highly directional noise
	noise := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.1, Octaves: 3}))
	// Scale X to 1.0, Scale Y to 0.05 to stretch vertically? Or horizontally?
	// If we want horizontal brush, we stretch X.
	// Scale(1, 0.05) -> Stretches Y (low freq Y).
	// We want lines. Lines are high freq in one direction, low in other.
	// High freq X (1.0), Low freq Y (0.01).
	brushed := NewScale(noise, ScaleX(1.0), ScaleY(0.02))

	metalBase := NewColorMap(brushed,
		ColorStop{0.0, color.RGBA{120, 120, 125, 255}},
		ColorStop{1.0, color.RGBA{180, 180, 190, 255}},
	)

	// Scratches
	scratchNoise := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.5}))
	// Rotate scratches
	scratches := NewRotate(NewScale(scratchNoise, ScaleX(1.0), ScaleY(0.01)), 15)

	scratchLayer := NewColorMap(scratches,
		ColorStop{0.0, color.RGBA{255, 255, 255, 40}},
		ColorStop{0.3, color.Transparent},
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

	l1 := NewBlend(metalBase, scratchLayer, BlendOverlay)
	l2 := NewBlend(l1, rivets, BlendNormal)

	return l2
}

// Fantasy Frame: Ornate (Fixed Bounds)
func GenerateFantasyFrame(rect image.Rectangle) image.Image {
	// Background: Dark Wood or Stone
	bg := NewRect(SetFillColor(color.RGBA{40, 20, 10, 255}), SetBounds(rect))

	// 1. Ornate Border Pattern (Gold Scrolls)
	scrolls := NewWorleyNoise(SetFrequency(0.15), NoiseSeed(300))
	goldScrolls := NewColorMap(scrolls,
		ColorStop{0.0, color.RGBA{220, 200, 50, 255}}, // Gold
		ColorStop{0.3, color.Transparent},
	)

	// Mask for border area
	// We want the border to be visible within 'rect'.
	// NewRect draws border *inside* bounds.
	// So we must pass 'rect' to NewRect.

	borderMask := NewRect(
		SetBounds(rect), // Explicit bounds
		SetLineSize(25),
		SetLineImageSource(goldScrolls),
		SetFillColor(color.Transparent),
	)

	// Corner Gems?
	gems := NewGeneric(func(x, y int) color.Color {
		// Use rect bounds
		w, h := rect.Dx(), rect.Dy()
		// Relative coords
		rx, ry := x - rect.Min.X, y - rect.Min.Y

		distTL := math.Sqrt(float64(rx*rx + ry*ry))
		distTR := math.Sqrt(float64((rx-w)*(rx-w) + ry*ry))
		distBL := math.Sqrt(float64(rx*rx + (ry-h)*(ry-h)))
		distBR := math.Sqrt(float64((rx-w)*(rx-w) + (ry-h)*(ry-h)))

		if distTL < 20 || distTR < 20 || distBL < 20 || distBR < 20 {
			return color.RGBA{255, 50, 50, 255} // Ruby
		}
		return color.Transparent
	})

	l1 := NewBlend(bg, borderMask, BlendNormal)
	l2 := NewBlend(l1, gems, BlendNormal)

	return l2
}

// We need to update ExampleNewFantasyFrame to use GenerateFantasyFrame or standard bounds
func ExampleNewFantasyFrame() image.Image {
	return GenerateFantasyFrame(image.Rect(0, 0, 150, 150))
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
	return GenerateCircuitImpl(rect)
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
	return GeneratePersianRugImpl(rect)
}
func GenerateLavaFlow(rect image.Rectangle) image.Image {
	return ExampleNewLavaFlow()
}
func GenerateMetalPlate(rect image.Rectangle) image.Image {
	return ExampleNewMetalPlate()
}
// GenerateFantasyFrame is defined above

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
