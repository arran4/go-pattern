package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure ThreadBands implements image.Image.
var _ image.Image = (*ThreadBands)(nil)

// ThreadBands renders orthogonal thread bands with alternating tones and subtle fiber noise.
// The pattern simulates a simple weave with configurable thread width, shadow strength, and color variation.
type ThreadBands struct {
	Null
	ThreadWidth      int
	CrossShadowDepth float64
	ColorVariance    float64
	LightThreadColor color.Color
	DarkThreadColor  color.Color
	Seed
}

func (t *ThreadBands) At(x, y int) color.Color {
	width := t.ThreadWidth
	if width <= 0 {
		width = 12
	}
	crossShadow := t.CrossShadowDepth
	if crossShadow < 0 {
		crossShadow = 0
	}
	colorVariance := t.ColorVariance
	if colorVariance < 0 {
		colorVariance = 0
	}

	warpIdx := int(math.Floor(float64(x) / float64(width)))
	weftIdx := int(math.Floor(float64(y) / float64(width)))
	localX := (x%width + width) % width
	localY := (y%width + width) % width

	light := toRGBAOrDefault(t.LightThreadColor, color.RGBA{214, 206, 191, 255})
	dark := toRGBAOrDefault(t.DarkThreadColor, color.RGBA{150, 141, 128, 255})

	warpColor := light
	if warpIdx%2 != 0 {
		warpColor = dark
	}
	weftColor := dark
	if weftIdx%2 != 0 {
		weftColor = light
	}

	warpOnTop := (warpIdx+weftIdx)%2 == 0

	var base color.RGBA
	var primaryPos, crossPos int
	if warpOnTop {
		base = warpColor
		primaryPos = localX
		crossPos = localY
	} else {
		base = weftColor
		primaryPos = localY
		crossPos = localX
	}

	profile := threadProfile(primaryPos, width)
	shadowProfile := threadProfile(crossPos, width)
	shadowFactor := clamp01(1 - crossShadow*(1-shadowProfile))
	fiberFactor := profile * shadowFactor

	varied := applyColorVariance(base, colorVariance, x, y, t.Seed.Seed)
	shaded := applyShade(varied, fiberFactor)

	return shaded
}

// NewThreadBands creates a thread weave pattern with sensible defaults.
func NewThreadBands(ops ...func(any)) image.Image {
	p := &ThreadBands{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		ThreadWidth:      12,
		CrossShadowDepth: 0.35,
		ColorVariance:    0.08,
		LightThreadColor: color.RGBA{214, 206, 191, 255},
		DarkThreadColor:  color.RGBA{150, 141, 128, 255},
		Seed:             Seed{Seed: 1},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoThreadBands produces a demo variant for readme.md pre-populated values.
func NewDemoThreadBands(ops ...func(any)) image.Image {
	return NewThreadBands(ops...)
}

// Options

type ThreadWidthOption struct{ ThreadWidth int }

func (t *ThreadWidthOption) SetThreadWidth(v int) { t.ThreadWidth = v }

type hasThreadWidth interface{ SetThreadWidth(int) }

func SetThreadWidth(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasThreadWidth); ok {
			h.SetThreadWidth(v)
		}
	}
}

type CrossShadow struct{ Depth float64 }

func (c *CrossShadow) SetCrossShadowDepth(v float64) { c.Depth = v }

type hasCrossShadowDepth interface{ SetCrossShadowDepth(float64) }

func SetCrossShadowDepth(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasCrossShadowDepth); ok {
			h.SetCrossShadowDepth(v)
		}
	}
}

type ThreadColorVariance struct{ Variance float64 }

func (c *ThreadColorVariance) SetColorVariance(v float64) { c.Variance = v }

type hasColorVariance interface{ SetColorVariance(float64) }

func SetColorVariance(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasColorVariance); ok {
			h.SetColorVariance(v)
		}
	}
}

type LightColor struct{ Color color.Color }

func (c *LightColor) SetLightThreadColor(col color.Color) { c.Color = col }

type hasLightThreadColor interface{ SetLightThreadColor(color.Color) }

func SetLightThreadColor(col color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasLightThreadColor); ok {
			h.SetLightThreadColor(col)
		}
	}
}

type DarkColor struct{ Color color.Color }

func (c *DarkColor) SetDarkThreadColor(col color.Color) { c.Color = col }

type hasDarkThreadColor interface{ SetDarkThreadColor(color.Color) }

func SetDarkThreadColor(col color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasDarkThreadColor); ok {
			h.SetDarkThreadColor(col)
		}
	}
}

// Option adapters for the ThreadBands pattern.
func (t *ThreadBands) SetThreadWidth(v int)              { t.ThreadWidth = v }
func (t *ThreadBands) SetCrossShadowDepth(v float64)     { t.CrossShadowDepth = v }
func (t *ThreadBands) SetColorVariance(v float64)        { t.ColorVariance = v }
func (t *ThreadBands) SetLightThreadColor(v color.Color) { t.LightThreadColor = v }
func (t *ThreadBands) SetDarkThreadColor(v color.Color)  { t.DarkThreadColor = v }

// Helpers

func threadProfile(pos, width int) float64 {
	if width <= 1 {
		return 1
	}
	half := float64(width-1) / 2
	dist := math.Abs(float64(pos)-half) / half
	curve := 1 - 0.35*dist
	return clamp01(curve)
}

func toRGBAOrDefault(c color.Color, fallback color.RGBA) color.RGBA {
	if c == nil {
		return fallback
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

func applyShade(c color.RGBA, factor float64) color.RGBA {
	f := clamp01(factor)
	return color.RGBA{
		R: clampUint8(float64(c.R) * f),
		G: clampUint8(float64(c.G) * f),
		B: clampUint8(float64(c.B) * f),
		A: c.A,
	}
}

func applyColorVariance(c color.RGBA, variance float64, x, y int, seed int64) color.RGBA {
	if variance <= 0 {
		return c
	}
	h := StableHash(x, y, uint64(seed))
	rv := ((float64(int64(h&0xff)) - 128) / 128.0) * variance
	gv := ((float64(int64(h>>8&0xff)) - 128) / 128.0) * variance
	bv := ((float64(int64(h>>16&0xff)) - 128) / 128.0) * variance

	return color.RGBA{
		R: clampUint8(float64(c.R) * (1 + rv)),
		G: clampUint8(float64(c.G) * (1 + gv)),
		B: clampUint8(float64(c.B) * (1 + bv)),
		A: c.A,
	}
}

func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(math.Round(v))
}
