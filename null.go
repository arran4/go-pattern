package pattern

import (
	"image"
	"image/color"
)

// Ensure Null implements the image.Image interface.
var _ image.Image = (*Null)(nil)

// Null is a pattern that returns a transparent color for all pixels.
type Null struct {
	bounds image.Rectangle
}

func (i *Null) ColorModel() color.Model {
	return color.RGBAModel
}

func (i *Null) Bounds() image.Rectangle {
	return i.bounds
}

func (i *Null) At(x, y int) color.Color {
	return color.RGBA{}
}

func (i *Null) SetBounds(bounds image.Rectangle) {
	i.bounds = bounds
}

func NewNull(ops ...func(any)) image.Image {
	p := &Null{
		bounds: image.Rect(0, 0, 255, 255),
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoNull produces a demo variant for readme.md pre-populated values
func NewDemoNull(ops ...func(any)) image.Image {
	return NewNull(ops...)
}
