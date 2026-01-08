package pattern

import (
	"crypto/rand"
	"image"
	"image/color"
	"math"
	mrand "math/rand"
	"sync"
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

// SetSeedUint64 sets the seed for the noise algorithm.
// It switches to HashNoise if the current algo is CryptoNoise.
func (n *Noise) SetSeedUint64(v uint64) {
	switch algo := n.algo.(type) {
	case *CryptoNoise:
		n.algo = &HashNoise{Seed: int64(v)}
	case *HashNoise:
		algo.Seed = int64(v)
	case *PerlinNoise:
		algo.Seed = int64(v)
	}
}

// SetSeed sets the seed for the noise algorithm.
// It switches to HashNoise if the current algo is CryptoNoise.
func (n *Noise) SetSeed(v int64) {
	switch algo := n.algo.(type) {
	case *CryptoNoise:
		n.algo = &HashNoise{Seed: v}
	case *HashNoise:
		algo.Seed = v
	case *PerlinNoise:
		algo.Seed = v
	}
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

// NoiseSeed sets the seed for the noise algorithm.
// If the current algorithm is CryptoNoise, it switches to HashNoise.
//
// Deprecated: Use Seed() instead.
func NoiseSeed(seed int64) func(any) {
	return func(i any) {
		if n, ok := i.(*Noise); ok {
			n.SetSeed(seed)
		}
	}
}

// --- Algorithms ---

// CryptoNoise uses crypto/rand.
type CryptoNoise struct{}

func (c *CryptoNoise) At(x, y int) color.Color {
	b := make([]byte, 1)
	if _, err := rand.Read(b); err != nil {
		return color.Black
	}
	return color.Gray{Y: b[0]}
}

// HashNoise uses a high-quality, stateless pseudo-random number generator based on coordinates.
// It produces a non-repeating pattern derived algorithmically.
type HashNoise struct {
	Seed int64
}

func (h *HashNoise) At(x, y int) color.Color {
	// Mix x, y, and Seed using a robust hash function.
	z := StableHash(x, y, uint64(h.Seed))
	return color.Gray{Y: uint8(z)}
}

// PerlinNoise implements Improved Perlin Noise with Fractional Brownian Motion (fBm).
type PerlinNoise struct {
	Seed        int64
	Octaves     int
	Persistence float64 // Alpha
	Lacunarity  float64 // Beta
	Frequency   float64

	p    [512]int
	once sync.Once
}

func (n *PerlinNoise) init() {
	n.once.Do(func() {
		// Default parameters if zero
		if n.Frequency == 0 {
			n.Frequency = 0.02
		}
		if n.Octaves == 0 {
			n.Octaves = 1
		}
		if n.Lacunarity == 0 {
			n.Lacunarity = 2.0
		}
		if n.Persistence == 0 {
			n.Persistence = 0.5
		}

		r := mrand.New(mrand.NewSource(n.Seed))
		perm := make([]int, 256)
		for i := range perm {
			perm[i] = i
		}
		r.Shuffle(len(perm), func(i, j int) {
			perm[i], perm[j] = perm[j], perm[i]
		})

		for i := 0; i < 256; i++ {
			n.p[i] = perm[i]
			n.p[i+256] = perm[i]
		}
	})
}

func (n *PerlinNoise) At(x, y int) color.Color {
	n.init()

	var total float64
	var maxAmplitude float64
	amplitude := 1.0
	frequency := n.Frequency

	for i := 0; i < n.Octaves; i++ {
		total += n.noise(float64(x)*frequency, float64(y)*frequency) * amplitude
		maxAmplitude += amplitude
		amplitude *= n.Persistence
		frequency *= n.Lacunarity
	}

	// Normalize result to [0, 1]
	// Perlin noise returns values roughly in [-1, 1]
	val := total / maxAmplitude

	// Map [-1, 1] to [0, 1]
	normalized := (val + 1.0) * 0.5
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	c := uint8(normalized * 255)
	return color.Gray{Y: c}
}

func (n *PerlinNoise) noise(x, y float64) float64 {
	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255

	x -= math.Floor(x)
	y -= math.Floor(y)

	u := fade(x)
	v := fade(y)

	A := n.p[X] + Y
	B := n.p[X+1] + Y

	return lerp(v, lerp(u, grad(n.p[A], x, y), grad(n.p[B], x-1, y)),
		lerp(u, grad(n.p[A+1], x, y-1), grad(n.p[B+1], x-1, y-1)))
}

// Sample returns the Perlin noise value in the range [0,1] at floating-point coordinates.
// This allows downstream patterns to query the same noise with custom coordinate transforms.
func (n *PerlinNoise) Sample(x, y float64) float64 {
	n.init()

	total := 0.0
	maxAmplitude := 0.0
	amplitude := 1.0
	frequency := n.Frequency

	for i := 0; i < n.Octaves; i++ {
		total += n.noise(x*frequency, y*frequency) * amplitude
		maxAmplitude += amplitude
		amplitude *= n.Persistence
		frequency *= n.Lacunarity
	}

	if maxAmplitude == 0 {
		return 0.5
	}

	normalized := (total/maxAmplitude + 1.0) * 0.5
	if normalized < 0 {
		return 0
	}
	if normalized > 1 {
		return 1
	}
	return normalized
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y float64) float64 {
	h := hash & 7
	switch h {
	case 0:
		return x + y
	case 1:
		return -x + y
	case 2:
		return x - y
	case 3:
		return -x - y
	case 4:
		return x
	case 5:
		return -x
	case 6:
		return y
	case 7:
		return -y
	default:
		return 0
	}
}
