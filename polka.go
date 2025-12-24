package pattern

import (
	"image"
	"image/color"
)

// Ensure Polka implements the image.Image interface.
var _ image.Image = (*Polka)(nil)

// Polka is a pattern that displays a grid of circles (polka dots).
type Polka struct {
	Null
	Radius
	Spacing
	FillColor
	SpaceColor
}

func (p *Polka) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *Polka) Bounds() image.Rectangle {
	return p.bounds
}

func (p *Polka) At(x, y int) color.Color {
	spacing := p.Spacing.Spacing
	if spacing <= 0 {
		return p.SpaceColor.SpaceColor
	}

	// Calculate center of the cell
	// We want the dot to be in the center of the square cell of size 'spacing'

	// Normalize coordinates to [0, spacing)
	// We use euclidean modulo to handle negative coordinates correctly
	dx := (x % spacing)
	if dx < 0 {
		dx += spacing
	}
	dy := (y % spacing)
	if dy < 0 {
		dy += spacing
	}

	// Center of the cell
	cx := spacing / 2
	cy := spacing / 2

	// Distance from center
	distX := dx - cx
	distY := dy - cy

	distSq := distX*distX + distY*distY
	radiusSq := p.Radius.Radius * p.Radius.Radius

	if distSq < radiusSq {
		return p.FillColor.FillColor
	}

	return p.SpaceColor.SpaceColor
}

// NewPolka creates a new Polka pattern.
// Default Radius is 10.
// Default Spacing is 40.
// Default FillColor is Black.
// Default SpaceColor is White.
func NewPolka(ops ...func(any)) image.Image {
	p := &Polka{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	p.Radius.Radius = 10
	p.Spacing.Spacing = 40
	p.FillColor.FillColor = color.Black
	p.SpaceColor.SpaceColor = color.White

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoPolka produces a demo variant for readme.md pre-populated values
func NewDemoPolka(ops ...func(any)) image.Image {
	return NewPolka(ops...)
}
