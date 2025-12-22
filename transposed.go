package pattern

import (
	"image"
	"image/color"
)

// Ensure Transposed implements the image.Image interface.
var _ image.Image = (*Transposed)(nil)

// Transposed is a pattern that transposes the X and Y coordinates of an underlying image.
type Transposed struct {
	img image.Image
	x   int
	y   int
}

func (t *Transposed) ColorModel() color.Model {
	return t.img.ColorModel()
}

func (t *Transposed) Bounds() image.Rectangle {
	b := t.img.Bounds()
	return image.Rect(b.Min.Y, b.Min.X, b.Max.Y, b.Max.X)
}

func (t *Transposed) At(x, y int) color.Color {
	return t.img.At(y+t.y, x+t.x)
}

// NewTransposed creates a new Transposed from an existing image.
func NewTransposed(img image.Image, x, y int, ops ...func(any)) image.Image {
	t := &Transposed{
		img: img,
		x:   x,
		y:   y,
	}
	for _, op := range ops {
		op(t)
	}
	return t
}

// NewDemoTransposed produces a demo variant for readme.md pre-populated values
func NewDemoTransposed(ops ...func(any)) image.Image {
	return NewTransposed(NewSimpleZoom(NewChecker(color.Black, color.White, ops...), 20, ops...), 10, 10, ops...)
}
