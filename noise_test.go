package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestPerlinNoise_Consistency(t *testing.T) {
	n1 := NewPerlinNoise(SetSeed(123), SetFrequency(0.1))
	n2 := NewPerlinNoise(SetSeed(123), SetFrequency(0.1))

	c1 := n1.At(10, 10)
	c2 := n2.At(10, 10)

	if c1 != c2 {
		t.Errorf("Expected same color for same seed, got %v and %v", c1, c2)
	}
}

func TestPerlinNoise_SeedDifference(t *testing.T) {
	n1 := NewPerlinNoise(SetSeed(123), SetFrequency(0.1))
	n2 := NewPerlinNoise(SetSeed(456), SetFrequency(0.1))

	// It's possible but unlikely that two seeds produce the exact same value at (10,10)
	// But let's check multiple points or assume they are different.
	// We'll check 5 points.
	same := true
	for i := 0; i < 5; i++ {
		if n1.At(i, i) != n2.At(i, i) {
			same = false
			break
		}
	}
	if same {
		t.Errorf("Expected different colors for different seeds")
	}
}

func TestNoise_Options(t *testing.T) {
	// Test that options are applied
	n := NewPerlinNoise(
		SetNoiseAlpha(0.5),
		SetNoiseBeta(1.5),
		SetNoiseN(5),
		SetFrequency(0.05),
		SetNoiseColorLow(color.RGBA{255, 0, 0, 255}),
		SetNoiseColorHigh(color.RGBA{0, 255, 0, 255}),
	).(*Noise)

	if n.alpha != 0.5 {
		t.Errorf("Expected alpha 0.5, got %v", n.alpha)
	}
	if n.beta != 1.5 {
		t.Errorf("Expected beta 1.5, got %v", n.beta)
	}
	if n.n != 5 {
		t.Errorf("Expected n 5, got %v", n.n)
	}
	if n.frequency != 0.05 {
		t.Errorf("Expected frequency 0.05, got %v", n.frequency)
	}
	if n.color1 != (color.RGBA{255, 0, 0, 255}) {
		t.Errorf("Expected color1 Red, got %v", n.color1)
	}
	if n.color2 != (color.RGBA{0, 255, 0, 255}) {
		t.Errorf("Expected color2 Green, got %v", n.color2)
	}
}

func TestNoise_Bounds(t *testing.T) {
	n := NewPerlinNoise(SetBounds(image.Rect(0, 0, 100, 100)))
	if n.Bounds() != image.Rect(0, 0, 100, 100) {
		t.Errorf("Expected bounds (0,0,100,100), got %v", n.Bounds())
	}
}
