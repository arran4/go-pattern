package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure FineGrid implements the image.Image interface.
var _ image.Image = (*FineGrid)(nil)

// FineGrid renders a thin grid on black with a soft outer glow and mild chromatic aberration.
type FineGrid struct {
	Null
	CellSize          int
	GlowRadius        float64
	Hue               float64
	AberrationOffset  int
	GlowStrength      float64
	LineCoreStrength  float64
	BackgroundOpacity float64
}

func (g *FineGrid) ColorModel() color.Model {
	return color.RGBAModel
}

func (g *FineGrid) Bounds() image.Rectangle {
	return g.bounds
}

func (g *FineGrid) At(x, y int) color.Color {
	cell := g.CellSize
	if cell <= 0 {
		cell = 12
	}

	glowRadius := g.GlowRadius
	if glowRadius <= 0 {
		glowRadius = 3.0
	}

	aberration := g.AberrationOffset
	if aberration == 0 {
		aberration = 1
	}

	baseCol := hueToRGB(g.Hue, 1.0, 1.0)

	rStrength := g.channelStrength(x+aberration, y, cell, glowRadius)
	gStrength := g.channelStrength(x, y, cell, glowRadius)
	bStrength := g.channelStrength(x-aberration, y, cell, glowRadius)

	r := clamp255(float64(baseCol.R) * rStrength)
	gr := clamp255(float64(baseCol.G) * gStrength)
	b := clamp255(float64(baseCol.B) * bStrength)

	alpha := clamp255(math.Max(rStrength, math.Max(gStrength, bStrength)) * 255.0 * math.Max(0, 1.0-g.BackgroundOpacity))

	return color.NRGBA{
		R: uint8(r),
		G: uint8(gr),
		B: uint8(b),
		A: uint8(alpha),
	}
}

func (g *FineGrid) channelStrength(x, y, cell int, glowRadius float64) float64 {
	localX := x - g.bounds.Min.X
	localY := y - g.bounds.Min.Y

	modX := positiveMod(localX, cell)
	modY := positiveMod(localY, cell)

	distX := math.Min(float64(modX), float64(cell-modX))
	distY := math.Min(float64(modY), float64(cell-modY))
	dist := math.Min(distX, distY)

	core := g.LineCoreStrength
	if core <= 0 {
		core = 1.0
	}
	coreStrength := math.Max(0.0, core-2.0*dist)

	glow := g.GlowStrength
	if glow <= 0 {
		glow = 0.85
	}
	glowStrength := glow * math.Exp(-(dist*dist)/(2.0*glowRadius*glowRadius))

	return math.Min(1.0, coreStrength+glowStrength)
}

func positiveMod(v, m int) int {
	if m == 0 {
		return 0
	}
	r := v % m
	if r < 0 {
		r += m
	}
	return r
}

// hueToRGB converts HSV (in degrees, 0..360) to an sRGB color with the given saturation and value.
func hueToRGB(h, s, v float64) color.RGBA {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return color.RGBA{
		R: uint8(clamp255((r + m) * 255)),
		G: uint8(clamp255((g + m) * 255)),
		B: uint8(clamp255((b + m) * 255)),
		A: 255,
	}
}


// Option helpers

type hasFineGridCellSize interface {
	SetFineGridCellSize(int)
}

// SetFineGridCellSize configures the cell spacing of the grid.
func SetFineGridCellSize(size int) func(any) {
	return func(i any) {
		if h, ok := i.(hasFineGridCellSize); ok {
			h.SetFineGridCellSize(size)
		}
	}
}

func (g *FineGrid) SetFineGridCellSize(size int) {
	g.CellSize = size
}

type hasFineGridGlowRadius interface {
	SetFineGridGlowRadius(float64)
}

// SetFineGridGlowRadius configures the softness of the glow.
func SetFineGridGlowRadius(radius float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasFineGridGlowRadius); ok {
			h.SetFineGridGlowRadius(radius)
		}
	}
}

func (g *FineGrid) SetFineGridGlowRadius(radius float64) {
	g.GlowRadius = radius
}

type hasFineGridHue interface {
	SetFineGridHue(float64)
}

// SetFineGridHue configures the hue (in degrees 0..360) of the glow.
func SetFineGridHue(h float64) func(any) {
	return func(i any) {
		if v, ok := i.(hasFineGridHue); ok {
			v.SetFineGridHue(h)
		}
	}
}

func (g *FineGrid) SetFineGridHue(h float64) {
	g.Hue = h
}

// SetFineGridAberration configures how many pixels each color channel is offset for chromatic aberration.
func SetFineGridAberration(offset int) func(any) {
	return func(i any) {
		if v, ok := i.(*FineGrid); ok {
			v.AberrationOffset = offset
		}
	}
}

// SetFineGridGlowStrength configures how strong the falloff glow is.
func SetFineGridGlowStrength(strength float64) func(any) {
	return func(i any) {
		if v, ok := i.(*FineGrid); ok {
			v.GlowStrength = strength
		}
	}
}

// SetFineGridLineStrength configures the thickness of the core line profile.
func SetFineGridLineStrength(strength float64) func(any) {
	return func(i any) {
		if v, ok := i.(*FineGrid); ok {
			v.LineCoreStrength = strength
		}
	}
}

// SetFineGridBackgroundFade lets the glow darken into the background (0=no fade, 1=full transparent).
func SetFineGridBackgroundFade(fade float64) func(any) {
	return func(i any) {
		if v, ok := i.(*FineGrid); ok {
			v.BackgroundOpacity = fade
		}
	}
}

// NewFineGrid builds a grid with glow and chromatic aberration.
func NewFineGrid(ops ...func(any)) image.Image {
	p := &FineGrid{
		Null: Null{
			bounds: image.Rect(0, 0, 512, 512),
		},
		CellSize:         12,
		GlowRadius:       3.0,
		Hue:              200,
		AberrationOffset: 1,
		GlowStrength:     0.85,
		LineCoreStrength: 1.2,
	}

	for _, op := range ops {
		op(p)
	}

	return p
}

// NewDemoFineGrid produces a demo variant for readme.md pre-populated values.
func NewDemoFineGrid(ops ...func(any)) image.Image {
	args := []func(any){
		SetFineGridCellSize(10),
		SetFineGridGlowRadius(3.5),
		SetFineGridHue(200),
		SetFineGridAberration(1),
		SetFineGridGlowStrength(0.9),
		SetFineGridLineStrength(1.3),
	}
	args = append(args, ops...)
	return NewFineGrid(args...)
}
