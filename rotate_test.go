package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestRotate_90(t *testing.T) {
	// 2x3 Image
	// (0,0) Red
	// Others Black
	w, h := 2, 3
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	red := color.RGBA{255, 0, 0, 255}
	img.Set(0, 0, red)

	// Rotate 90 CW
	r := NewRotate(img, 90)

	// Bounds should be 3x2
	if r.Bounds().Dx() != 3 || r.Bounds().Dy() != 2 {
		t.Errorf("Bounds expected 3x2, got %dx%d", r.Bounds().Dx(), r.Bounds().Dy())
	}

	// (0,0) Red -> Becomes (h-1, 0) = (2, 0) ?
	// Let's trace my logic again.
	// NewRotate(90): sx = dy, sy = h - 1 - dx
	// At(destX, destY)
	// At(2, 0): dx=2, dy=0.
	// sx = 0. sy = 3 - 1 - 2 = 0. -> Src(0,0).
	// So (2,0) should be Red.

	if !colorsEqual(r.At(2, 0), red) {
		t.Errorf("At(2,0) expected Red, got %v", r.At(2, 0))
	}

	// Check origin (0,0)
	// At(0,0): dx=0, dy=0. sx=0, sy=2. Src(0,2). Should be Black.
	if colorsEqual(r.At(0, 0), red) {
		t.Errorf("At(0,0) expected Black, got Red")
	}
}

func TestRotate_180(t *testing.T) {
	// 2x2 Image
	// (0,0) Red
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	red := color.RGBA{255, 0, 0, 255}
	img.Set(0, 0, red)

	r := NewRotate(img, 180)

	// Bounds 2x2
	if r.Bounds() != img.Bounds() {
		t.Errorf("Bounds mismatch")
	}

	// (0,0) -> (1,1)
	// At(1,1): dx=1, dy=1. sx=2-1-1=0. sy=2-1-1=0. Src(0,0).
	if !colorsEqual(r.At(1, 1), red) {
		t.Errorf("At(1,1) expected Red, got %v", r.At(1, 1))
	}
}

func TestRotate_270(t *testing.T) {
	// 2x3 Image
	// (0,0) Red
	img := image.NewRGBA(image.Rect(0, 0, 2, 3))
	red := color.RGBA{255, 0, 0, 255}
	img.Set(0, 0, red)

	r := NewRotate(img, 270)

	// Bounds 3x2
	if r.Bounds().Dx() != 3 || r.Bounds().Dy() != 2 {
		t.Errorf("Bounds expected 3x2, got %dx%d", r.Bounds().Dx(), r.Bounds().Dy())
	}

	// 270 CW (or 90 CCW)
	// (0,0) Top-Left -> Bottom-Left (0, 2)? No, dest height is 2.
	// (0,0) -> (0, 1) in 3x2 image?
	// Let's trace.
	// Case 270: sx = w - 1 - dy, sy = dx
	// At(0, 1): dx=0, dy=1.
	// sx = 2 - 1 - 1 = 0. sy = 0. Src(0,0).
	// So (0,1) should be Red.

	if !colorsEqual(r.At(0, 1), red) {
		t.Errorf("At(0,1) expected Red, got %v", r.At(0, 1))
	}
}
