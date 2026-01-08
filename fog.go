package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Fog implements the image.Image interface.
var _ image.Image = (*Fog)(nil)

// Fog renders a tinted, noise-driven fog with radial falloff.
// The center stays clearer than the edges to focus the viewer's eye.
type Fog struct {
	Null
	FillColor
	Density
	FloatCenter
	FalloffCurve

	useFloatCenter bool
	algo           NoiseAlgorithm
}

func (f *Fog) ColorModel() color.Model {
	return color.NRGBAModel
}

func (f *Fog) At(x, y int) color.Color {
	b := f.Bounds()
	if b.Empty() {
		return color.NRGBA{}
	}

	noiseVal := 0.0
	if f.algo != nil {
		// Convert to grayscale to normalize arbitrary NoiseAlgorithms.
		g := color.GrayModel.Convert(f.algo.At(x, y)).(color.Gray)
		noiseVal = float64(g.Y) / 255.0
	}

	radial := f.radialWeight(x, y)
	strength := clampFloat(noiseVal * f.Density.Density * radial)

	tr, tg, tb, ta := f.FillColor.FillColor.RGBA()
	tintR := float64(tr) / 65535.0
	tintG := float64(tg) / 65535.0
	tintB := float64(tb) / 65535.0
	tintA := float64(ta) / 65535.0

	return color.NRGBA{
		R: uint8(clampFloat(tintR*strength) * 255),
		G: uint8(clampFloat(tintG*strength) * 255),
		B: uint8(clampFloat(tintB*strength) * 255),
		A: uint8(clampFloat(strength*tintA) * 255),
	}
}

// SetNoiseAlgorithm allows swapping the underlying noise source.
func (f *Fog) SetNoiseAlgorithm(algo NoiseAlgorithm) {
	f.algo = algo
}

// SetFloatCenter selects the normalized center for the falloff.
func (f *Fog) SetFloatCenter(x, y float64) {
	f.FloatCenter.CenterX = x
	f.FloatCenter.CenterY = y
	f.useFloatCenter = true
}

func (f *Fog) radialWeight(x, y int) float64 {
	b := f.Bounds()
	if b.Empty() {
		return 0
	}

	cx := float64(b.Min.X + b.Dx()/2) // Default to center
	cy := float64(b.Min.Y + b.Dy()/2) // Default to center
	if f.useFloatCenter {
		cx = float64(b.Min.X) + f.FloatCenter.CenterX*float64(b.Dx())
		cy = float64(b.Min.Y) + f.FloatCenter.CenterY*float64(b.Dy())
	}

	dx := float64(x) - cx
	dy := float64(y) - cy

	dist := math.Sqrt(dx*dx + dy*dy)
	maxDist := math.Sqrt(float64(b.Dx()*b.Dx()+b.Dy()*b.Dy())) / 2.0
	if maxDist == 0 {
		return 0
	}

	falloff := f.FalloffCurve.FalloffCurve
	if falloff <= 0 {
		falloff = 1.0
	}

	// Keep a faint base layer (20%) so the center isn't completely empty.
	norm := clampFloat(dist / maxDist)
	t := clampFloat(math.Pow(norm, falloff))
	return clampFloat(0.2 + 0.8*t)
}

// NewFog creates a fog pattern using Perlin fBm and radial attenuation.
func NewFog(ops ...func(any)) image.Image {
	f := &Fog{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		FillColor: FillColor{
			FillColor: color.RGBA{200, 210, 225, 255},
		},
		Density: Density{
			Density: 0.8,
		},
		FalloffCurve: FalloffCurve{
			FalloffCurve: 1.6,
		},
		FloatCenter: FloatCenter{
			CenterX: 0.5,
			CenterY: 0.5,
		},
		algo: &PerlinNoise{
			Seed:        2024,
			Octaves:     5,
			Frequency:   0.012,
			Persistence: 0.55,
			Lacunarity:  2.1,
		},
	}

	for _, op := range ops {
		op(f)
	}

	if f.FalloffCurve.FalloffCurve <= 0 {
		f.FalloffCurve.FalloffCurve = 1.0
	}

	return f
}
