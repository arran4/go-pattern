package pattern

import (
	"image"
	"image/color"
)

// Ensure Rect implements the image.Image interface.
var _ image.Image = (*Rect)(nil)

// Rect is a pattern that draws a filled rectangle.
type Rect struct {
	Null
	FillColor
	LineSize
	LineColor
	LineImageSource
}

func (r *Rect) At(x, y int) color.Color {
	// Check if we are inside the bounds.
	// Since Rect embeds Null, and we default bounds to something,
	// we should respect them.
	// Note: image.Image.At contract says "The bounds of the image do not necessarily contain the point (x, y)."
	// It doesn't strictly say it must return zero color outside, but that's typical for "bounded" images
	// unless they are infinite patterns.
	// The `Null` struct has `bounds`.
	if !(image.Point{x, y}.In(r.bounds)) {
		return color.RGBA{}
	}

	ls := r.LineSize.LineSize
	if ls > 0 {
		min := r.bounds.Min
		max := r.bounds.Max
		// Top or Bottom border
		// Top: y in [min.Y, min.Y + ls)
		// Bottom: y in [max.Y - ls, max.Y)
		if y < min.Y+ls || y >= max.Y-ls {
			return r.getLineColor(x, y)
		}
		// Left or Right border
		// Left: x in [min.X, min.X + ls)
		// Right: x in [max.X - ls, max.X)
		if x < min.X+ls || x >= max.X-ls {
			return r.getLineColor(x, y)
		}
	}

	return r.FillColor.FillColor
}

func (r *Rect) getLineColor(x, y int) color.Color {
	if r.LineImageSource.LineImageSource != nil {
		return r.LineImageSource.LineImageSource.At(x, y)
	}
	if r.LineColor.LineColor != nil {
		return r.LineColor.LineColor
	}
	// Fallback if LineColor is nil but LineSize > 0?
	// Standard might be Black, or Transparent.
	// Let's assume Black as per other patterns.
	return color.Black
}

// NewRect creates a new Rect pattern with the given options.
func NewRect(ops ...func(any)) image.Image {
	p := &Rect{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.FillColor.FillColor = color.Black
	p.LineColor.LineColor = color.Black

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoRect produces a demo variant for readme.md pre-populated values
func NewDemoRect(ops ...func(any)) image.Image {
	return NewRect(ops...)
}
