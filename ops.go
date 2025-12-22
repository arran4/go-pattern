package pattern

import "image"

type hasBounds interface {
	SetBounds(image.Rectangle)
}

func SetBounds(bounds image.Rectangle) func(ai any) {
	return func(i any) {
		if i, ok := i.(hasBounds); ok {
			i.SetBounds(bounds)
		}
	}
}
