package pattern

import (
	"testing"
)

func TestNoiseStability(t *testing.T) {
	// Create two noise patterns with the same seed
	n1 := NewNoise(NoiseSeed(12345))
	n2 := NewNoise(NoiseSeed(12345))

	// Check a few pixels to ensure they are identical
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			c1 := n1.At(x, y)
			c2 := n2.At(x, y)
			if c1 != c2 {
				t.Errorf("Pixels at (%d, %d) differ: %v vs %v", x, y, c1, c2)
			}
		}
	}
}

func TestNoiseVariety(t *testing.T) {
	// Create two noise patterns with different seeds
	n1 := NewNoise(NoiseSeed(12345))
	n2 := NewNoise(NoiseSeed(67890))

	// Check a few pixels to ensure they are NOT all identical
	// (There is a small chance of collision for a single pixel, but unlikely for all 100)
	identicalCount := 0
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			c1 := n1.At(x, y)
			c2 := n2.At(x, y)
			if c1 == c2 {
				identicalCount++
			}
		}
	}

	if identicalCount == 100 {
		t.Errorf("All 100 tested pixels were identical despite different seeds")
	}
}

func TestNoiseAlgoSwitch(t *testing.T) {
	// Default NewNoise uses CryptoNoise
	n := NewNoise().(*Noise)
	if _, ok := n.algo.(*CryptoNoise); !ok {
		t.Errorf("Expected default algo to be CryptoNoise, got %T", n.algo)
	}

	// Applying NoiseSeed should switch to HashNoise
	NoiseSeed(123)(n)
	if _, ok := n.algo.(*HashNoise); !ok {
		t.Errorf("Expected algo to switch to HashNoise after seeding, got %T", n.algo)
	}
	if h, ok := n.algo.(*HashNoise); ok {
		if h.Seed != 123 {
			t.Errorf("Expected seed 123, got %d", h.Seed)
		}
	}
}

func TestPerlinNoiseSeed(t *testing.T) {
	n := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Seed: 1})).(*Noise)

	// Apply new seed
	NoiseSeed(999)(n)

	if p, ok := n.algo.(*PerlinNoise); ok {
		if p.Seed != 999 {
			t.Errorf("Expected Perlin seed to update to 999, got %d", p.Seed)
		}
	} else {
		t.Errorf("Expected algo to remain PerlinNoise, got %T", n.algo)
	}
}
