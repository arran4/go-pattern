package pattern

import (
	"image"
	"image/color"
)

// Ensure SimpleZoom implements the image.Image interface.
var _ image.Image = (*SimpleZoom)(nil)

// SimpleZoom is a pattern that zooms in on an underlying image.
type SimpleZoom struct {
	Null
	img    image.Image
	factor int
}

func (s *SimpleZoom) ColorModel() color.Model {
	return s.img.ColorModel()
}

func (s *SimpleZoom) Bounds() image.Rectangle {
	return s.bounds
}

func (s *SimpleZoom) At(x, y int) color.Color {
	return s.img.At(x/s.factor, y/s.factor)
}

const maxZoomDimension = 16384

// NewSimpleZoom creates a new SimpleZoom with the given image and zoom factor.
func NewSimpleZoom(img image.Image, factor int, ops ...func(any)) image.Image {
	if factor <= 0 {
		factor = 1
	}
	b := img.Bounds()

	// Safety check to prevent integer overflow and excessive memory allocation
	w := int64(b.Dx())
	h := int64(b.Dy())

	if w > 0 && h > 0 {
		if int64(factor) > maxZoomDimension/w {
			factor = int(maxZoomDimension / w)
		}
		if int64(factor) > maxZoomDimension/h {
			factor = int(maxZoomDimension / h)
		}
		if factor < 1 {
			factor = 1
		}
	}

	p := &SimpleZoom{
		img:    img,
		factor: factor,
		Null: Null{
			bounds: image.Rect(b.Min.X, b.Min.Y, b.Max.X*factor, b.Max.Y*factor),
		},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoSimpleZoom produces a demo variant for readme.md pre-populated values
func NewDemoSimpleZoom(img image.Image, ops ...func(any)) image.Image {
	return NewSimpleZoom(img, 2, ops...)
}
