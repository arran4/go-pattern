package pattern

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

// Ensure Noise implements the image.Image interface.
var _ image.Image = (*Noise)(nil)

// Noise generates procedural noise using a native Perlin implementation.
type Noise struct {
	Null
	seed        int64
	alpha       float64 // Persistence
	beta        float64 // Lacunarity
	n           int     // Octaves
	frequency   float64 // Base frequency
	color1      color.Color // Low value color
	color2      color.Color // High value color

	// Internal perlin state
	p []int
}

func (n *Noise) ColorModel() color.Model {
	return color.RGBAModel
}

func (n *Noise) Bounds() image.Rectangle {
	return n.bounds
}

func (n *Noise) At(x, y int) color.Color {
	val := n.fbm(float64(x)*n.frequency, float64(y)*n.frequency)

	// val is roughly -1 to 1. Normalize to 0..1
	norm := (val + 1.0) / 2.0
	if norm < 0 {
		norm = 0
	}
	if norm > 1 {
		norm = 1
	}

	return interpolateColor(n.color1, n.color2, norm)
}

// fbm implements Fractional Brownian Motion
func (n *Noise) fbm(x, y float64) float64 {
	total := 0.0
	frequency := 1.0
	amplitude := 1.0
	maxValue := 0.0 // Used for normalizing result

	for i := 0; i < n.n; i++ {
		total += n.perlin(x*frequency, y*frequency) * amplitude

		maxValue += amplitude

		// Use beta as lacunarity (frequency multiplier)
		// Use alpha as persistence (amplitude multiplier - usually < 1)

		amplitude *= n.alpha
		frequency *= n.beta
	}

	if maxValue > 0 {
		return total / maxValue
	}
	return total
}


// Perlin Noise Implementation

func (n *Noise) initPerlin() {
	r := rand.New(rand.NewSource(n.seed))
	n.p = make([]int, 512)
	permutation := make([]int, 256)
	for i := 0; i < 256; i++ {
		permutation[i] = i
	}

	// Shuffle
	for i := 255; i > 0; i-- {
		j := r.Intn(i + 1)
		permutation[i], permutation[j] = permutation[j], permutation[i]
	}

	for i := 0; i < 256; i++ {
		n.p[i] = permutation[i]
		n.p[i+256] = permutation[i]
	}
}

func (n *Noise) perlin(x, y float64) float64 {
	if n.p == nil {
		n.initPerlin()
	}

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

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y float64) float64 {
	h := hash & 15
	grad := 1.0 + float64(h&7) // Gradient value 1.0, 2.0, ..., 8.0
	if (h & 8) != 0 {
		grad = -grad
	}
	// This is not the standard grad function.
	// Standard Improved Perlin uses vectors (1,1), (-1,1), etc.
	// Let's use the standard switch based one.

	switch h & 0xF {
	case 0x0: return  x + y
	case 0x1: return -x + y
	case 0x2: return  x - y
	case 0x3: return -x - y
	case 0x4: return  x
	case 0x5: return -x
	case 0x6: return  x
	case 0x7: return -x
	case 0x8: return  y
	case 0x9: return -y
	case 0xA: return  y
	case 0xB: return -y
	case 0xC: return  y + x
	case 0xD: return -y + x
	case 0xE: return  y - x
	case 0xF: return -y - x
	default: return 0 // never happens
	}
}

// Configuration options

// Seed configures the seed for the noise generator.
type Seed struct {
	Seed int64
}

func (s *Seed) SetSeed(v int64) {
	s.Seed = v
}

type hasSeed interface {
	SetSeed(int64)
}

func SetSeed(v int64) func(any) {
	return func(i any) {
		if h, ok := i.(hasSeed); ok {
			h.SetSeed(v)
		}
	}
}

// Frequency configures the base frequency for the noise generator.
type Frequency struct {
	Frequency float64
}

func (s *Frequency) SetFrequency(v float64) {
	s.Frequency = v
}

type hasFrequency interface {
	SetFrequency(float64)
}

func SetFrequency(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasFrequency); ok {
			h.SetFrequency(v)
		}
	}
}

// NoiseAlpha configures the persistence for Perlin noise.
type NoiseAlpha struct {
	NoiseAlpha float64
}

func (s *NoiseAlpha) SetNoiseAlpha(v float64) {
	s.NoiseAlpha = v
}

type hasNoiseAlpha interface {
	SetNoiseAlpha(float64)
}

func SetNoiseAlpha(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasNoiseAlpha); ok {
			h.SetNoiseAlpha(v)
		}
	}
}

// NoiseBeta configures the lacunarity for Perlin noise.
type NoiseBeta struct {
	NoiseBeta float64
}

func (s *NoiseBeta) SetNoiseBeta(v float64) {
	s.NoiseBeta = v
}

type hasNoiseBeta interface {
	SetNoiseBeta(float64)
}

func SetNoiseBeta(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasNoiseBeta); ok {
			h.SetNoiseBeta(v)
		}
	}
}

// NoiseN configures the number of octaves for Perlin noise.
type NoiseN struct {
	NoiseN int
}

func (s *NoiseN) SetNoiseN(v int) {
	s.NoiseN = v
}

type hasNoiseN interface {
	SetNoiseN(int)
}

func SetNoiseN(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasNoiseN); ok {
			h.SetNoiseN(v)
		}
	}
}

// NewPerlinNoise creates a new Perlin noise pattern.
func NewPerlinNoise(ops ...func(any)) image.Image {
	p := &Noise{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		seed:      1,
		alpha:     0.5, // Persistence (usually < 1)
		beta:      2.0, // Lacunarity (usually > 1)
		n:         3,
		frequency: 0.1,
		color1:    color.Black,
		color2:    color.White,
	}
	for _, op := range ops {
		op(p)
	}
	// Init perlin immediately? Or lazy?
	// Lazy is safer for options setting seed.
	// But let's init here if seed is final, or re-init if seed changes?
	// The options are applied before this return.
	p.initPerlin()
	return p
}

// Implement option interfaces
func (n *Noise) SetSeed(v int64) {
	n.seed = v
	n.initPerlin() // Re-init if seed changes
}

func (n *Noise) SetFrequency(v float64) {
	n.frequency = v
}

func (n *Noise) SetNoiseAlpha(v float64) {
	n.alpha = v
}

func (n *Noise) SetNoiseBeta(v float64) {
	n.beta = v
}

func (n *Noise) SetNoiseN(v int) {
	n.n = v
}

type NoiseColorLow struct {
	NoiseColorLow color.Color
}

func (s *NoiseColorLow) SetNoiseColorLow(v color.Color) {
	s.NoiseColorLow = v
}

type hasNoiseColorLow interface {
	SetNoiseColorLow(color.Color)
}

func SetNoiseColorLow(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasNoiseColorLow); ok {
			h.SetNoiseColorLow(v)
		}
	}
}

type NoiseColorHigh struct {
	NoiseColorHigh color.Color
}

func (s *NoiseColorHigh) SetNoiseColorHigh(v color.Color) {
	s.NoiseColorHigh = v
}

type hasNoiseColorHigh interface {
	SetNoiseColorHigh(color.Color)
}

func SetNoiseColorHigh(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasNoiseColorHigh); ok {
			h.SetNoiseColorHigh(v)
		}
	}
}

func (n *Noise) SetNoiseColorLow(v color.Color) {
	n.color1 = v
}

func (n *Noise) SetNoiseColorHigh(v color.Color) {
	n.color2 = v
}
