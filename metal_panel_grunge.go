package pattern

import (
	"image"
	"image/color"
	"math"
)

// MetalPanelGrunge renders brushed metal panels with grunge darkening around seams.
// It exposes controls for streak direction (degrees), scratch density, and seam spacing.
type MetalPanelGrunge struct {
	Null
	Seed

	Direction      float64
	ScratchDensity float64
	SeamSpacing    int

	baseNoise    *PerlinNoise
	scratchNoise *PerlinNoise
	grungeNoise  *PerlinNoise
}

func (p *MetalPanelGrunge) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *MetalPanelGrunge) SetSeed(v int64) {
	p.Seed.Seed = v
	p.initNoise()
}

func (p *MetalPanelGrunge) SetSeedUint64(v uint64) {
	p.Seed.Seed = int64(v)
	p.initNoise()
}

func (p *MetalPanelGrunge) SetStreakDirection(deg float64) {
	p.Direction = deg
}

func (p *MetalPanelGrunge) SetScratchDensity(d float64) {
	p.ScratchDensity = d
}

func (p *MetalPanelGrunge) SetSeamSpacing(px int) {
	p.SeamSpacing = px
}

func (p *MetalPanelGrunge) initNoise() {
	seed := p.Seed.Seed
	p.baseNoise = &PerlinNoise{
		Seed:        seed,
		Frequency:   0.08,
		Octaves:     4,
		Persistence: 0.55,
		Lacunarity:  2.0,
	}
	p.scratchNoise = &PerlinNoise{
		Seed:        seed + 17,
		Frequency:   0.35,
		Octaves:     3,
		Persistence: 0.5,
		Lacunarity:  2.4,
	}
	p.grungeNoise = &PerlinNoise{
		Seed:        seed + 71,
		Frequency:   0.18,
		Octaves:     2,
		Persistence: 0.6,
		Lacunarity:  2.2,
	}
}

func (p *MetalPanelGrunge) At(x, y int) color.Color {
	if p.baseNoise == nil || p.scratchNoise == nil || p.grungeNoise == nil {
		p.initNoise()
	}

	b := p.Bounds()
	fx := float64(x - b.Min.X)
	fy := float64(y - b.Min.Y)

	dir := p.Direction * math.Pi / 180.0
	cosD := math.Cos(dir)
	sinD := math.Sin(dir)

	// Rotate the sampling frame so streaks follow the requested direction.
	along := fx*cosD + fy*sinD
	across := -fx*sinD + fy*cosD

	basePrimary := p.baseNoise.Sample(across*1.2, along*0.15)
	baseSecondary := p.baseNoise.Sample(across*0.25, along*0.02)
	baseValue := clamp01(0.45 + 0.4*basePrimary + 0.15*(baseSecondary-0.5))

	// Hairline scratches aligned with the brushing direction.
	density := clamp01(p.ScratchDensity)
	scratchVal := p.scratchNoise.Sample(across*2.5, along*0.35)
	scratchThreshold := 1.0 - 0.35*density
	denom := 1.0 - scratchThreshold
	if denom < 1e-5 {
		denom = 1e-5
	}
	scratchAlpha := clamp01((scratchVal - scratchThreshold) / denom)
	scratched := baseValue*(1-0.25*scratchAlpha) + 0.05*scratchAlpha

	// Seam mask with grunge modulation.
	spacing := float64(p.SeamSpacing)
	if spacing <= 0 {
		spacing = 140
	}
	vx := math.Mod(math.Abs(fx), spacing)
	vy := math.Mod(math.Abs(fy), spacing)
	distV := math.Min(vx, spacing-vx)
	distH := math.Min(vy, spacing-vy)
	dist := math.Min(distV, distH)
	edgeWidth := math.Max(2.5, spacing*0.02)
	seamMask := math.Exp(-dist * dist / (2 * edgeWidth * edgeWidth))

	grungeVal := p.grungeNoise.Sample(float64(x)*0.4, float64(y)*0.4)
	seamGrunge := clamp01(seamMask * (0.65 + 0.7*(grungeVal-0.5)))

	// Combine layers.
	finalVal := scratched * (1 - 0.35*seamGrunge)
	finalVal = clamp01(finalVal + 0.05*math.Sin(along*0.02))

	sheen := p.baseNoise.Sample(across*0.5, along*0.05) - 0.5

	r := clamp01(finalVal*0.95 + 0.04*sheen)
	g := clamp01(finalVal*1.00 + 0.02*sheen)
	bc := clamp01(finalVal*1.05 + 0.06*sheen)

	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(bc * 255),
		A: 255,
	}
}

// Options
type hasStreakDirection interface{ SetStreakDirection(float64) }

func SetStreakDirection(deg float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasStreakDirection); ok {
			h.SetStreakDirection(deg)
		}
	}
}

type hasScratchDensity interface{ SetScratchDensity(float64) }

func SetScratchDensity(d float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasScratchDensity); ok {
			h.SetScratchDensity(d)
		}
	}
}

type hasSeamSpacing interface{ SetSeamSpacing(int) }

func SetSeamSpacing(px int) func(any) {
	return func(i any) {
		if h, ok := i.(hasSeamSpacing); ok {
			h.SetSeamSpacing(px)
		}
	}
}

// NewMetalPanelGrunge creates the pattern with defaults tailored for brushed metal panels.
func NewMetalPanelGrunge(ops ...func(any)) image.Image {
	p := &MetalPanelGrunge{
		Null:           Null{bounds: image.Rect(0, 0, 255, 255)},
		Direction:      8,
		ScratchDensity: 0.4,
		SeamSpacing:    140,
		Seed:           Seed{Seed: 1337},
	}
	p.initNoise()
	for _, op := range ops {
		op(p)
	}
	return p
}
