package pattern

import (
	"image"
)

func GenerateBlueNoise(ops ...func(any)) image.Image {
	return NewBlueNoise(ops...)
}

func GenerateBlueNoiseReferences() (map[string]image.Image, []string) {
	return nil, nil
}
