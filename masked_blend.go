package pattern

import (
	"image"
	"image/color"
)

// Ensure MaskedBlend implements the image.Image interface.
var _ image.Image = (*MaskedBlend)(nil)

// MaskedBlend combines two images using a mask image and a predicate.
type MaskedBlend struct {
	Null
	Background image.Image
	Foreground image.Image
	Mask       image.Image
	Predicate  ColorPredicate
}

func (m *MaskedBlend) At(x, y int) color.Color {
	if m.Background == nil {
		return color.RGBA{}
	}
	if m.Foreground == nil {
		return m.Background.At(x, y)
	}
	if m.Mask == nil {
		return m.Background.At(x, y)
	}

	bg := m.Background.At(x, y)
	fg := m.Foreground.At(x, y)
	maskColor := m.Mask.At(x, y)

	t := m.Predicate(maskColor)
	return InterpolateColor(bg, fg, t)
}

func (m *MaskedBlend) SetPredicate(p ColorPredicate) {
	m.Predicate = p
}

// NewMaskedBlend creates a new MaskedBlend pattern.
// By default, it uses PredicateFuzzyGray which works well for grayscale masks.
func NewMaskedBlend(bg, fg, mask image.Image, ops ...func(any)) image.Image {
	p := &MaskedBlend{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Background: bg,
		Foreground: fg,
		Mask:       mask,
		Predicate:  PredicateFuzzyGray(),
	}
	for _, op := range ops {
		op(p)
	}
	return p
}
