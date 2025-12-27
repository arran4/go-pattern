package pattern

import (
	"image"
	"image/color"
)

// Ensure Blend implements the image.Image interface.
var _ image.Image = (*Blend)(nil)

type BlendMode int

const (
	BlendAdd BlendMode = iota
	BlendMultiply
	BlendAverage
	BlendScreen
	BlendOverlay
)

// Blend combines two images using a specified blend mode.
type Blend struct {
	Null
	Image1 image.Image
	Image2 image.Image
	Mode   BlendMode
}

func (b *Blend) At(x, y int) color.Color {
	c1 := b.Image1.At(x, y)
	c2 := b.Image2.At(x, y)

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	// Normalize to [0, 1] float
	fr1, fg1, fb1, fa1 := float64(r1)/65535, float64(g1)/65535, float64(b1)/65535, float64(a1)/65535
	fr2, fg2, fb2, fa2 := float64(r2)/65535, float64(g2)/65535, float64(b2)/65535, float64(a2)/65535

	var or, og, ob, oa float64
	oa = (fa1 + fa2) / 2 // Simple alpha blending for now, or max? Let's assume opaque for texture gen.
	if oa > 1 {
		oa = 1
	}

	switch b.Mode {
	case BlendAdd:
		or = fr1 + fr2
		og = fg1 + fg2
		ob = fb1 + fb2
	case BlendMultiply:
		or = fr1 * fr2
		og = fg1 * fg2
		ob = fb1 * fb2
	case BlendAverage:
		or = (fr1 + fr2) / 2
		og = (fg1 + fg2) / 2
		ob = (fb1 + fb2) / 2
	case BlendScreen:
		or = 1 - (1-fr1)*(1-fr2)
		og = 1 - (1-fg1)*(1-fg2)
		ob = 1 - (1-fb1)*(1-fb2)
	case BlendOverlay:
		or = overlay(fr1, fr2)
		og = overlay(fg1, fg2)
		ob = overlay(fb1, fb2)
	default:
		or, og, ob = fr1, fg1, fb1
	}

	// Clamp
	if or > 1 { or = 1 }
	if og > 1 { og = 1 }
	if ob > 1 { ob = 1 }

	return color.RGBA64{
		R: uint16(or * 65535),
		G: uint16(og * 65535),
		B: uint16(ob * 65535),
		A: uint16(oa * 65535),
	}
}

func overlay(a, b float64) float64 {
	if a < 0.5 {
		return 2 * a * b
	}
	return 1 - 2*(1-a)*(1-b)
}

// NewBlend creates a new Blend pattern.
func NewBlend(i1, i2 image.Image, mode BlendMode, ops ...func(any)) image.Image {
	p := &Blend{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Image1: i1,
		Image2: i2,
		Mode:   mode,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}
