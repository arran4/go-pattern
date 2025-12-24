package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestCircle_At(t *testing.T) {
	// 10x10 image. Center should be at (5,5).
	// Pixel range 0-9. Center logic: 0+10 = 10 -> cx=5.
	// 2*coord logic:
	// cx2 = 0 + 10 = 10.
	// cy2 = 0 + 10 = 10.
	// radius = 10 / 2 = 5.
	// diameter = 10. radiusSq = 100 (in 2x space, actually diameter^2).

	// Test center pixel (4,4) -> (9,9)
	// dx2 = 9 - 10 = -1
	// dy2 = 9 - 10 = -1
	// distSq = 1 + 1 = 2 <= 100. Inside.

	// Test pixel (0,5) -> (1, 11)
	// dx2 = 1 - 10 = -9
	// dy2 = 11 - 10 = 1
	// distSq = 81 + 1 = 82 <= 100. Inside.

	// Test pixel (0,0) -> (1, 1)
	// dx2 = 1 - 10 = -9
	// dy2 = 1 - 10 = -9
	// distSq = 81 + 81 = 162 > 100. Outside.

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
		{0, 5, color.Black},
		{5, 0, color.Black},
		{0, 0, color.White},
		{9, 9, color.White}, // 2*9+1 = 19. diff 9. sq 81. 81+81 = 162. Outside.
		{2, 2, color.Black}, // 2*2+1 = 5. diff -5. sq 25. 25+25 = 50 <= 100. Inside.
	}

	for _, tt := range tests {
		got := c.At(tt.x, tt.y)
		if got != tt.want {
			t.Errorf("At(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.want)
		}
	}
}

func TestCircle_OddSize(t *testing.T) {
	// 11x11 image.
	// cx2 = 0 + 11 = 11. (Center at 5.5)
	// diameter = 11. radiusSq = 121.

	// Pixel 5 (center). 2*5+1 = 11. diff 0. Inside.

	c := NewCircle(
		SetBounds(image.Rect(0, 0, 11, 11)),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
	)

	if c.At(5, 5) != color.Black {
		t.Errorf("Center pixel should be black")
	}

	// Pixel 0, 5. 2*0+1 = 1. diff -10. sq 100. dy=0. 100 <= 121. Inside.
	if c.At(0, 5) != color.Black {
		t.Errorf("Edge pixel (0,5) should be black")
	}

	// Pixel 0,0. 2*0+1=1. diff -10. sq 100. 100+100=200 > 121. Outside.
	if c.At(0, 0) != color.White {
		t.Errorf("Corner pixel (0,0) should be white")
	}
}
