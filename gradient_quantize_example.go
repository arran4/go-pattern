package pattern

import (
	"image"
	"image/color"
)

func GenerateGradientQuantization(ops ...func(any)) image.Image {
	grad := NewLinearGradient(
		SetStartColor(color.Black),
		SetEndColor(color.White),
	)
	// Quantize to 4 levels
	return NewQuantize(grad, 4, ops...)
}

func GenerateGradientQuantizationReferences() (map[string]image.Image, []string) {
	return nil, nil
}
