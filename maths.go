package pattern

import (
	"image"
	"image/color"
)

// Ensure Maths implements the image.Image interface.
var _ image.Image = (*Maths)(nil)

// MathsFunc is the function signature for the pattern generator.
type MathsFunc func(x, y int) color.Color

// Maths is a pattern that generates colors based on a provided function.
type Maths struct {
	Null
	Func MathsFunc
}

func (m *Maths) At(x, y int) color.Color {
	if m.Func == nil {
		return color.RGBA{}
	}
	return m.Func(x, y)
}

// NewMaths creates a new Maths pattern with the given function.
func NewMaths(f MathsFunc, ops ...func(any)) image.Image {
	p := &Maths{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Func: f,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}
