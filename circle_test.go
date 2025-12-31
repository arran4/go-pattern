package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestCircle_At(t *testing.T) {
	// 10x10 image. Center (5,5).
	// Legacy mode (LineSize 0).
	c := NewCircle(
		SetBounds(image.Rect(0, 0, 10, 10)),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
	)

	tests := []struct {
		x, y int
		want color.Color
	}{
		{4, 4, color.Black},
		{5, 5, color.Black},
		{0, 0, color.White},
	}

	for _, tt := range tests {
		got := c.At(tt.x, tt.y)
		if got != tt.want {
			t.Errorf("Legacy At(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.want)
		}
	}
}

func TestCircle_BorderAndFill(t *testing.T) {
	// 20x20 image. Center (10,10). Radius 10.
	// LineSize 2.
	// Border: Radius 10 to 8.
	// Fill: Radius 8 to 0.

	c := NewCircle(
		SetBounds(image.Rect(0, 0, 20, 20)),
		SetLineSize(2),
		SetLineColor(color.Black), // Border
		SetFillColor(color.RGBA{255, 0, 0, 255}), // Red Fill
		SetSpaceColor(color.White), // Background
	)

	// Center pixel (10,10) should be Fill (Red).
	if c.At(10, 10) != (color.RGBA{255, 0, 0, 255}) {
		t.Errorf("Center pixel should be Red Fill")
	}

	// Radius 9. (e.g. x=19, y=10. dx=9).
	// 2*19+1 = 39. Center 20. Diff 19. 19/2 = 9.5?
	// Let's use coordinate logic.
	// bounds (0,0)-(20,20). cx2=20, cy2=20.
	// Diameter 20. OuterSq = 400.
	// LineSize 2. InnerDiameter = 20 - 4 = 16. InnerSq = 256.

	// Test Point (19, 10).
	// x=19 -> 2x+1 = 39. dx2 = 19.
	// y=10 -> 2y+1 = 21. dy2 = 1.
	// DistSq = 361 + 1 = 362.
	// 362 <= 400 (Inside Outer).
	// 362 > 256 (Outside Inner).
	// Should be Border (Black).
	if c.At(19, 10) != color.Black {
		t.Errorf("Pixel (19,10) should be Border (Black), got %v", c.At(19, 10))
	}

	// Test Point (14, 10).
	// x=14 -> 2x+1 = 29. dx2 = 9.
	// y=10 -> 2y+1 = 21. dy2 = 1.
	// DistSq = 81 + 1 = 82.
	// 82 <= 256. Inside Inner.
	// Should be Fill (Red).
	if c.At(14, 10) != (color.RGBA{255, 0, 0, 255}) {
		t.Errorf("Pixel (14,10) should be Fill (Red), got %v", c.At(14, 10))
	}
}

func TestCircle_FillImage(t *testing.T) {
	// Use Checker as fill source.
	// We explicitly set SpaceSize(1) to restore legacy 1x1 behavior for this test.
	check := NewChecker(color.White, color.Black, SetSpaceSize(1))
	// Circle with check fill, no border.
	c := NewCircle(
		SetBounds(image.Rect(0, 0, 10, 10)),
		SetFillImageSource(check),
	)

	// Center (5,5). Checker at (5,5). 5%2 != 5%2? 5%2=1. 1==1. -> White.
	if c.At(5, 5) != color.White {
		t.Errorf("Center pixel should be White from Checker")
	}
	// (5,4). 4%2=0. 1!=0. -> Black.
	if c.At(5, 4) != color.Black {
		t.Errorf("Pixel (5,4) should be Black from Checker")
	}
}

func TestCircle_OddSize(t *testing.T) {
	c := NewCircle(
		SetBounds(image.Rect(0, 0, 11, 11)),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
	)

	if c.At(5, 5) != color.Black {
		t.Errorf("Center pixel should be black")
	}
}
