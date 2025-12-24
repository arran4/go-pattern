package pattern

import (
	"image"
	"image/color"
)

// Ensure Mirror implements the image.Image interface.
var _ image.Image = (*Mirror)(nil)

// Mirror is a pattern that mirrors (flips) the X and/or Y coordinates of an underlying image.
type Mirror struct {
	img        image.Image
	horizontal bool
	vertical   bool
}

func (m *Mirror) ColorModel() color.Model {
	return m.img.ColorModel()
}

func (m *Mirror) Bounds() image.Rectangle {
	return m.img.Bounds()
}

func (m *Mirror) At(x, y int) color.Color {
	b := m.img.Bounds()
	// If mirroring horizontally, we want to fetch the pixel from the opposite side.
	// The coordinate system is relative to the bounds.
	// If range is [Min.X, Max.X), width is Max.X - Min.X
	// x' = Max.X - 1 - (x - Min.X)
	srcX := x
	srcY := y

	if m.horizontal {
		srcX = b.Max.X - 1 - (x - b.Min.X)
	}
	if m.vertical {
		srcY = b.Max.Y - 1 - (y - b.Min.Y)
	}
	return m.img.At(srcX, srcY)
}

// NewMirror creates a new Mirror from an existing image.
// horizontal: flips left-right
// vertical: flips top-bottom
func NewMirror(img image.Image, horizontal, vertical bool, ops ...func(any)) image.Image {
	m := &Mirror{
		img:        img,
		horizontal: horizontal,
		vertical:   vertical,
	}
	for _, op := range ops {
		op(m)
	}
	return m
}
