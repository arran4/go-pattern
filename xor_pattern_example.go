package pattern

import (
	"image"
)

func GenerateXorPattern(ops ...func(any)) image.Image {
	return NewXorPattern(ops...)
}

func GenerateXorPatternReferences() (map[string]image.Image, []string) {
	return nil, nil
}
