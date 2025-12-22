package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestScale(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// Set (0,0) to white
	img.Set(0, 0, color.White)

	// Scale by 2x
	s := NewScale(img, ScaleFactor(2.0))

	expected := image.Rect(0, 0, 20, 20)
	if s.Bounds() != expected {
		t.Errorf("Expected bounds %v, got %v", expected, s.Bounds())
	}

	// Scale to fixed size
	s2 := NewScale(img, ScaleSize(5, 5))
	expected2 := image.Rect(0, 0, 5, 5)
	if s2.Bounds() != expected2 {
		t.Errorf("Expected bounds %v, got %v", expected2, s2.Bounds())
	}
}
