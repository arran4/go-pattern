package pattern

import (
	"image"
	"image/color"
)

var _ image.Image = (*Padding)(nil)

type Padding struct {
	img      image.Image
	bounds   image.Rectangle
	bgPattern image.Image
	padding  int
}

func (p *Padding) ColorModel() color.Model {
	return p.img.ColorModel()
}

func (p *Padding) Bounds() image.Rectangle {
	return p.bounds
}

func (p *Padding) At(x, y int) color.Color {
	// Check if we are inside the inner image area
	// The inner image is centered or placed at padding offset?
	// Usually padding is added around.

	innerMinX := p.bounds.Min.X + p.padding
	innerMinY := p.bounds.Min.Y + p.padding
	innerMaxX := p.bounds.Max.X - p.padding
	innerMaxY := p.bounds.Max.Y - p.padding

	if x >= innerMinX && x < innerMaxX && y >= innerMinY && y < innerMaxY {
		// Map to inner image coordinates
		// Assume inner image starts at innerMinX, innerMinY in this view
		// So we need to map (x,y) to img relative.
		// If img.Bounds().Min is (0,0), then:
		// lx = x - innerMinX
		// ly = y - innerMinY

		lx := x - innerMinX + p.img.Bounds().Min.X
		ly := y - innerMinY + p.img.Bounds().Min.Y

		pt := image.Point{X: lx, Y: ly}
		if pt.In(p.img.Bounds()) {
			return p.img.At(lx, ly)
		}
	}

	// Outside inner area (or transparent part of inner image?), draw background
	if p.bgPattern != nil {
		return p.bgPattern.At(x, y)
	}
	return color.RGBA{}
}

func NewPadding(img image.Image, padding int, bgPattern image.Image) image.Image {
	b := img.Bounds()
	w := b.Dx() + 2*padding
	h := b.Dy() + 2*padding

	return &Padding{
		img:       img,
		bounds:    image.Rect(0, 0, w, h),
		bgPattern: bgPattern,
		padding:   padding,
	}
}
