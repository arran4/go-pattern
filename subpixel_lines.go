package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure SubpixelLines implements the image.Image interface.
var _ image.Image = (*SubpixelLines)(nil)

// SubpixelLines renders alternating dark/light horizontal bands with subtle RGB
// channel offsets and a vignette falloff.
type SubpixelLines struct {
	Null
	LineThickness  int
	OffsetStrength float64
	VignetteRadius float64
}

func (p *SubpixelLines) At(x, y int) color.Color {
	thickness := p.LineThickness
	if thickness <= 0 {
		thickness = 1
	}

	rSample := p.sampleLine(float64(y)+p.OffsetStrength, thickness)
	gSample := p.sampleLine(float64(y), thickness)
	bSample := p.sampleLine(float64(y)-p.OffsetStrength, thickness)

	phase := math.Mod(float64(x), 3.0) / 3.0
	rWeight := 0.9 + 0.1*math.Cos(phase*math.Pi*2.0)
	gWeight := 0.9 + 0.08*math.Cos((phase+1.0/3.0)*math.Pi*2.0)
	bWeight := 0.9 + 0.1*math.Cos((phase+2.0/3.0)*math.Pi*2.0)

	r := clamp01(rSample * rWeight)
	g := clamp01(gSample * gWeight)
	b := clamp01(bSample * bWeight)

	v := p.vignette(float64(x), float64(y))
	r *= v
	g *= v
	b *= v

	return color.NRGBA{
		R: uint8(r*255 + 0.5),
		G: uint8(g*255 + 0.5),
		B: uint8(b*255 + 0.5),
		A: 255,
	}
}

func (p *SubpixelLines) sampleLine(y float64, thickness int) float64 {
	period := float64(thickness * 2)
	if period <= 0 {
		return 0
	}

	pos := math.Mod(y, period)
	if pos < 0 {
		pos += period
	}

	light := 0.88
	dark := 0.18

	var base float64
	if pos < float64(thickness) {
		base = light
	} else {
		base = dark
	}

	edge := math.Min(pos, period-pos) / float64(thickness)
	soften := 1.0 - smoothstep(0.0, 0.6, edge)

	return clamp01(base - soften*0.08)
}

func (p *SubpixelLines) vignette(x, y float64) float64 {
	if p.VignetteRadius <= 0 {
		return 1
	}

	b := p.Bounds()
	cx := float64(b.Min.X) + float64(b.Dx())/2.0
	cy := float64(b.Min.Y) + float64(b.Dy())/2.0

	dx := x - cx
	dy := y - cy
	dist := math.Sqrt(dx*dx + dy*dy)

	maxDim := math.Min(float64(b.Dx()), float64(b.Dy())) / 2.0
	if maxDim == 0 {
		return 1
	}

	radius := p.VignetteRadius * maxDim
	if radius <= 0 {
		return 1
	}

	outer := radius + maxDim*0.25
	fade := smoothstep(radius, outer, dist)

	return clamp01(1.0 - fade)
}


func smoothstep(edge0, edge1, x float64) float64 {
	if edge1 == edge0 {
		return 0
	}
	t := (x - edge0) / (edge1 - edge0)
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	return t * t * (3 - 2*t)
}

// NewSubpixelLines creates a new SubpixelLines pattern.
func NewSubpixelLines(ops ...func(any)) image.Image {
	p := &SubpixelLines{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		LineThickness:  2,
		OffsetStrength: 0.75,
		VignetteRadius: 0.85,
	}

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoSubpixelLines produces a demo variant for readme.md pre-populated values.
func NewDemoSubpixelLines(ops ...func(any)) image.Image {
	return NewSubpixelLines(ops...)
}

func (p *SubpixelLines) SetLineThickness(thickness int) {
	p.LineThickness = thickness
}

type hasLineThickness interface {
	SetLineThickness(int)
}

// SetLineThickness configures the thickness of each horizontal band.
func SetLineThickness(thickness int) func(any) {
	return func(i any) {
		if v, ok := i.(hasLineThickness); ok {
			v.SetLineThickness(thickness)
		}
	}
}

func (p *SubpixelLines) SetOffsetStrength(strength float64) {
	p.OffsetStrength = strength
}

type hasOffsetStrength interface {
	SetOffsetStrength(float64)
}

// SetOffsetStrength controls the per-channel vertical offset used when sampling
// the alternating lines.
func SetOffsetStrength(strength float64) func(any) {
	return func(i any) {
		if v, ok := i.(hasOffsetStrength); ok {
			v.SetOffsetStrength(strength)
		}
	}
}

func (p *SubpixelLines) SetVignetteRadius(radius float64) {
	p.VignetteRadius = radius
}

type hasVignetteRadius interface {
	SetVignetteRadius(float64)
}

// SetVignetteRadius adjusts the radius (in normalized 0..1 units of the
// shortest canvas dimension) where the vignette falloff begins.
func SetVignetteRadius(radius float64) func(any) {
	return func(i any) {
		if v, ok := i.(hasVignetteRadius); ok {
			v.SetVignetteRadius(radius)
		}
	}
}
