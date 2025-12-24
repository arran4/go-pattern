package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestMirror_Horizontal(t *testing.T) {
	// 2x1 Image: [Red, Blue]
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	img.Set(0, 0, red)
	img.Set(1, 0, blue)

	m := NewMirror(img, true, false)

	// Expected: [Blue, Red]
	if !colorsEqual(m.At(0, 0), blue) {
		t.Errorf("At(0,0) expected Blue, got %v", m.At(0, 0))
	}
	if !colorsEqual(m.At(1, 0), red) {
		t.Errorf("At(1,0) expected Red, got %v", m.At(1, 0))
	}
}

func TestMirror_Vertical(t *testing.T) {
	// 1x2 Image:
	// [Red]
	// [Blue]
	img := image.NewRGBA(image.Rect(0, 0, 1, 2))
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	img.Set(0, 0, red)
	img.Set(0, 1, blue)

	m := NewMirror(img, false, true)

	// Expected:
	// [Blue]
	// [Red]
	if !colorsEqual(m.At(0, 0), blue) {
		t.Errorf("At(0,0) expected Blue, got %v", m.At(0, 0))
	}
	if !colorsEqual(m.At(0, 1), red) {
		t.Errorf("At(0,1) expected Red, got %v", m.At(0, 1))
	}
}

func TestMirror_Both(t *testing.T) {
	// 2x2 Image:
	// [Red, Green]
	// [Blue, White]
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	red := color.RGBA{255, 0, 0, 255}
	green := color.RGBA{0, 255, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	white := color.RGBA{255, 255, 255, 255}
	img.Set(0, 0, red)
	img.Set(1, 0, green)
	img.Set(0, 1, blue)
	img.Set(1, 1, white)

	m := NewMirror(img, true, true)

	// Expected:
	// [White, Blue]
	// [Green, Red]

	// (0,0) -> was (1,1) White
	if !colorsEqual(m.At(0, 0), white) {
		t.Errorf("At(0,0) expected White, got %v", m.At(0, 0))
	}
	// (1,0) -> was (0,1) Blue
	if !colorsEqual(m.At(1, 0), blue) {
		t.Errorf("At(1,0) expected Blue, got %v", m.At(1, 0))
	}
	// (0,1) -> was (1,0) Green
	if !colorsEqual(m.At(0, 1), green) {
		t.Errorf("At(0,1) expected Green, got %v", m.At(0, 1))
	}
	// (1,1) -> was (0,0) Red
	if !colorsEqual(m.At(1, 1), red) {
		t.Errorf("At(1,1) expected Red, got %v", m.At(1, 1))
	}
}

func colorsEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
