package pattern

import (
	"image"
	"image/color"
)

// Ensure Circle implements the image.Image interface.
var _ image.Image = (*Circle)(nil)

// Circle is a pattern that draws a circle fitting within its bounds.
// It supports a border (LineSize, LineColor, LineImageSource) and a fill (FillColor, FillImageSource).
type Circle struct {
	Null
	LineSize
	LineColor
	LineImageSource
	FillColor
	FillImageSource
	SpaceColor
}

func (p *Circle) At(x, y int) color.Color {
	b := p.Bounds()
	width := b.Dx()
	height := b.Dy()

	// Center coordinate in 2x space:
	cx2 := b.Min.X + b.Max.X
	cy2 := b.Min.Y + b.Max.Y

	// Point (x, y) center in 2x space is (2*x+1, 2*y+1).
	dx2 := (2*x + 1) - cx2
	dy2 := (2*y + 1) - cy2

	// Distance squared in 2x space.
	distSq := dx2*dx2 + dy2*dy2

	diameter := width
	if height < diameter {
		diameter = height
	}

	// Outer radius squared in 2x space.
	// radius = diameter/2.
	// In 2x space, radius2x = diameter.
	// radiusSq = diameter^2.
	outerRadiusSq := diameter * diameter

	if distSq > outerRadiusSq {
		// Outside
		if p.SpaceColor.SpaceColor != nil {
			return p.SpaceColor.SpaceColor
		}
		return color.RGBA{}
	}

	// Inside Outer Circle.

	// Check for Border.
	ls := p.LineSize.LineSize
	if ls > 0 {
		// Border Logic.
		// Inner Radius = Radius - LineSize.
		// In 2x space: innerRadius2x = diameter - 2*LineSize.
		innerDiameter := diameter - 2*ls
		if innerDiameter < 0 {
			innerDiameter = 0
		}
		innerRadiusSq := innerDiameter * innerDiameter

		if distSq > innerRadiusSq {
			// In the Border.
			if p.LineImageSource.LineImageSource != nil {
				return p.LineImageSource.LineImageSource.At(x, y)
			}
			return p.LineColor.LineColor
		}

		// Inside Inner Circle (Fill).
		if p.FillImageSource.FillImageSource != nil {
			return p.FillImageSource.FillImageSource.At(x, y)
		}
		if p.FillColor.FillColor != nil {
			return p.FillColor.FillColor
		}
		// If no fill specified, transparent?
		return color.RGBA{}
	}

	// LineSize == 0. Legacy/Solid mode.
	// The entire circle is filled.

	// Prioritize: FillImage > FillColor > LineImage > LineColor
	if p.FillImageSource.FillImageSource != nil {
		return p.FillImageSource.FillImageSource.At(x, y)
	}
	if p.FillColor.FillColor != nil {
		return p.FillColor.FillColor
	}
	if p.LineImageSource.LineImageSource != nil {
		return p.LineImageSource.LineImageSource.At(x, y)
	}
	return p.LineColor.LineColor
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
	p.LineSize.LineSize = 0
	// FillColor, FillImage, LineImage, SpaceColor default to nil (transparent/zero)

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoCircle produces a demo variant for readme.md pre-populated values
func NewDemoCircle(ops ...func(any)) image.Image {
	return NewCircle(ops...)
}
