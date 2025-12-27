package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Grass implements the image.Image interface.
var _ image.Image = (*Grass)(nil)

// Grass is a pattern that generates procedural grass blades.
type Grass struct {
	Null
	BladeWidth  float64
	BladeHeight float64
	FillColor   // Blade color
	Source      image.Image // Background
	Wind        image.Image // Controls bend
	Density     image.Image // Controls density
	Seed        int64
}

// SetBladeWidth sets the average width of a grass blade at the base.
func SetBladeWidth(w float64) func(any) {
	return func(p any) {
		if g, ok := p.(*Grass); ok {
			g.BladeWidth = w
		}
	}
}

// SetBladeHeight sets the average height of a grass blade.
func SetBladeHeight(h float64) func(any) {
	return func(p any) {
		if g, ok := p.(*Grass); ok {
			g.BladeHeight = h
		}
	}
}

// SetWindSource sets the pattern used for wind (bending).
func SetWindSource(src image.Image) func(any) {
	return func(p any) {
		if g, ok := p.(*Grass); ok {
			g.Wind = src
		}
	}
}

// SetDensitySource sets the pattern used for density.
func SetDensitySource(src image.Image) func(any) {
	return func(p any) {
		if g, ok := p.(*Grass); ok {
			g.Density = src
		}
	}
}

// grassHash returns a pseudo-random float64 in [0, 1) based on coordinates and seed.
func grassHash(x, y int, seed int64) float64 {
	var h = uint64(seed) ^ (uint64(x)*73856093) ^ (uint64(y)*19349663)
	h ^= h >> 15
	h *= 0x735a2d97a6428a11
	h ^= h >> 15
	return float64(h&0xffffff) / 16777215.0
}

func (p *Grass) At(x, y int) color.Color {
	// 1. Sample background
	var r, g, b, a uint32
	if p.Source != nil {
		r, g, b, a = p.Source.At(x, y).RGBA()
	}

	// 2. Render Grass
	// We use a grid system.
	// Grid size must be large enough to contain a blade, or we need to search further.
	// A blade of height H and max lean can cover some area.
	// For simplicity, we search 3x3 neighbors and ensure grid size is close to blade height.
	gridSize := int(p.BladeHeight)
	if gridSize < 10 {
		gridSize = 10
	}

	gx := int(math.Floor(float64(x) / float64(gridSize)))
	gy := int(math.Floor(float64(y) / float64(gridSize)))

	// Blade rendering accumulators
	// We need to handle occlusion/blending.
	// Since we are 2D, we can just blend them on top.
	// Ideally we draw back-to-front.
	// Y-sorting: Blades with higher Y (lower on screen) are "closer" and should be drawn last.
	// In our loop, we iterate neighbors. We should process them in Y order?
	// The neighbor loop is small (3x3).
	// We can collect valid blades and sort them?
	// Or just iterate dy from -1 to 1? (Top to bottom neighbors -> Back to front?)
	// Yes, dy=-1 is "above" (smaller y), dy=1 is "below" (larger y).
	// But `At(x, y)` is a single pixel. We are checking which blades cover *this* pixel.
	// The pixel has a fixed Z-order relative to the blades.
	// Actually, for a single pixel, it will only hit one blade usually, or maybe overlap.
	// If multiple blades overlap this pixel, the one "in front" wins.
	// "In front" means the blade's root Y is largest.
	// So we want the blade with the largest `rootY` that covers this pixel.

	var bestBladeDist float64 = -1.0
	var bestBladeColor color.Color
	var hit bool

	// Check neighbors
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			cx := gx + dx
			cy := gy + dy

			// Deterministic RNG for this cell
			h1 := grassHash(cx, cy, p.Seed)
			h2 := grassHash(cx, cy, p.Seed+1)
			h3 := grassHash(cx, cy, p.Seed+2)
			h4 := grassHash(cx, cy, p.Seed+3)

			// Density check
			// We define the "cell center" or random position as the sample point for density.
			rootX := float64(cx*gridSize) + h1*float64(gridSize)
			rootY := float64(cy*gridSize) + h2*float64(gridSize)

			if p.Density != nil {
				_, _, _, da := p.Density.At(int(rootX), int(rootY)).RGBA()
				// Map alpha/luminance to probability
				prob := float64(da) / 65535.0
				if h3 > prob {
					continue
				}
			} else {
				// Default random density (e.g. 1 blade per cell max? or probability?)
				// Let's assume 1 per cell if no density map, or use a threshold.
				if h3 < 0.2 { // Some empty spots
					continue
				}
			}

			// Blade parameters
			height := p.BladeHeight * (0.8 + 0.4*h4) // +/- 20%
			width := p.BladeWidth * (0.8 + 0.4*h1)

			// Initial random lean
			lean := (h2 - 0.5) * 0.5 // -0.25 to 0.25 radians

			// Wind effect
			bend := 0.0
			if p.Wind != nil {
				// Sample wind at blade root
				g16 := color.Gray16Model.Convert(p.Wind.At(int(rootX), int(rootY))).(color.Gray16)
				windVal := float64(g16.Y) / 65535.0
				// Wind bends the grass. Let's say wind moves x positive.
				// Bend factor
				bend = windVal * 20.0 // Amount of x-displacement at tip
				// Also affect lean?
				lean += windVal * 0.5
			}

			// Render check
			// Transform pixel (x, y) to local blade space
			// Translate so root is (0,0)
			px := float64(x) - rootX
			py := float64(y) - rootY

			// Rotate by lean
			// We want to un-rotate the point to align blade with Y axis (up)
			// Blade grows UP (negative Y in image coords usually).
			// Let's define blade going UP (-Y).
			// Local Y should be positive going up.
			// So `localY = -(py)`?
			// Let's standard math:
			// Rotated point:
			cosA := math.Cos(-lean) // Un-rotate
			sinA := math.Sin(-lean)

			lx := px*cosA - py*sinA
			ly := px*sinA + py*cosA

			// Now blade is along -Y axis (if we consider standard image coords where +Y is down)
			// Wait, let's say blade goes from (0,0) to (0, -height).
			// So we want `t` from 0 to 1 as we go from 0 to -height.
			// t = -ly / height.

			t := -ly / height

			if t < 0.0 || t > 1.0 {
				continue
			}

			// Curve (bend)
			// Quadratic bezier style bend.
			// As t goes 0->1, x shifts by `bend * t^2`.
			centerOffset := bend * t * t

			lx -= centerOffset

			// Width profile
			// Taper: w = width * (1 - t)
			halfW := (width * (1.0 - t)) / 2.0

			if math.Abs(lx) < halfW {
				// Hit!
				// We need to resolve Z-order.
				// Larger rootY means closer to camera (lower on screen).
				if rootY > bestBladeDist {
					bestBladeDist = rootY

					// Shading
					// Simple gradient based on t
					// Darker at bottom, lighter at top
					// Or use BladeColor
					br, bg, bb, ba := p.FillColor.FillColor.RGBA()

					// Apply lighting/shading
					// Fake ambient occlusion at bottom
					light := 0.5 + 0.5*t
					// Highlight on one side?
					if lx < 0 {
						light *= 0.8 // Shadow side
					}

					fr := float64(br) * light
					fg := float64(bg) * light
					fb := float64(bb) * light
					fa := float64(ba)

					bestBladeColor = color.RGBA64{
						R: uint16(math.Min(65535, fr)),
						G: uint16(math.Min(65535, fg)),
						B: uint16(math.Min(65535, fb)),
						A: uint16(fa),
					}
					hit = true
				}
			}
		}
	}

	if hit {
		// Blend blade over background
		return blendColors(color.RGBA64{
			R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a),
		}, bestBladeColor)
	}

	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

// blendColors blends foreground over background.
func blendColors(bg, fg color.Color) color.Color {
	br, bg_g, bb, ba := bg.RGBA()
	fr, fg_g, fb, fa := fg.RGBA()

	// Convert to 0-1
	a := float64(fa) / 65535.0
	// invA := 1.0 - a

	// Premultiplied alpha handling?
	// The RGBA() method returns premultiplied alpha values.
	// result = fg + bg * (1 - fg_alpha)

	outR := float64(fr) + float64(br)*(1.0-a)
	outG := float64(fg_g) + float64(bg_g)*(1.0-a)
	outB := float64(fb) + float64(bb)*(1.0-a)
	outA := float64(fa) + float64(ba)*(1.0-a)

	return color.RGBA64{
		R: uint16(math.Min(65535, outR)),
		G: uint16(math.Min(65535, outG)),
		B: uint16(math.Min(65535, outB)),
		A: uint16(math.Min(65535, outA)),
	}
}

func (p *Grass) Bounds() image.Rectangle {
	return p.bounds
}

func (p *Grass) ColorModel() color.Model {
	return color.RGBA64Model
}

// NewGrass creates a new Grass pattern.
func NewGrass(ops ...func(any)) image.Image {
	p := &Grass{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Seed: 1,
		BladeWidth: 6.0,
		BladeHeight: 30.0,
	}
	p.FillColor.FillColor = color.RGBA{10, 150, 10, 255} // Default green

	for _, op := range ops {
		op(p)
	}
	return p
}
