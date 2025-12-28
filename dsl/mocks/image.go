package mocks

import (
	"image"
	"image/color"
)

type MockImage struct {
	BoundsVal image.Rectangle
	ColorVal  color.Color
}

func (m *MockImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (m *MockImage) Bounds() image.Rectangle {
	return m.BoundsVal
}

func (m *MockImage) At(x, y int) color.Color {
	return m.ColorVal
}

func NewMockImage(w, h int, c color.Color) *MockImage {
	return &MockImage{
		BoundsVal: image.Rect(0, 0, w, h),
		ColorVal:  c,
	}
}
