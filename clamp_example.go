package pattern

import (
	"image"
	"image/color"
)

var (
	ClampOutputFilename = "clamp.png"
	ClampZoomLevels     = []int{}
	ClampOrder          = 20
	ClampBaseLabel      = "Clamp"
)

func init() {
	RegisterGenerator("Clamp", func(bounds image.Rectangle) image.Image {
		return ExampleNewClamp(SetBounds(bounds))
	})
	RegisterReferences("Clamp", GenerateClampReferences)
}

func ExampleNewClamp(ops ...func(any)) image.Image {
	// Create a small source image
	src := NewCircle(SetFillColor(color.RGBA{255, 0, 0, 255}))
	// Crop it effectively makes it smaller if we assume it was larger,
	// or we can just use it.
	// Let's crop it to have a specific small bound.
	src = NewCrop(src, image.Rect(100, 100, 150, 150))

	// Clamp to 255x255 (default) or whatever the caller wants?
	// The example function is expected to return the image.
	// If called by bootstrap, it might not set bounds via ops.

	// Let's return a Clamped version of the cropped circle.
	// We want to show the clamping effect, so the output image needs to be larger than the source crop.
	// The bootstrap tool usually asks for 150x150.
	// Our crop is 50x50.
	// So we return NewClamp(src, bounds).
	// But NewClamp takes bounds in constructor.

	// We need to know the target bounds.
	// We can use a default if not provided, but usually ExampleNew* is used for the main demo image.

	// Pass ops to NewClamp so SetBounds can override the default
	return NewClamp(src, image.Rect(50, 50, 200, 200), ops...)
}

func GenerateClampReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Reference": func(bounds image.Rectangle) image.Image {
			src := NewGopher()
			cropped := NewCrop(src, image.Rect(50, 20, 100, 70))
			return NewClamp(cropped, bounds)
		},
	}, []string{"Reference"}
}
