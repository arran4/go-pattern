package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure CrossHatch implements the image.Image interface.
var _ image.Image = (*CrossHatch)(nil)

// CrossHatch is a pattern that draws layered diagonal hatch lines.
type CrossHatch struct {
	Null
	SpaceSize
	LineSize
	LineColor
	SpaceColor
	LineImageSource
	Angles
}

func (p *CrossHatch) SetAngle(v float64) {
	p.Angles.Angles = []float64{v}
}

func (p *CrossHatch) At(x, y int) color.Color {
	ls := p.LineSize.LineSize
	ss := p.SpaceSize.SpaceSize
	period := float64(ls + ss)
	if period == 0 {
		return p.LineColor.LineColor
	}

	fx, fy := float64(x), float64(y)

	// Check if the point falls on any of the lines defined by the angles
	for _, angle := range p.Angles.Angles {
		// Convert angle to radians
		theta := angle * math.Pi / 180.0

		// Rotate coordinates: distance along the normal vector
		// d = x * cos(theta) + y * sin(theta)
		d := fx*math.Cos(theta) + fy*math.Sin(theta)

		// Modulo arithmetic for floating point
		// We want (d % period) but in a way that handles negatives correctly and is consistent

		// Use math.Floor to handle negative numbers correctly for modulo behavior
		// Python's % operator behavior: a % n = a - n * floor(a / n)
		m := d - period*math.Floor(d/period)

		// Check if within line width
		// We center the line or start from 0?
		// Usually 0 to LineSize is line, LineSize to Period is space.
		if m < float64(ls) {
			if p.LineImageSource.LineImageSource != nil {
				return p.LineImageSource.LineImageSource.At(x, y)
			}
			return p.LineColor.LineColor
		}
	}

	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
	}
	return color.RGBA{}
}

// NewCrossHatch creates a new CrossHatch pattern.
func NewCrossHatch(ops ...func(any)) image.Image {
	p := &CrossHatch{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.LineSize.LineSize = 1
	p.SpaceSize.SpaceSize = 5
	p.LineColor.LineColor = color.Black
	p.Angles.Angles = []float64{45} // Default to single hatch if not specified? Or maybe -45 too.

	for _, op := range ops {
		op(p)
	}

	// If the user provided single Angle via SetAngle (though CrossHatch doesn't embed Angle struct, but we could support it)
	// Currently CrossHatch embeds Angles.
	// If we want to support SetAngle as well, we should probably check if we want to add that interface.
	// But SetAngles covers it.

	return p
}

// NewDemoCrossHatch produces a demo variant for readme.md pre-populated values
func NewDemoCrossHatch(ops ...func(any)) image.Image {
	// Default demo with cross hatching
	ops = append([]func(any){
		SetAngles(45, -45),
		SetLineSize(2),
		SetSpaceSize(10),
	}, ops...)
	return NewCrossHatch(ops...)
}
