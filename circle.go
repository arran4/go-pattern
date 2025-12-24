package pattern

import (
	"image"
	"image/color"
)

// Ensure Circle implements the image.Image interface.
var _ image.Image = (*Circle)(nil)

// Circle is a pattern that draws a circle fitting within its bounds.
type Circle struct {
	Null
	LineColor
	SpaceColor
	LineImageSource
}

func (p *Circle) At(x, y int) color.Color {
	b := p.Bounds()
	width := b.Dx()
	height := b.Dy()

	// Calculate diameter and radius (squared)
	// We use 2*coordinate logic to avoid floating point math.
	// Center coordinate in 2x space:
	cx2 := b.Min.X + b.Max.X
	cy2 := b.Min.Y + b.Max.Y

	// Radius in 2x space would be diameter (since radius = diameter/2, 2*radius = diameter).
	// We want to check distance from center.
	// Point (x, y) center in 2x space is (2*x+1, 2*y+1).

	dx2 := (2*x + 1) - cx2
	dy2 := (2*y + 1) - cy2

	// Max allowed distance (squared) in 2x space.
	// The circle diameter corresponds to min(width, height).
	// The radius in 1x space is min(width, height)/2.
	// The radius in 2x space is min(width, height).
	// So we compare dx2^2 + dy2^2 <= diameter^2.

	diameter := width
	if height < diameter {
		diameter = height
	}

	radiusSq := diameter * diameter

	if dx2*dx2 + dy2*dy2 <= radiusSq {
		if p.LineImageSource.LineImageSource != nil {
			return p.LineImageSource.LineImageSource.At(x, y)
		}
		return p.LineColor.LineColor
	}

	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
	}
	return color.RGBA{}
}

// NewCircle creates a new Circle pattern.
func NewCircle(ops ...func(any)) image.Image {
	p := &Circle{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.LineColor.LineColor = color.Black
	// SpaceColor defaults to nil (transparent)

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoCircle produces a demo variant for readme.md pre-populated values
func NewDemoCircle(ops ...func(any)) image.Image {
	return NewCircle(ops...)
}
