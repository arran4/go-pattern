package pattern

import (
	"image"
	"image/color"
)

// FillColor configures the fill color in a pattern.
type FillColor struct {
	FillColor color.Color
}

func (s *FillColor) SetFillColor(v color.Color) {
	s.FillColor = v
}

type hasFillColor interface {
	SetFillColor(color.Color)
}

// SetFillColor creates an option to set the fill color.
func SetFillColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasFillColor); ok {
			h.SetFillColor(v)
		}
	}
}

// FillImageSource configures an image source for fill in a pattern.
type FillImageSource struct {
	FillImageSource image.Image
}

func (s *FillImageSource) SetFillImageSource(v image.Image) {
	s.FillImageSource = v
}

type hasFillImageSource interface {
	SetFillImageSource(image.Image)
}

// SetFillImageSource creates an option to set the fill image source.
func SetFillImageSource(v image.Image) func(any) {
	return func(i any) {
		if h, ok := i.(hasFillImageSource); ok {
			h.SetFillImageSource(v)
		}
	}
}
