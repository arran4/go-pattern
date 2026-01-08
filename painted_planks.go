package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure PaintedPlanks implements the image.Image interface.
var _ image.Image = (*PaintedPlanks)(nil)

// PaintedPlanks composes a set of vertical boards with wood grain and
// a chipped paint overlay.
type PaintedPlanks struct {
	Null

	BaseWidth      int
	WidthVariance  float64
	GrainIntensity float64
	PaintWear      float64

	PaintColor color.RGBA
	SeamColor  color.RGBA
	WoodDark   color.RGBA
	WoodLight  color.RGBA

	Seed

	grainNoise *PerlinNoise
	fiberNoise *PerlinNoise
	paintNoise *PerlinNoise
	chipNoise  *PerlinNoise
}

func (p *PaintedPlanks) At(x, y int) color.Color {
	baseWidth := p.BaseWidth
	if baseWidth <= 4 {
		baseWidth = 60
	}

	plankIndex, start, end := p.plankInfo(x, float64(baseWidth))
	width := end - start
	if width <= 0 {
		width = 1
	}
	localX := float64(x) - start
	localT := clamp01(localX / width)

	seamDistance := math.Min(localX, width-localX)
	if seamDistance < 1.0 {
		return p.SeamColor
	}

	grainValue := p.sampleNoise(p.grainNoise, x+plankIndex*17, y)
	fiberValue := p.sampleNoise(p.fiberNoise, x+plankIndex*11, y*3+plankIndex*5)
	grainMix := clamp01(0.55 + (grainValue-0.5)*2*p.GrainIntensity*0.6 + (fiberValue-0.5)*0.35)
	plankTint := p.plankTint(plankIndex)
	woodMix := clamp01(grainMix + plankTint)
	wood := lerpRGBA(p.WoodDark, p.WoodLight, woodMix)

	paintPresence := p.paintMask(localT, plankIndex, x, y, seamDistance/width)
	finalColor := lerpRGBA(wood, p.PaintColor, paintPresence)
	return finalColor
}

func (p *PaintedPlanks) plankInfo(x int, baseWidth float64) (int, float64, float64) {
	idx := int(math.Floor(float64(x) / baseWidth))
	for attempts := 0; attempts < 6; attempts++ {
		start := float64(idx)*baseWidth + p.boundaryJitter(idx, baseWidth)
		end := float64(idx+1)*baseWidth + p.boundaryJitter(idx+1, baseWidth)
		if end <= start {
			end = start + 1
		}
		if float64(x) >= start && float64(x) < end {
			return idx, start, end
		}
		if float64(x) < start {
			idx--
		} else {
			idx++
		}
	}
	start := float64(idx) * baseWidth
	return idx, start, start + baseWidth
}

func (p *PaintedPlanks) boundaryJitter(idx int, baseWidth float64) float64 {
	amp := clamp01(p.WidthVariance) * 0.5 * baseWidth
	h := StableHash(idx, 0, uint64(p.Seed.Seed))
	n := (float64(int64(h&0xffff)) - 32768.0) / 32768.0
	return n * amp
}

func (p *PaintedPlanks) plankTint(idx int) float64 {
	h := StableHash(idx, 1, uint64(p.Seed.Seed))
	// Map to [-0.08, 0.08]
	return (float64(int64(h&0xffff)) - 32768.0) / 32768.0 * 0.08
}

func (p *PaintedPlanks) paintMask(localT float64, plankIndex, x, y int, edgeRatio float64) float64 {
	baseWear := clamp01(p.PaintWear)

	paintTexture := p.sampleNoise(p.paintNoise, x+plankIndex*23, y)
	chips := p.sampleNoise(p.chipNoise, x*2+plankIndex*13, y*2+7*plankIndex)
	edgeWear := clamp01(1.0 - edgeRatio*2)

	coverage := 1 - baseWear
	coverage += (paintTexture - 0.5) * 0.45
	coverage += (chips - 0.5) * 0.35
	coverage -= edgeWear * 0.35

	return clamp01(coverage)
}

func (p *PaintedPlanks) sampleNoise(n *PerlinNoise, x, y int) float64 {
	if n == nil {
		return 0.5
	}
	c := n.At(x, y)
	r, _, _, _ := c.RGBA()
	return float64(r) / 65535.0
}

func lerpFloat(a, b, t float64) float64 {
	return a + (b-a)*t
}

// NewPaintedPlanks creates a plank wall pattern with configurable wood grain and paint wear.
func NewPaintedPlanks(ops ...func(any)) image.Image {
	p := &PaintedPlanks{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		BaseWidth:      70,
		WidthVariance:  0.25,
		GrainIntensity: 0.65,
		PaintWear:      0.35,
		PaintColor:     color.RGBA{190, 205, 219, 255},
		SeamColor:      color.RGBA{50, 40, 32, 255},
		WoodDark:       color.RGBA{104, 76, 55, 255},
		WoodLight:      color.RGBA{193, 160, 120, 255},
		Seed:           Seed{Seed: 1337},
	}
	for _, op := range ops {
		op(p)
	}
	p.configureNoises()
	return p
}

func (p *PaintedPlanks) configureNoises() {
	base := p.Seed.Seed
	p.grainNoise = &PerlinNoise{
		Seed:        base + 11,
		Frequency:   0.08,
		Octaves:     4,
		Persistence: 0.55,
		Lacunarity:  2.2,
	}
	p.fiberNoise = &PerlinNoise{
		Seed:        base + 19,
		Frequency:   0.18,
		Octaves:     3,
		Persistence: 0.6,
		Lacunarity:  2.1,
	}
	p.paintNoise = &PerlinNoise{
		Seed:        base + 29,
		Frequency:   0.025,
		Octaves:     3,
		Persistence: 0.6,
		Lacunarity:  2.0,
	}
	p.chipNoise = &PerlinNoise{
		Seed:        base + 37,
		Frequency:   0.11,
		Octaves:     2,
		Persistence: 0.65,
		Lacunarity:  2.3,
	}
}

// Options

type hasPlankBaseWidth interface{ SetPlankBaseWidth(int) }

func SetPlankBaseWidth(w int) func(any) {
	return func(i any) {
		if h, ok := i.(hasPlankBaseWidth); ok {
			h.SetPlankBaseWidth(w)
		}
	}
}

type hasPlankWidthVariance interface{ SetPlankWidthVariance(float64) }

func SetPlankWidthVariance(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasPlankWidthVariance); ok {
			h.SetPlankWidthVariance(v)
		}
	}
}

type hasGrainIntensity interface{ SetGrainIntensity(float64) }

func SetGrainIntensity(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasGrainIntensity); ok {
			h.SetGrainIntensity(v)
		}
	}
}

type hasPaintWear interface{ SetPaintWear(float64) }

func SetPaintWear(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasPaintWear); ok {
			h.SetPaintWear(v)
		}
	}
}

type hasPaintColor interface{ SetPaintColor(color.RGBA) }

func SetPaintColor(c color.RGBA) func(any) {
	return func(i any) {
		if h, ok := i.(hasPaintColor); ok {
			h.SetPaintColor(c)
		}
	}
}

// Implement setters
func (p *PaintedPlanks) SetPlankBaseWidth(w int)         { p.BaseWidth = w }
func (p *PaintedPlanks) SetPlankWidthVariance(v float64) { p.WidthVariance = v }
func (p *PaintedPlanks) SetGrainIntensity(v float64)     { p.GrainIntensity = v }
func (p *PaintedPlanks) SetPaintWear(v float64)          { p.PaintWear = v }
func (p *PaintedPlanks) SetPaintColor(c color.RGBA)      { p.PaintColor = c }

