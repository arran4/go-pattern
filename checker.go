package pattern

import (
	"image"
	"image/color"
)

// Ensure Checker implements the image.Image interface.
var _ image.Image = (*Checker)(nil)

// Checker is a pattern that alternates between two colors in a checkerboard fashion.
type Checker struct {
	Null
	color1, color2 color.Color
}

func (c *Checker) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *Checker) Bounds() image.Rectangle {
	return c.bounds
}

func (c *Checker) At(x, y int) color.Color {
	if x%2 == y%2 {
		return c.color1
	}
	return c.color2
}

// NewChecker creates a new Checker with the given colors and square size.
func NewChecker(color1, color2 color.Color, ops ...func(any)) image.Image {
	p := &Checker{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		color1: color1,
		color2: color2,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoChecker produces a demo variant for readme.md pre-populated values
func NewDemoChecker(ops ...func(any)) image.Image {
	return NewChecker(color.Black, color.White, ops...)
}
