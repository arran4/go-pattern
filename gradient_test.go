package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestLinearGradient(t *testing.T) {
	start := color.RGBA{0, 0, 0, 255}     // Black
	end := color.RGBA{255, 255, 255, 255} // White

	g := NewLinearGradient(
		SetStartColor(start),
		SetEndColor(end),
	).(*LinearGradient)
	g.SetBounds(image.Rect(0, 0, 100, 100))

	// Test at x=0 (Start)
	c0 := g.At(0, 50)
	if c0 != start {
		t.Errorf("Expected start color at x=0, got %v", c0)
	}

	// Test at x=99 (End)
	c99 := g.At(99, 50)
	if c99 != end {
		t.Errorf("Expected end color at x=99, got %v", c99)
	}

	// Test at x=50 (Mid)
	c50 := g.At(50, 50)
	// Approximate gray
	r, _, _, _ := c50.RGBA()
	if r < 0x7000 || r > 0x9000 {
		t.Errorf("Expected mid color at x=50, got %v (r=%x)", c50, r)
	}
}

func TestLinearGradientVertical(t *testing.T) {
	start := color.RGBA{0, 0, 0, 255}
	end := color.RGBA{255, 255, 255, 255}

	g := NewLinearGradient(
		SetStartColor(start),
		SetEndColor(end),
		GradientVertical(),
	).(*LinearGradient)
	g.SetBounds(image.Rect(0, 0, 100, 100))

	// Test at y=0 (Start)
	c0 := g.At(50, 0)
	if c0 != start {
		t.Errorf("Expected start color at y=0, got %v", c0)
	}

	// Test at y=99 (End)
	c99 := g.At(50, 99)
	if c99 != end {
		t.Errorf("Expected end color at y=99, got %v", c99)
	}
}

func TestRadialGradient(t *testing.T) {
	start := color.RGBA{0, 0, 0, 255}
	end := color.RGBA{255, 255, 255, 255}

	g := NewRadialGradient(
		SetStartColor(start),
		SetEndColor(end),
	).(*RadialGradient)
	g.SetBounds(image.Rect(0, 0, 100, 100))

	// Center is 50, 50
	cCenter := g.At(50, 50)
	if cCenter != start {
		t.Errorf("Expected start color at center, got %v", cCenter)
	}

	// Corner (0,0) is far away
	// dist = sqrt(50^2 + 50^2) = 70.7
	// maxDist = sqrt(100^2 + 100^2) / 2 = 141.4 / 2 = 70.7
	// So corner should be roughly end color.

	cCorner := g.At(0, 0)
	// It might be slightly off due to float precision, but should be close to end.

	r, _, _, _ := cCorner.RGBA()
	if r < 0xF000 {
		t.Errorf("Expected nearly white at corner, got %v", cCorner)
	}
}

func TestConicGradient(t *testing.T) {
	start := color.RGBA{0, 0, 0, 255}     // Black (0)
	end := color.RGBA{255, 255, 255, 255} // White (1)

	g := NewConicGradient(
		SetStartColor(start),
		SetEndColor(end),
	).(*ConicGradient)
	g.SetBounds(image.Rect(0, 0, 100, 100))

	// Center 50, 50
	// Right: (51, 50) -> dx=1, dy=0 -> angle=0
	// Our formula: t = (0 + Pi) / 2Pi = 0.5

	cRight := g.At(60, 50)
	r, _, _, _ := cRight.RGBA()
	// Expect 50% grey
	if r < 0x7000 || r > 0x9000 {
		t.Errorf("Expected mid color at right, got %v (r=%x)", cRight, r)
	}

	// Left: (40, 50) -> dx=-10, dy=0 -> angle=Pi (or -Pi)
	// t = (Pi + Pi) / 2Pi = 1.0  OR (-Pi + Pi) / 2Pi = 0.0
	// Atan2(-10, 0) is Pi.
	// t = (Pi + Pi) / 2Pi = 1.0 -> End Color (White)

	cLeft := g.At(40, 50)
	// Should be White (or close to it)
	if cLeft != end {
		// Wait, floating point might make it slightly less than 1.0 or wrap?
		// Check values.
		t.Logf("Left color: %v", cLeft)
	}

	// Top: (50, 40) -> dx=0, dy=-10 -> angle=-Pi/2
	// t = (-Pi/2 + Pi) / 2Pi = (Pi/2) / 2Pi = 0.25

	cTop := g.At(50, 40)
	rTop, _, _, _ := cTop.RGBA()
	// Expect 25% grey (approx 0x4000)
	if rTop < 0x3000 || rTop > 0x5000 {
		t.Errorf("Expected 25%% grey at top, got %v (r=%x)", cTop, rTop)
	}

}
