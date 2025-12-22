package pattern

import (
	"image"
	"image/color"
)

// Ensure Crop implements the image.Image interface.
var _ image.Image = (*Crop)(nil)

type Crop struct {
	img  image.Image
	rect image.Rectangle
}

func (c *Crop) ColorModel() color.Model {
	return c.img.ColorModel()
}

func (c *Crop) Bounds() image.Rectangle {
	return c.rect
}

func (c *Crop) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(c.rect)) {
		return color.RGBA{}
	}
	return c.img.At(x, y)
}

func NewCrop(img image.Image, rect image.Rectangle) image.Image {
	// Intersection of requested crop and actual image bounds
	// intersect := rect.Intersect(img.Bounds())
	// Usually crop implies we want the result to have these bounds.
	// But if the user asks for a crop outside the image, it should probably be transparent?
	// Or should we just limit it?
	// Standard sub-image behavior is usually limited to bounds.

	// However, if we want to "Crop" to a specific window, we might mean "View" this window.
	// If the window is larger, we get transparency.
	// Let's stick to the requested rect as bounds. At() will handle out-of-bounds.

	return &Crop{
		img:  img,
		rect: rect,
	}
}
