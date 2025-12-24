package pattern

import (
	"image"
)

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
