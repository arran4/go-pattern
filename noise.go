package pattern

import (
	"image"
	"image/color"

	"github.com/aquilax/go-perlin"
	"github.com/ojrac/opensimplex-go"
)

// Ensure Noise implements the image.Image interface.
var _ image.Image = (*Noise)(nil)

// NoiseType distinguishes between different noise algorithms.
type NoiseType int

const (
	PerlinNoise NoiseType = iota
	OpenSimplexNoise
)

// Noise generates procedural noise using Perlin or OpenSimplex algorithms.
type Noise struct {
	Null
	noiseType   NoiseType
	seed        int64
	alpha       float64 // For Perlin: weight when summing octaves (persistence)
	beta        float64 // For Perlin: frequency multiplier (lacunarity)
	n           int     // For Perlin: number of octaves
	frequency   float64 // Base frequency
	perlinGen   *perlin.Perlin
	simplexGen  opensimplex.Noise
	color1      color.Color // Low value color (approx -1)
	color2      color.Color // High value color (approx 1)
}

func (n *Noise) ColorModel() color.Model {
	return color.RGBAModel
}

func (n *Noise) Bounds() image.Rectangle {
	return n.bounds
}

func (n *Noise) At(x, y int) color.Color {
	var val float64

	switch n.noiseType {
	case PerlinNoise:
		// go-perlin expects float coordinates. We scale by frequency.
		val = n.perlinGen.Noise2D(float64(x)*n.frequency, float64(y)*n.frequency)
	case OpenSimplexNoise:
		val = n.simplexGen.Eval2(float64(x)*n.frequency, float64(y)*n.frequency)
	}

	// val is roughly -1 to 1. Normalize to 0..1 for color interpolation
	norm := (val + 1.0) / 2.0
	if norm < 0 {
		norm = 0
	}
	if norm > 1 {
		norm = 1
	}

	return interpolateColor(n.color1, n.color2, norm)
}

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

// NoiseAlpha configures the alpha (persistence) for Perlin noise.
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

// NoiseBeta configures the beta (lacunarity) for Perlin noise.
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

// NoiseN configures the number of iterations (octaves) for Perlin noise.
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
		noiseType: PerlinNoise,
		seed:      1,
		alpha:     2,
		beta:      2,
		n:         3,
		frequency: 0.1,
		color1:    color.Black,
		color2:    color.White,
	}
	for _, op := range ops {
		op(p)
	}
	p.perlinGen = perlin.NewPerlin(p.alpha, p.beta, int32(p.n), p.seed)
	return p
}

// NewSimplexNoise creates a new OpenSimplex noise pattern.
func NewSimplexNoise(ops ...func(any)) image.Image {
	p := &Noise{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		noiseType: OpenSimplexNoise,
		seed:      1,
		frequency: 0.1,
		color1:    color.Black,
		color2:    color.White,
	}
	for _, op := range ops {
		op(p)
	}
	p.simplexGen = opensimplex.New(p.seed)
	return p
}

// Implement option interfaces
func (n *Noise) SetSeed(v int64) {
	n.seed = v
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
