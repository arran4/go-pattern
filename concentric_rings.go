package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure ConcentricRings implements the image.Image interface.
var _ image.Image = (*ConcentricRings)(nil)

// ConcentricRings generates concentric rings using sqrt(x^2 + y^2) % n.
type ConcentricRings struct {
	Null
	Colors []color.Color
}

func (p *ConcentricRings) At(x, y int) color.Color {
	if len(p.Colors) == 0 {
		return color.Black
	}
	// Calculate distance from origin (0,0).
	// Users might want to center this.
	// But the request says "(sqrt(x^2+y^2) % n)".
	// If the user uses a transform/crop/translate, they can center it.
	// Or we can add center coordinates. For now, strict adherence to formula.

	dist := math.Sqrt(float64(x*x + y*y))
	// n := float64(len(p.Colors))

	// The request says "bands with palette mapping".
	// Usually this means each band has width W, and we cycle through colors.
	// Or if the request literally means "modulo n" where n is palette size?
	// sqrt(x^2+y^2) is a float. Modulo n on a float?
	// Likely means floor(dist) % n.

	idx := int(math.Floor(dist)) % len(p.Colors)
	if idx < 0 {
		idx += len(p.Colors)
	}
	return p.Colors[idx]
}

// NewConcentricRings creates a new ConcentricRings pattern.
func NewConcentricRings(colors []color.Color, ops ...func(any)) image.Image {
	if len(colors) == 0 {
		colors = []color.Color{color.Black, color.White}
	}
	p := &ConcentricRings{
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
