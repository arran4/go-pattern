package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestVoronoi_At(t *testing.T) {
	points := []image.Point{
		{10, 10},
		{20, 10},
	}
	colors := []color.Color{
		color.White,
		color.Black,
	}

	v := NewVoronoi(points, colors)

	// Point exactly at first seed
	c := v.At(10, 10)
	if c != color.White {
		t.Errorf("Expected White at (10, 10), got %v", c)
	}

	// Point exactly at second seed
	c = v.At(20, 10)
	if c != color.Black {
		t.Errorf("Expected Black at (20, 10), got %v", c)
	}

	// Point closer to first seed
	c = v.At(14, 10) // distance to 10 is 4, to 20 is 6
	if c != color.White {
		t.Errorf("Expected White at (14, 10), got %v", c)
	}

	// Point closer to second seed
	c = v.At(16, 10) // distance to 10 is 6, to 20 is 4
	if c != color.Black {
		t.Errorf("Expected Black at (16, 10), got %v", c)
	}

	// Equidistant point (implementation dependent, usually first one found if strictly less, or last one if <=)
	// My implementation uses strictly less (<), so it keeps the first one if distances are equal.
	// 10 -> 10, 10 -> 5 distance sq = 25
	// 20 -> 10, 10 -> 5 distance sq = 25
	c = v.At(15, 10)
	if c != color.White {
		t.Errorf("Expected White at (15, 10) due to order preference, got %v", c)
	}
}
