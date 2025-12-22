package pattern

import (
	"image"
	"image/color"
)

var _ image.Image = (*Tile)(nil)

type Tile struct {
	img    image.Image
	bounds image.Rectangle
}

func (t *Tile) ColorModel() color.Model {
	return t.img.ColorModel()
}

func (t *Tile) Bounds() image.Rectangle {
	return t.bounds
}

func (t *Tile) At(x, y int) color.Color {
	// Map x, y to source image coordinates by modulo
	b := t.img.Bounds()
	w := b.Dx()
	h := b.Dy()
	if w == 0 || h == 0 {
		return color.RGBA{}
	}

	// We need to handle negative coordinates if the requested bounds start before 0
	// Standard modulo in Go can return negative.

	// Shift x,y to be relative to 0,0 conceptually, then mod.
	// But actually, we just want the offset from Min.

	// Let's assume the source image starts at (0,0) for tiling purposes,
	// or we align the tile grid to (0,0).

	localX := x % w
	if localX < 0 {
		localX += w
	}
	localY := y % h
	if localY < 0 {
		localY += h
	}

	// Add source Min offset back
	return t.img.At(b.Min.X+localX, b.Min.Y+localY)
}

func NewTile(img image.Image, bounds image.Rectangle) image.Image {
	return &Tile{
		img:    img,
		bounds: bounds,
	}
}
