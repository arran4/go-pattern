package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestNewSimpleZoom_DoS(t *testing.T) {
	// Create a 10x10 base image
	myImg := &MyImage{image.Rect(0, 0, 10, 10)}

	// Use a huge factor
	hugeFactor := 1_000_000
	zoomed := NewSimpleZoom(myImg, hugeFactor)

	b := zoomed.Bounds()

	// Expect clamped bounds
	// Max dimension is 16384.

	if b.Dx() > maxZoomDimension {
		t.Errorf("Expected width <= %d, got %d", maxZoomDimension, b.Dx())
	}

	if b.Dy() > maxZoomDimension {
		t.Errorf("Expected height <= %d, got %d", maxZoomDimension, b.Dy())
	}
}

func TestNewSimpleZoom_Normal(t *testing.T) {
	myImg := &MyImage{image.Rect(0, 0, 10, 10)}
	factor := 2
	zoomed := NewSimpleZoom(myImg, factor)
	b := zoomed.Bounds()
	if b.Dx() != 20 {
		t.Errorf("Expected width 20, got %d", b.Dx())
	}
}

type MyImage struct {
	bounds image.Rectangle
}

func (m *MyImage) ColorModel() color.Model { return color.RGBAModel }
func (m *MyImage) Bounds() image.Rectangle { return m.bounds }
func (m *MyImage) At(_, _ int) color.Color { return color.Black }
