package pattern

import (
	"image"
	"image/color"
)

// Ensure Clamp implements the image.Image interface.
var _ image.Image = (*Clamp)(nil)

// Clamp extends the edges of the source image to fill the specified bounds.
// If a pixel is outside the source image, the closest edge pixel is used.
type Clamp struct {
	img    image.Image
	bounds image.Rectangle
}

func (c *Clamp) ColorModel() color.Model {
	if c.img != nil {
		return c.img.ColorModel()
	}
	return color.RGBAModel
}

func (c *Clamp) Bounds() image.Rectangle {
	return c.bounds
}

func (c *Clamp) At(x, y int) color.Color {
	if c.img == nil {
		return color.Transparent
	}

	b := c.img.Bounds()

	// If the requested point is inside, return directly
	if x >= b.Min.X && x < b.Max.X && y >= b.Min.Y && y < b.Max.Y {
		return c.img.At(x, y)
	}

	// Clamp x
	cx := x
	if cx < b.Min.X {
		cx = b.Min.X
	} else if cx >= b.Max.X {
		cx = b.Max.X - 1
	}

	// Clamp y
	cy := y
	if cy < b.Min.Y {
		cy = b.Min.Y
	} else if cy >= b.Max.Y {
		cy = b.Max.Y - 1
	}

	return c.img.At(cx, cy)
}

// NewClamp creates a new Clamp pattern.
// bounds defines the new size of the image.
func NewClamp(img image.Image, bounds image.Rectangle, ops ...func(any)) image.Image {
	c := &Clamp{
		img:    img,
		bounds: bounds,
	}
	for _, op := range ops {
		op(c)
	}
	return c
}

// SetBounds sets the bounds of the image.
func (c *Clamp) SetBounds(bounds image.Rectangle) {
	c.bounds = bounds
}
