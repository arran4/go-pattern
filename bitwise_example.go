package pattern

import (
	"image"
	"image/color"
)

func GenerateBitwiseAnd(ops ...func(any)) image.Image {
	h := NewHorizontalLine(SetLineSize(50), SetSpaceSize(50), SetLineColor(color.RGBA{255, 0, 0, 255}))
	v := NewVerticalLine(SetLineSize(50), SetSpaceSize(50), SetLineColor(color.RGBA{0, 255, 0, 255}))
	return NewBitwiseAnd([]image.Image{h, v}, ops...)
}

func GenerateBitwiseAndReferences() (map[string]image.Image, []string) {
	return nil, nil
}
