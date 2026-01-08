package pattern

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Ensure GlyphRing implements the image.Image interface.
var _ image.Image = (*GlyphRing)(nil)

// GlyphRing renders a circular band with glyphs stamped around its circumference and a soft outer glow.
// Parameters of interest:
//   - Radius (via SetRadius): controls the ring radius relative to the bounds.
//   - Density (via SetDensity): controls how many runes appear around the circle.
//   - GlowColor (via SetGlowColor): sets the tint of the exterior glow.
//
// Additional styling controls (rune, band, and background colors plus sizing) are provided via setters.
type GlyphRing struct {
	Null
	Radius
	Density
	Seed
	GlowColor     color.Color
	RuneColor     color.Color
	BandColor     color.Color
	Background    color.Color
	BandThickness int
	GlowSize      int
	GlyphSize     int
	rendered      *image.RGBA
}

// SetGlowColor sets the outer glow color.
func (p *GlyphRing) SetGlowColor(v color.Color) { p.GlowColor = v }

// SetRuneColor sets the glyph fill color.
func (p *GlyphRing) SetRuneColor(v color.Color) { p.RuneColor = v }

// SetBandColor sets the ring body color.
func (p *GlyphRing) SetBandColor(v color.Color) { p.BandColor = v }

// SetBackgroundColor sets the background color.
func (p *GlyphRing) SetBackgroundColor(v color.Color) { p.Background = v }

// SetBandThickness adjusts the ring thickness.
func (p *GlyphRing) SetBandThickness(v int) { p.BandThickness = v }

// SetGlowSize adjusts the glow falloff width.
func (p *GlyphRing) SetGlowSize(v int) { p.GlowSize = v }

// SetGlyphSize adjusts the glyph scale.
func (p *GlyphRing) SetGlyphSize(v int) { p.GlyphSize = v }

type hasGlowColor interface{ SetGlowColor(color.Color) }
type hasRuneColor interface{ SetRuneColor(color.Color) }
type hasBandColor interface{ SetBandColor(color.Color) }
type hasBackgroundColor interface{ SetBackgroundColor(color.Color) }
type hasBandThickness interface{ SetBandThickness(int) }
type hasGlowSize interface{ SetGlowSize(int) }
type hasGlyphSize interface{ SetGlyphSize(int) }

// SetGlowColor creates an option to set the glow color.
func SetGlowColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasGlowColor); ok {
			h.SetGlowColor(v)
		}
	}
}

// SetRuneColor creates an option to set the rune color.
func SetRuneColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasRuneColor); ok {
			h.SetRuneColor(v)
		}
	}
}

// SetBandColor creates an option to set the band color.
func SetBandColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasBandColor); ok {
			h.SetBandColor(v)
		}
	}
}

// SetBackgroundColor creates an option to set the background color.
func SetBackgroundColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasBackgroundColor); ok {
			h.SetBackgroundColor(v)
		}
	}
}

// SetBandThickness creates an option to set the ring thickness.
func SetBandThickness(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasBandThickness); ok {
			h.SetBandThickness(v)
		}
	}
}

// SetGlowSize creates an option to set the glow falloff width.
func SetGlowSize(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasGlowSize); ok {
			h.SetGlowSize(v)
		}
	}
}

// SetGlyphSize creates an option to set the glyph scale.
func SetGlyphSize(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasGlyphSize); ok {
			h.SetGlyphSize(v)
		}
	}
}

func (p *GlyphRing) ColorModel() color.Model {
	return color.NRGBAModel
}

func (p *GlyphRing) Bounds() image.Rectangle {
	return p.bounds
}

func (p *GlyphRing) At(x, y int) color.Color {
	return p.rendered.At(x, y)
}

// NewGlyphRing constructs a GlyphRing with optional configuration.
func NewGlyphRing(ops ...func(any)) image.Image {
	p := &GlyphRing{
		Null: Null{
			bounds: image.Rect(0, 0, 400, 400),
		},
		GlowColor:     color.NRGBA{R: 80, G: 190, B: 255, A: 200},
		RuneColor:     color.NRGBA{R: 40, G: 255, B: 210, A: 255},
		BandColor:     color.NRGBA{R: 210, G: 200, B: 120, A: 255},
		Background:    color.NRGBA{R: 8, G: 6, B: 16, A: 255},
		BandThickness: 18,
		GlowSize:      32,
		GlyphSize:     14,
	}
	p.Density.Density = 0.06
	p.Seed.Seed = 99

	for _, op := range ops {
		op(p)
	}

	p.render()
	return p
}

// NewDemoGlyphRing produces a demo variant for readme.md pre-populated values.
func NewDemoGlyphRing(ops ...func(any)) image.Image {
	return NewGlyphRing(ops...)
}

func (p *GlyphRing) render() {
	b := p.Bounds()
	dst := image.NewRGBA(b)

	// Normalize sizing defaults against the current bounds.
	minDim := b.Dx()
	if b.Dy() < minDim {
		minDim = b.Dy()
	}

	radius := p.Radius.Radius
	if radius == 0 {
		radius = minDim/2 - int(float64(p.GlowSize)*0.8)
	}
	if radius < 8 {
		radius = 8
	}

	if p.BandThickness <= 0 {
		p.BandThickness = int(math.Max(6, float64(minDim)/30))
	}

	if p.GlowSize <= 0 {
		p.GlowSize = p.BandThickness * 2
	}

	if p.GlyphSize <= 0 {
		p.GlyphSize = p.BandThickness
	}

	if p.Density.Density <= 0 {
		p.Density.Density = 0.05
	}

	cx := float64(b.Min.X + b.Dx()/2)
	cy := float64(b.Min.Y + b.Dy()/2)

	inner := float64(radius) - float64(p.BandThickness)/2
	outer := float64(radius) + float64(p.BandThickness)/2
	glowOuter := outer + float64(p.GlowSize)

	bg := colorToNRGBA(p.Background)
	ring := colorToNRGBA(p.BandColor)
	glow := colorToNRGBA(p.GlowColor)

	draw.Draw(dst, b, image.NewUniform(bg), image.Point{}, draw.Src)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			dx := float64(x) - cx + 0.5
			dy := float64(y) - cy + 0.5
			dist := math.Hypot(dx, dy)

			col := bg
			if dist >= inner && dist <= outer {
				col = ring
			}

			if dist > outer && dist < glowOuter {
				falloff := 1 - (dist-outer)/float64(p.GlowSize)
				if falloff < 0 {
					falloff = 0
				}
				col = addScaled(col, glow, falloff*falloff)
			}

			dst.Set(x, y, col)
		}
	}

	p.renderGlyphs(dst, radius, cx, cy)
	p.rendered = dst
}

func (p *GlyphRing) renderGlyphs(dst *image.RGBA, radius int, cx, cy float64) {
	circumference := float64(radius) * 2 * math.Pi
	glyphCount := int(math.Max(4, circumference*p.Density.Density))
	glyphColor := colorToNRGBA(p.RuneColor)

	for i := 0; i < glyphCount; i++ {
		h := StableHash(i, radius, uint64(p.Seed.Seed))
		glyph := runeGlyphs[int(h%uint64(len(runeGlyphs)))]

		baseAngle := float64(i) * 2 * math.Pi / float64(glyphCount)
		jitter := (float64((h>>16)&0xFFFF)/65535.0 - 0.5) * (2 * math.Pi / float64(glyphCount)) * 0.5
		angle := baseAngle + jitter

		p.placeGlyph(dst, glyph, glyphColor, radius, angle, cx, cy)
	}
}

func (p *GlyphRing) placeGlyph(dst *image.RGBA, glyph []string, glyphColor color.NRGBA, radius int, angle, cx, cy float64) {
	if len(glyph) == 0 {
		return
	}

	gw := len(glyph[0])
	gh := len(glyph)
	scale := float64(p.GlyphSize) / float64(gw)

	sinA, cosA := math.Sincos(angle + math.Pi/2)
	baseX := cx + math.Cos(angle)*float64(radius)
	baseY := cy + math.Sin(angle)*float64(radius)

	for y := 0; y < gh; y++ {
		row := glyph[y]
		for x := 0; x < gw; x++ {
			if row[x] != '#' {
				continue
			}

			localX := (float64(x) - float64(gw)/2 + 0.5) * scale
			localY := (float64(y) - float64(gh)/2 + 0.5) * scale

			rx := localX*cosA - localY*sinA
			ry := localX*sinA + localY*cosA

			px := int(math.Round(baseX + rx))
			py := int(math.Round(baseY + ry))
			if !(image.Point{px, py}.In(dst.Bounds())) {
				continue
			}

			existing := colorToNRGBA(dst.RGBAAt(px, py))
			dst.Set(px, py, overNRGBA(existing, glyphColor))
		}
	}
}

func colorToNRGBA(c color.Color) color.NRGBA {
	if n, ok := c.(color.NRGBA); ok {
		return n
	}
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

func addScaled(base, add color.NRGBA, t float64) color.NRGBA {
	clamp := func(v float64) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v + 0.5)
	}
	return color.NRGBA{
		R: clamp(float64(base.R) + float64(add.R)*t),
		G: clamp(float64(base.G) + float64(add.G)*t),
		B: clamp(float64(base.B) + float64(add.B)*t),
		A: clamp(float64(base.A)),
	}
}

func overNRGBA(dst, src color.NRGBA) color.NRGBA {
	a := float64(src.A) / 255.0
	inv := 1 - a
	blend := func(s, d uint8) uint8 {
		return uint8(math.Min(255, math.Round(float64(s)*a+float64(d)*inv)))
	}
	return color.NRGBA{
		R: blend(src.R, dst.R),
		G: blend(src.G, dst.G),
		B: blend(src.B, dst.B),
		A: uint8(math.Min(255, math.Round(float64(src.A)+float64(dst.A)*inv))),
	}
}

var runeGlyphs = [][]string{
	{
		"..#..",
		".###.",
		"..#..",
		"..#..",
		".###.",
	},
	{
		"#...#",
		".#.#.",
		"..#..",
		".#.#.",
		"#...#",
	},
	{
		".###.",
		"#...#",
		"..#..",
		"#...#",
		".###.",
	},
	{
		"#.#.#",
		".#.#.",
		"..#..",
		".#.#.",
		"#.#.#",
	},
	{
		"..#..",
		".#.#.",
		"#.#.#",
		".#.#.",
		"..#..",
	},
}
