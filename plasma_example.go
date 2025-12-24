package pattern

import (
	"image"
)

func GeneratePlasma(ops ...func(any)) image.Image {
	return NewPlasma(ops...)
}

func GeneratePlasmaReferences() (map[string]image.Image, []string) {
	return nil, nil
}
