package pattern

import (
	"image"
	"image/color"
)

// Ensure MaskedBlend implements the image.Image interface.
var _ image.Image = (*MaskedBlend)(nil)

// MaskedBlend blends a foreground and background image based on a mask image.
type MaskedBlend struct {
	Null
	Background image.Image
	Foreground image.Image
	Mask       image.Image
	Predicate  ColorPredicate
}

// At returns the color at (x, y).
func (mb *MaskedBlend) At(x, y int) color.Color {
	// If mask is nil, return background (effectively mask=0)
	if mb.Mask == nil {
		if mb.Background != nil {
			return mb.Background.At(x, y)
		}
		return color.Transparent
	}

	cMask := mb.Mask.At(x, y)
	predicate := mb.Predicate
	if predicate == nil {
		// Fallback if somehow predicate is nil, though constructor sets it.
		predicate = PredicateFuzzyGray()
	}

	t := predicate(cMask)

	// Optimization for strict 0 or 1
	if t <= 0 {
		if mb.Background != nil {
			return mb.Background.At(x, y)
		}
		return color.Transparent
	}
	if t >= 1 {
		if mb.Foreground != nil {
			return mb.Foreground.At(x, y)
		}
		return color.Transparent
	}

	var cBg, cFg color.Color
	if mb.Background != nil {
		cBg = mb.Background.At(x, y)
	} else {
		cBg = color.Transparent
	}
	if mb.Foreground != nil {
		cFg = mb.Foreground.At(x, y)
	} else {
		cFg = color.Transparent
	}

	r0, g0, b0, a0 := cBg.RGBA()
	r1, g1, b1, a1 := cFg.RGBA()

	// Interpolate
	r := float64(r0) + t*(float64(r1)-float64(r0))
	g := float64(g0) + t*(float64(g1)-float64(g0))
	b := float64(b0) + t*(float64(b1)-float64(b0))
	a := float64(a0) + t*(float64(a1)-float64(a0))

	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

// NewMaskedBlend creates a new MaskedBlend pattern.
// bg: Background image
// fg: Foreground image
// mask: Mask image (determines how much of fg is shown)
func NewMaskedBlend(bg, fg, mask image.Image, ops ...func(any)) image.Image {
	mb := &MaskedBlend{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Background: bg,
		Foreground: fg,
		Mask:       mask,
		Predicate:  PredicateFuzzyGray(),
	}

	for _, op := range ops {
		op(mb)
	}
	return mb
}

// SetPredicate sets the color predicate used to determine the mask value.
func (mb *MaskedBlend) SetPredicate(p ColorPredicate) {
	mb.Predicate = p
}
