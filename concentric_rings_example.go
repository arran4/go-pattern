package pattern

import (
	"image"
	"image/color"
)

func GenerateConcentricRings(ops ...func(any)) image.Image {
	return NewConcentricRings([]color.Color{
		color.Black,
		color.White,
		color.RGBA{255, 0, 0, 255},
	}, ops...)
}

func GenerateConcentricRingsReferences() (map[string]image.Image, []string) {
	return nil, nil
}
