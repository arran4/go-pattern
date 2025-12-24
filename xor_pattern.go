package pattern

import (
	"image"
	"image/color"
)

// Ensure XorPattern implements the image.Image interface.
var _ image.Image = (*XorPattern)(nil)

// XorPattern generates an XOR texture pattern.
type XorPattern struct {
	Null
	Colors []color.Color
}

func (p *XorPattern) At(x, y int) color.Color {
	// Simple XOR pattern: (x ^ y)
	val := x ^ y

	if len(p.Colors) > 0 {
		idx := val % len(p.Colors)
		return p.Colors[idx]
	}

	v := uint8(val)
	return color.RGBA{R: v, G: v, B: v, A: 255}
}

// NewXorPattern creates a new XorPattern.
func NewXorPattern(ops ...func(any)) image.Image {
	p := &XorPattern{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}
