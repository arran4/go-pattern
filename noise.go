package pattern

import (
	"crypto/rand"
	"image"
	"image/color"
)

// Ensure Noise implements the image.Image interface.
var _ image.Image = (*Noise)(nil)

// NoiseAlgorithm defines the source of randomness for the Noise pattern.
type NoiseAlgorithm interface {
	At(x, y int) color.Color
}

// Noise pattern generates random noise based on a selected algorithm.
type Noise struct {
	Null
	algo NoiseAlgorithm
}

func (n *Noise) At(x, y int) color.Color {
	if n.algo != nil {
		return n.algo.At(x, y)
	}
	return color.RGBA{0, 0, 0, 255}
}

// NewNoise creates a new Noise pattern.
func NewNoise(ops ...func(any)) image.Image {
	p := &Noise{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		algo: &CryptoNoise{}, // Default to crypto noise
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NoiseAlgorithm configuration
type hasNoiseAlgorithm interface {
	SetNoiseAlgorithm(NoiseAlgorithm)
}

func (n *Noise) SetNoiseAlgorithm(algo NoiseAlgorithm) {
	n.algo = algo
}

func SetNoiseAlgorithm(algo NoiseAlgorithm) func(any) {
	return func(i any) {
		if n, ok := i.(hasNoiseAlgorithm); ok {
			n.SetNoiseAlgorithm(algo)
		}
	}
}

// --- Algorithms ---

// CryptoNoise uses crypto/rand.
type CryptoNoise struct{}

func (c *CryptoNoise) At(x, y int) color.Color {
	b := make([]byte, 1)
	rand.Read(b)
	return color.Gray{Y: b[0]}
}

// HashNoise uses a high-quality, stateless pseudo-random number generator based on coordinates.
// It produces a non-repeating pattern derived algorithmically.
type HashNoise struct {
	Seed int64
}

func (h *HashNoise) At(x, y int) color.Color {
	// Mix x, y, and Seed using a robust hash function (SplitMix64-like steps)
	// to ensure good visual randomness without repetition.
	z := uint64(int64(x)*0x9e3779b9 + int64(y)*0x632be59b + h.Seed)
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	z = z ^ (z >> 31)
	return color.Gray{Y: uint8(z)}
}
