package pattern

import (
	"image"
	"image/color"
)

// Ensure ModuloStripe implements the image.Image interface.
var _ image.Image = (*ModuloStripe)(nil)

// ModuloStripe generates a pattern based on (x + y) % n.
type ModuloStripe struct {
	Null
	Colors []color.Color
}

func (p *ModuloStripe) At(x, y int) color.Color {
	if len(p.Colors) == 0 {
		return color.Black
	}
	n := len(p.Colors)
	// (x + y) can be negative if bounds are negative?
	// x, y are usually coordinates.
	// We should handle modulo correctly for negative numbers if necessary, but pattern usually assumes positive space or wrapping.
	// Go's % operator returns negative result for negative dividend.
	// So we do ((x + y) % n + n) % n
	sum := x + y
	idx := sum % n
	if idx < 0 {
		idx += n
	}
	return p.Colors[idx]
}

// NewModuloStripe creates a new ModuloStripe pattern.
func NewModuloStripe(colors []color.Color, ops ...func(any)) image.Image {
	if len(colors) == 0 {
		colors = []color.Color{color.Black, color.White}
	}
	p := &ModuloStripe{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Colors: colors,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}
