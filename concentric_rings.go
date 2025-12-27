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
	FrequencyX
	FrequencyY
	Colors []color.Color
}

func (p *ConcentricRings) At(x, y int) color.Color {
	if len(p.Colors) == 0 {
		return color.Black
	}

	dx := float64(x - p.CenterX)
	dy := float64(y - p.CenterY)

	// Apply separate frequencies if set, otherwise use global Frequency
	fx := p.Frequency.Frequency
	if p.FrequencyX.FrequencyX != 0 {
		fx = p.FrequencyX.FrequencyX
	}
	fy := p.Frequency.Frequency
	if p.FrequencyY.FrequencyY != 0 {
		fy = p.FrequencyY.FrequencyY
	}

	// Calculate distance with scaling
	// sqrt((dx*fx)^2 + (dy*fy)^2)
	dist := math.Sqrt((dx*fx)*(dx*fx) + (dy*fy)*(dy*fy))

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
