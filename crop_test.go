package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestCrop(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with white
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.White)
		}
	}

	cropRect := image.Rect(10, 10, 20, 20)
	c := NewCrop(img, cropRect)

	if c.Bounds() != cropRect {
		t.Errorf("Expected bounds %v, got %v", cropRect, c.Bounds())
	}

	// Test At
	// Inside crop
	if _, _, _, a := c.At(15, 15).RGBA(); a == 0 {
		t.Error("Expected visible pixel at 15,15")
	}

	// Outside crop
	if _, _, _, a := c.At(5, 5).RGBA(); a != 0 {
		t.Error("Expected transparent pixel at 5,5")
	}
}
