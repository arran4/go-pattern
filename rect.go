package pattern

import (
	"image"
	"image/color"
)

// Ensure Rect implements the image.Image interface.
var _ image.Image = (*Rect)(nil)

// Rect is a pattern that draws a filled rectangle.
type Rect struct {
	Null
	FillColor
}

func (r *Rect) At(x, y int) color.Color {
	return r.FillColor.FillColor
}

// NewRect creates a new Rect pattern with the given options.
func NewRect(ops ...func(any)) image.Image {
	p := &Rect{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.FillColor.FillColor = color.Black

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoRect produces a demo variant for readme.md pre-populated values
func NewDemoRect(ops ...func(any)) image.Image {
	return NewRect(ops...)
}
