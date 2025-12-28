package pattern

import (
	"image"
)

var (
	ScaleOutputFilename = "scale.png"
	ScaleZoomLevels     = []int{}
	ScaleOrder          = 22
	ScaleBaseLabel      = "Scale"
)

func init() {
	RegisterGenerator("Scale", func(bounds image.Rectangle) image.Image {
		return ExampleNewScale(SetBounds(bounds))
	})
	RegisterReferences("Scale", GenerateScaleReferences)
}

func ExampleNewScale(ops ...func(any)) image.Image {
	src := NewGopher()
	// Scale down by 0.5
	return NewScale(src, ScaleX(0.5), ScaleY(0.5))
}

func GenerateScaleReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"ZoomIn": func(bounds image.Rectangle) image.Image {
			src := NewGopher()
			// Zoom in (2.0)
			s := NewScale(src, ScaleX(2.0), ScaleY(2.0))
			// Usually we want to crop it to the view
			return NewCrop(s, bounds)
		},
		"Stretch": func(bounds image.Rectangle) image.Image {
			src := NewGopher()
			// Stretch X only
			s := NewScale(src, ScaleX(2.0), ScaleY(1.0))
			return NewCrop(s, bounds)
		},
	}, []string{"ZoomIn", "Stretch"}
}
