package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestPadding(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.White)

	padding := 5
	p := NewPadding(img, padding, nil)

	expectedW := 10 + 2*padding
	// expectedH := 10 + 2*padding

	if p.Bounds().Dx() != expectedW {
		t.Errorf("Expected width %d, got %d", expectedW, p.Bounds().Dx())
	}

	// Check content position.
	// Image starts at (5,5) in local coordinates?
	// Bounds start at (0,0).
	// Padding puts image at (padding, padding).

	// Check (padding, padding) -> should be original (0,0) -> White
	if r, _, _, _ := p.At(padding, padding).RGBA(); r == 0 {
		t.Error("Expected pixel at padding,padding to be white")
	}

	// Check (0,0) -> should be transparent (nil background)
	if _, _, _, a := p.At(0, 0).RGBA(); a != 0 {
		t.Error("Expected pixel at 0,0 to be transparent")
	}
}

func TestPaddingWithBackground(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	bg := image.NewUniform(color.Black)

	p := NewPadding(img, 5, bg)

	// Check (0,0) -> Black
	r, g, b, _ := p.At(0, 0).RGBA()
	if r != 0 || g != 0 || b != 0 {
		// Note: Black is 0,0,0
		// Uniform Black might return alpha?
		// color.Black is opaque black.
		t.Errorf("Expected black background, got r=%d g=%d b=%d", r, g, b)
	}
	// Check alpha
	_, _, _, a := p.At(0, 0).RGBA()
	if a == 0 {
		t.Error("Expected background to be opaque")
	}
}
