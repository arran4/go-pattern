package pattern

import (
	"image"
	"image/color"
)

// Ensure Bayer2x2Dither implements the image.Image interface.
var _ image.Image = (*Bayer2x2Dither)(nil)

// Bayer2x2Dither applies a 2x2 Bayer ordered dither to an image using a provided palette.
type Bayer2x2Dither struct {
	Null
	img     image.Image
	palette []color.Color
}

func (p *Bayer2x2Dither) ColorModel() color.Model {
	if p.img != nil {
		return p.img.ColorModel()
	}
	return color.RGBAModel
}

func (p *Bayer2x2Dither) Bounds() image.Rectangle {
	return p.bounds
}

func (p *Bayer2x2Dither) At(x, y int) color.Color {
	if p.img == nil {
		return color.RGBA{}
	}

	c := p.img.At(x, y)
	r, g, b, _ := c.RGBA()

	// Calculate luminance: 0.299*R + 0.587*G + 0.114*B
	// Using coefficients from color.GrayModel
	lum := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	// Normalize to 0.0 - 1.0
	lumFloat := float64(lum) / 65535.0

	n := len(p.palette)
	// Safety check, though constructor ensures palette is not empty
	if n == 0 {
		return c
	}
	if n == 1 {
		return p.palette[0]
	}

	// Map luminance to palette range [0, n-1]
	val := lumFloat * float64(n-1)
	base := int(val)
	rem := val - float64(base)

	// Bayer 2x2 matrix
	// 0 2
	// 3 1
	// Divided by 4.
	matrix := [2][2]float64{
		{0.0 / 4.0, 2.0 / 4.0},
		{3.0 / 4.0, 1.0 / 4.0},
	}

	// Handle negative coordinates correctly for modulo
	mx := x % 2
	if mx < 0 {
		mx += 2
	}
	my := y % 2
	if my < 0 {
		my += 2
	}

	threshold := matrix[my][mx]

	idx := base
	if rem > threshold {
		idx++
	}

	if idx >= n {
		idx = n - 1
	}

	return p.palette[idx]
}

// NewBayer2x2Dither creates a new Bayer2x2Dither pattern.
// src is the source image.
// palette is the list of colors to map to based on luminance.
func NewBayer2x2Dither(src image.Image, palette []color.Color, ops ...func(any)) image.Image {
	bounds := image.Rect(0, 0, 100, 100)
	if src != nil {
		bounds = src.Bounds()
	}

	// If palette is empty, default to Black and White
	if len(palette) == 0 {
		palette = []color.Color{color.Black, color.White}
	}

	p := &Bayer2x2Dither{
		img:     src,
		palette: palette,
		Null: Null{
			bounds: bounds,
		},
	}

	for _, op := range ops {
		op(p)
	}

	return p
}
