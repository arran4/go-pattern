package pattern

import (
	"image"
	"image/color"
)

func GenerateModuloStripe(ops ...func(any)) image.Image {
	return NewModuloStripe([]color.Color{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
	}, ops...)
}

func GenerateModuloStripeReferences() (map[string]image.Image, []string) {
	return nil, nil
}
