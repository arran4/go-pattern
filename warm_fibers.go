package pattern

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// Ensure WarmFibersGradient implements the image.Image interface.
var _ image.Image = (*WarmFibersGradient)(nil)

// WarmFibersGradient blends a warm vertical gradient with subtle fibrous noise and vignette.
type WarmFibersGradient struct {
	Null
	StartColor
	EndColor
	Seed

	FiberDensity     float64
	VignetteStrength float64

	noiseStrength float64
	perlin        *PerlinNoise
}

// NewWarmFibersGradient creates a new WarmFibersGradient pattern.
func NewWarmFibersGradient(ops ...func(any)) image.Image {
	w := &WarmFibersGradient{
		Null:       Null{bounds: image.Rect(0, 0, 255, 255)},
		StartColor: StartColor{StartColor: color.RGBA{240, 192, 150, 255}},
		EndColor:   EndColor{EndColor: color.RGBA{143, 72, 45, 255}},
		Seed:       Seed{Seed: 7},

		FiberDensity:     1.0,
		VignetteStrength: 0.35,
		noiseStrength:    0.14,
		perlin: &PerlinNoise{
			Frequency:   0.02,
			Octaves:     4,
			Persistence: 0.55,
			Lacunarity:  2.1,
		},
	}

	for _, op := range ops {
		op(w)
	}

	w.syncNoiseConfig()

	return w
}

// SetSeed ensures the internal noise uses the provided seed.
func (w *WarmFibersGradient) SetSeed(v int64) {
	w.Seed.Seed = v
	w.ensurePerlin()
	w.perlin.Seed = v
	w.perlin.once = sync.Once{}
}

// SetSeedUint64 ensures the internal noise uses the provided seed.
func (w *WarmFibersGradient) SetSeedUint64(v uint64) {
	w.SetSeed(int64(v))
}

// SetFiberDensity adjusts the density of vertical fibers.
func (w *WarmFibersGradient) SetFiberDensity(v float64) {
	if v <= 0 {
		v = 0.1
	}
	w.FiberDensity = v
}

// SetVignetteStrength controls the vignette falloff.
func (w *WarmFibersGradient) SetVignetteStrength(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	w.VignetteStrength = v
}

// At returns the color at (x, y).
func (w *WarmFibersGradient) At(x, y int) color.Color {
	b := w.Bounds()
	if b.Empty() {
		return color.RGBA{}
	}

	t := 0.0
	if b.Dy() > 1 {
		t = float64(y-b.Min.Y) / float64(b.Dy()-1)
	}

	base := lerpColor(w.StartColor.StartColor, w.EndColor.EndColor, t)
	br, bg, bb, ba := base.RGBA()
	r := float64(br) / 65535
	g := float64(bg) / 65535
	bv := float64(bb) / 65535
	a := float64(ba) / 65535

	noiseVal := w.sampleNoise(float64(x), float64(y))
	fiberShift := (noiseVal - 0.5) * w.noiseStrength
	vignette := w.vignetteFactor(x, y)

	factor := clampFloat64((1+fiberShift)*vignette, 0, 1.5)

	outR := clampFloat64(r*factor, 0, 1)
	outG := clampFloat64(g*factor, 0, 1)
	outB := clampFloat64(bv*factor, 0, 1)

	return color.RGBA64{
		R: uint16(outR * 65535),
		G: uint16(outG * 65535),
		B: uint16(outB * 65535),
		A: uint16(a * 65535),
	}
}

func (w *WarmFibersGradient) sampleNoise(x, y float64) float64 {
	w.syncNoiseConfig()
	baseFreq := w.perlin.Frequency
	if baseFreq == 0 {
		baseFreq = 0.02
	}
	xFreq := baseFreq * w.FiberDensity
	if xFreq <= 0 {
		xFreq = baseFreq
	}
	yFreq := baseFreq * 0.35

	octaves := w.perlin.Octaves
	if octaves == 0 {
		octaves = 1
	}
	persistence := w.perlin.Persistence
	if persistence == 0 {
		persistence = 0.5
	}
	lacunarity := w.perlin.Lacunarity
	if lacunarity == 0 {
		lacunarity = 2.0
	}

	var total float64
	var maxAmplitude float64
	amplitude := 1.0

	for i := 0; i < octaves; i++ {
		freqMul := math.Pow(lacunarity, float64(i))
		total += w.perlin.noise(x*xFreq*freqMul, y*yFreq*freqMul) * amplitude
		maxAmplitude += amplitude
		amplitude *= persistence
	}

	val := total / maxAmplitude
	normalized := (val + 1.0) * 0.5
	return clampFloat64(normalized, 0, 1)
}

func (w *WarmFibersGradient) vignetteFactor(x, y int) float64 {
	if w.VignetteStrength <= 0 {
		return 1
	}
	b := w.Bounds()
	cx := float64(b.Min.X) + float64(b.Dx())/2
	cy := float64(b.Min.Y) + float64(b.Dy())/2
	dx := float64(x) - cx
	dy := float64(y) - cy
	dist := math.Sqrt(dx*dx + dy*dy)
	maxDist := math.Sqrt(float64(b.Dx()*b.Dx()+b.Dy()*b.Dy())) / 2
	if maxDist == 0 {
		return 1
	}
	ratio := dist / maxDist
	if ratio > 1 {
		ratio = 1
	}
	factor := 1 - w.VignetteStrength*ratio*ratio
	return clampFloat64(factor, 0, 1)
}

func (w *WarmFibersGradient) ensurePerlin() {
	if w.perlin == nil {
		w.perlin = &PerlinNoise{}
	}
}

func (w *WarmFibersGradient) syncNoiseConfig() {
	w.ensurePerlin()

	if w.perlin.Frequency == 0 {
		w.perlin.Frequency = 0.02
	}
	if w.perlin.Octaves == 0 {
		w.perlin.Octaves = 4
	}
	if w.perlin.Persistence == 0 {
		w.perlin.Persistence = 0.55
	}
	if w.perlin.Lacunarity == 0 {
		w.perlin.Lacunarity = 2.1
	}

	w.perlin.Seed = w.Seed.Seed
}

func clampFloat64(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// SetFiberDensity creates an option to control fiber density.
func SetFiberDensity(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetFiberDensity(float64) }); ok {
			h.SetFiberDensity(v)
		}
	}
}

// SetVignetteStrength creates an option to control vignette falloff.
func SetVignetteStrength(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetVignetteStrength(float64) }); ok {
			h.SetVignetteStrength(v)
		}
	}
}
