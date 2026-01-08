package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure WindRidges implements the image.Image interface.
var _ image.Image = (*WindRidges)(nil)

// WindRidges turns white noise into stretched, wind-swept streaks with soft shadow ridges.
//
// It starts from a noise source (default HashNoise), streaks it along a wind angle, then
// applies a subtle perpendicular shadow to give the impression of depth.
//
// Tunable parameters:
//   - Angle: wind direction in degrees (0° is to the right, 90° is downward).
//   - StreakLength: how far each sample smears along the wind direction.
//   - Contrast: exponent applied to the final luminance for crisper ridges (>1) or softer (<1).
//
// ShadowDistance and ShadowStrength are secondary controls to keep the ridges soft by default.
type WindRidges struct {
	Null
	Noise          image.Image
	Angle          float64
	StreakLength   int
	Contrast       float64
	ShadowDistance float64
	ShadowStrength float64
}

// NewWindRidges builds a wind-swept noise field.
func NewWindRidges(ops ...func(any)) image.Image {
	p := &WindRidges{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Noise:          NewNoise(SetNoiseAlgorithm(&HashNoise{Seed: 404})),
		Angle:          20.0,
		StreakLength:   18,
		Contrast:       1.1,
		ShadowDistance: 2.0,
		ShadowStrength: 0.35,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

func (p *WindRidges) At(x, y int) color.Color {
	streak := p.sampleStreak(x, y)
	contrast := p.Contrast
	if contrast == 0 {
		contrast = 1.0
	}
	if contrast != 1.0 {
		streak = math.Pow(streak, contrast)
	}

	shaded := clamp01(streak + p.shadowRidge(x, y))
	v := uint8(shaded * 255.0)
	return color.RGBA{R: v, G: v, B: v, A: 255}
}

func (p *WindRidges) sampleStreak(x, y int) float64 {
	length := p.StreakLength
	if length < 2 {
		length = 8
	}

	rad := p.angleRad()
	dx := math.Cos(rad)
	dy := math.Sin(rad)

	total := 0.0
	weightSum := 0.0
	decay := 1.0 / float64(length)

	for i := 0; i < length; i++ {
		weight := math.Exp(-float64(i) * decay)
		nx := int(float64(x) - dx*float64(i))
		ny := int(float64(y) - dy*float64(i))
		total += p.sampleNoise(nx, ny) * weight
		weightSum += weight
	}

	if weightSum == 0 {
		return 0
	}
	return total / weightSum
}

func (p *WindRidges) shadowRidge(x, y int) float64 {
	strength := p.ShadowStrength
	if strength == 0 {
		strength = 0.3
	}
	distance := p.ShadowDistance
	if distance == 0 {
		distance = 2.0
	}

	rad := p.angleRad()
	px := -math.Sin(rad)
	py := math.Cos(rad)

	ahead := p.sampleNoise(int(float64(x)+px*distance), int(float64(y)+py*distance))
	behind := p.sampleNoise(int(float64(x)-px*distance), int(float64(y)-py*distance))

	return (ahead - behind) * strength * 0.5
}

func (p *WindRidges) sampleNoise(x, y int) float64 {
	if p.Noise == nil {
		return 0
	}
	gray := color.GrayModel.Convert(p.Noise.At(x, y)).(color.Gray)
	return float64(gray.Y) / 255.0
}

func (p *WindRidges) angleRad() float64 {
	return p.Angle * math.Pi / 180.0
}

// SetWindAngle sets the wind direction in degrees.
func SetWindAngle(angle float64) func(any) {
	return func(i any) {
		if p, ok := i.(*WindRidges); ok {
			p.Angle = angle
		}
	}
}

// SetStreakLength sets how far the streaking extends along the wind.
func SetStreakLength(length int) func(any) {
	return func(i any) {
		if p, ok := i.(*WindRidges); ok {
			p.StreakLength = length
		}
	}
}

// SetWindContrast adjusts the luminance exponent applied after streaking.
func SetWindContrast(contrast float64) func(any) {
	return func(i any) {
		if p, ok := i.(*WindRidges); ok {
			p.Contrast = contrast
		}
	}
}

// SetWindNoise swaps the underlying noise source.
func SetWindNoise(img image.Image) func(any) {
	return func(i any) {
		if p, ok := i.(*WindRidges); ok {
			p.Noise = img
		}
	}
}
