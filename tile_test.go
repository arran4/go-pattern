package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestTile(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.White) // Top-left pixel is white

	// Tile to 20x20 (2x2 tiles)
	tileRect := image.Rect(0, 0, 20, 20)
	tl := NewTile(img, tileRect)

	if tl.Bounds() != tileRect {
		t.Errorf("Expected bounds %v, got %v", tileRect, tl.Bounds())
	}

	// Check (0,0) -> White
	if r, _, _, _ := tl.At(0, 0).RGBA(); r == 0 {
		t.Error("Expected pixel at 0,0 to be white")
	}

	// Check (10,0) -> Should be start of next tile -> White
	if r, _, _, _ := tl.At(10, 0).RGBA(); r == 0 {
		t.Error("Expected pixel at 10,0 to be white")
	}

	// Check (10,10) -> White
	if r, _, _, _ := tl.At(10, 10).RGBA(); r == 0 {
		t.Error("Expected pixel at 10,10 to be white")
	}
}
