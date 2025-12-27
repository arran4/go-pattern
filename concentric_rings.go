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
	Center
	Frequency
	Colors []color.Color
}

func (p *ConcentricRings) At(x, y int) color.Color {
	if len(p.Colors) == 0 {
		return color.Black
	}

	dx := float64(x - p.CenterX)
	dy := float64(y - p.CenterY)
	dist := math.Sqrt(dx*dx + dy*dy)

	// Apply frequency scaling if set
	if p.Frequency.Frequency != 0 {
		dist *= p.Frequency.Frequency
	}

	idx := int(math.Floor(dist)) % len(p.Colors)
	// Handle negative indices if any (unlikely with dist >= 0, but good practice)
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
	p.Frequency.Frequency = 1.0 // Default to 1.0 (identity scaling)
	for _, op := range ops {
		op(p)
	}
	return p
}
