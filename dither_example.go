package pattern

import (
	"image"
	"image/color"
)

func GenerateBayerDither(ops ...func(any)) image.Image {
	// Dither a gradient
	grad := NewLinearGradient(
		SetStartColor(color.Black),
		SetEndColor(color.White),
	)
	return NewBayerDither(grad, 4, ops...)
}

func GenerateBayerDitherReferences() (map[string]image.Image, []string) {
	return nil, nil
}
