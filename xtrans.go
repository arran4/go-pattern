package pattern

import (
	"image"
	"image/color"
)

// Ensure XTrans implements the image.Image interface.
var _ image.Image = (*XTrans)(nil)

// XTrans applies a Fujifilm X-Trans color filter array pattern to an image.
type XTrans struct {
	Null
	Input image.Image
}

// NewXTrans creates a new XTrans pattern.
// If input is nil, it renders the raw colors of the X-Trans color filter array.
func NewXTrans(input image.Image, ops ...func(any)) image.Image {
	p := &XTrans{
		Input: input,
		Null: Null{
			bounds: image.Rect(0, 0, 100, 100),
		},
	}
	if input != nil {
		p.bounds = input.Bounds()
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// The standard X-Trans CFA 6x6 pattern.
// 0 = Red, 1 = Green, 2 = Blue.
var xtransPattern = [6][6]uint8{
	{1, 2, 1, 1, 0, 1},
	{0, 1, 0, 2, 1, 2},
	{1, 2, 1, 1, 0, 1},
	{1, 0, 1, 1, 2, 1},
	{2, 1, 2, 0, 1, 0},
	{1, 0, 1, 1, 2, 1},
}

// At returns the color of the pixel at (x, y) based on the X-Trans pattern.
func (p *XTrans) At(x, y int) color.Color {
	mx := x % 6
	if mx < 0 {
		mx += 6
	}
	my := y % 6
	if my < 0 {
		my += 6
	}

	cfa := xtransPattern[my][mx]

	if p.Input == nil {
		switch cfa {
		case 0:
			return color.RGBA{255, 0, 0, 255}
		case 1:
			return color.RGBA{0, 255, 0, 255}
		case 2:
			return color.RGBA{0, 0, 255, 255}
		}
	}

	c := p.Input.At(x, y)
	r, g, b, a := c.RGBA()

	switch cfa {
	case 0:
		return color.RGBA{uint8(r >> 8), 0, 0, uint8(a >> 8)}
	case 1:
		return color.RGBA{0, uint8(g >> 8), 0, uint8(a >> 8)}
	case 2:
		return color.RGBA{0, 0, uint8(b >> 8), uint8(a >> 8)}
	}
	return color.Transparent
}
