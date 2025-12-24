package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure ScreenTone implements the image.Image interface.
var _ image.Image = (*ScreenTone)(nil)

// ScreenTone is a pattern that displays a grid of dots (halftone) at a specified angle.
type ScreenTone struct {
	Null
	Radius      // Size of the dots
	Spacing     // Distance between dots (frequency)
	Angle       // Angle of the grid in degrees
	FillColor
	SpaceColor
}

func (p *ScreenTone) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *ScreenTone) Bounds() image.Rectangle {
	return p.bounds
}

func (p *ScreenTone) At(x, y int) color.Color {
	spacing := float64(p.Spacing.Spacing)
	if spacing <= 0 {
		return p.SpaceColor.SpaceColor
	}

	// Convert angle to radians
	theta := p.Angle.Angle * math.Pi / 180.0
	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)

	// Rotate coordinates (inverse rotation to map screen to texture space)
	// We want to find which cell (u, v) corresponds to screen (x, y)
	xf := float64(x)
	yf := float64(y)

	u := xf*cosTheta + yf*sinTheta
	v := -xf*sinTheta + yf*cosTheta

	// Normalize coordinates to [0, spacing) within the rotated grid
	// Using math.Mod for floating point modulo
	du := math.Mod(u, spacing)
	if du < 0 {
		du += spacing
	}
	dv := math.Mod(v, spacing)
	if dv < 0 {
		dv += spacing
	}

	// Center of the cell
	c := spacing / 2.0

	// Distance from center
	distU := du - c
	distV := dv - c

	distSq := distU*distU + distV*distV
	radius := float64(p.Radius.Radius)
	radiusSq := radius * radius

	if distSq < radiusSq {
		return p.FillColor.FillColor
	}

	return p.SpaceColor.SpaceColor
}

// NewScreenTone creates a new ScreenTone pattern.
// Default Radius is 2.
// Default Spacing is 10.
// Default Angle is 45 degrees.
// Default FillColor is Black.
// Default SpaceColor is White.
func NewScreenTone(ops ...func(any)) image.Image {
	p := &ScreenTone{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	p.Angle.Angle = 45.0
	p.Radius.Radius = 2
	p.Spacing.Spacing = 10
	p.FillColor.FillColor = color.Black
	p.SpaceColor.SpaceColor = color.White

	for _, op := range ops {
		op(p)
	}
	return p
}
