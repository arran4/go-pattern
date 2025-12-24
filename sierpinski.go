package pattern

import (
	"image"
	"image/color"
)

// Ensure implementation of image.Image
var _ image.Image = (*SierpinskiTriangle)(nil)
var _ image.Image = (*SierpinskiCarpet)(nil)

// SierpinskiTriangle represents a pattern generated using the Sierpinski Triangle fractal algorithm.
// It uses a bitwise operation equivalent to Pascal's Triangle modulo 2.
type SierpinskiTriangle struct {
	Null
	FillColor
	SpaceColor
}

// At returns the color at the given coordinates.
// If (x & y) == 0, it returns the fill color (part of the triangle).
// Otherwise, it returns the space color.
func (p *SierpinskiTriangle) At(x, y int) color.Color {
	if x < 0 || y < 0 {
		return p.SpaceColor.SpaceColor
	}
	if (x & y) == 0 {
		return p.FillColor.FillColor
	}
	return p.SpaceColor.SpaceColor
}

// NewSierpinskiTriangle creates a new SierpinskiTriangle pattern.
func NewSierpinskiTriangle(ops ...func(any)) image.Image {
	p := &SierpinskiTriangle{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.FillColor.FillColor = color.Black
	p.SpaceColor.SpaceColor = color.RGBA{0, 0, 0, 0} // Transparent

	for _, op := range ops {
		op(p)
	}
	return p
}

// SierpinskiCarpet represents a pattern generated using the Sierpinski Carpet fractal algorithm.
type SierpinskiCarpet struct {
	Null
	FillColor
	SpaceColor
}

// At returns the color at the given coordinates.
// It recursively checks if the coordinate belongs to a "hole" in the carpet.
func (p *SierpinskiCarpet) At(x, y int) color.Color {
	if x < 0 || y < 0 {
		return p.SpaceColor.SpaceColor
	}
	tx, ty := x, y
	for tx > 0 || ty > 0 {
		if tx%3 == 1 && ty%3 == 1 {
			return p.SpaceColor.SpaceColor
		}
		tx /= 3
		ty /= 3
	}
	return p.FillColor.FillColor
}

// NewSierpinskiCarpet creates a new SierpinskiCarpet pattern.
func NewSierpinskiCarpet(ops ...func(any)) image.Image {
	p := &SierpinskiCarpet{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.FillColor.FillColor = color.Black
	p.SpaceColor.SpaceColor = color.RGBA{0, 0, 0, 0} // Transparent

	for _, op := range ops {
		op(p)
	}
	return p
}
