package pattern

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// FrostCracks renders icy fissures by thresholding the gradient magnitude of an fBm height field.
// The density controls how readily cracks appear, the glow boosts light within the fissures,
// and blur softens the edges to avoid aliasing.
type FrostCracks struct {
	Null
	Seed
	Density
	GlowAmount  float64
	BlurAmount  int
	Frequency   float64
	Octaves     int
	Persistence float64
	Lacunarity  float64

	once   sync.Once
	pixels []color.RGBA
}

func (f *FrostCracks) SetGlowAmount(v float64) {
	f.GlowAmount = v
}

func (f *FrostCracks) SetBlurAmount(v int) {
	f.BlurAmount = v
}

type hasGlowAmount interface {
	SetGlowAmount(float64)
}

type hasBlurAmount interface {
	SetBlurAmount(int)
}

// SetGlowAmount configures how strongly light blooms inside the cracks.
func SetGlowAmount(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasGlowAmount); ok {
			h.SetGlowAmount(v)
		}
	}
}

// SetBlurAmount configures how many 3x3 blur passes are applied to the gradient mask.
func SetBlurAmount(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasBlurAmount); ok {
			h.SetBlurAmount(v)
		}
	}
}

// NewFrostCracks builds a frozen crack pattern based on fractal Brownian motion.
func NewFrostCracks(ops ...func(any)) image.Image {
	p := &FrostCracks{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Seed:       Seed{Seed: 1337},
		Density:    Density{Density: 0.6},
		GlowAmount: 0.5,
		BlurAmount: 2,
		// Sensible defaults for fBm
		Frequency:   0.012,
		Octaves:     5,
		Persistence: 0.55,
		Lacunarity:  2.1,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

func (f *FrostCracks) ColorModel() color.Model {
	return color.RGBAModel
}

func (f *FrostCracks) Bounds() image.Rectangle {
	return f.bounds
}

func (f *FrostCracks) At(x, y int) color.Color {
	f.generate()
	if !image.Pt(x, y).In(f.bounds) {
		return color.RGBA{}
	}

	w := f.bounds.Dx()
	gx := x - f.bounds.Min.X
	gy := y - f.bounds.Min.Y
	return f.pixels[gy*w+gx]
}

func (f *FrostCracks) generate() {
	f.once.Do(func() {
		w := f.bounds.Dx()
		h := f.bounds.Dy()

		density := clamp01(f.Density.Density)
		glow := clamp01(f.GlowAmount)
		blurPasses := f.BlurAmount
		if blurPasses < 0 {
			blurPasses = 0
		}

		freq := f.Frequency
		if freq <= 0 {
			freq = 0.008 + 0.02*density
		}
		octaves := f.Octaves
		if octaves <= 0 {
			octaves = 5
		}
		persistence := f.Persistence
		if persistence <= 0 {
			persistence = 0.55
		}
		lacunarity := f.Lacunarity
		if lacunarity <= 0 {
			lacunarity = 2.1
		}

		height := make([]float64, w*h)
		perlin := &PerlinNoise{
			Seed:        f.Seed.Seed,
			Octaves:     octaves,
			Persistence: persistence,
			Lacunarity:  lacunarity,
			Frequency:   freq,
		}

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := y*w + x
				c := perlin.At(f.bounds.Min.X+x, f.bounds.Min.Y+y)
				r, _, _, _ := c.RGBA()
				height[idx] = float64(r) / 65535.0
			}
		}

		grad := make([]float64, w*h)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				// Sobel on the scalar height field
				v00 := sampleClamp(height, w, h, x-1, y-1)
				v10 := sampleClamp(height, w, h, x, y-1)
				v20 := sampleClamp(height, w, h, x+1, y-1)
				v01 := sampleClamp(height, w, h, x-1, y)
				v21 := sampleClamp(height, w, h, x+1, y)
				v02 := sampleClamp(height, w, h, x-1, y+1)
				v12 := sampleClamp(height, w, h, x, y+1)
				v22 := sampleClamp(height, w, h, x+1, y+1)

				gx := -1*v00 + 1*v20 + -2*v01 + 2*v21 + -1*v02 + 1*v22
				gy := -1*v00 + -2*v10 + -1*v20 + 1*v02 + 2*v12 + 1*v22

				mag := math.Sqrt(gx*gx + gy*gy)
				grad[y*w+x] = mag
			}
		}

		for i := 0; i < blurPasses; i++ {
			grad = blur3x3(grad, w, h)
		}

		maxGrad := 0.0
		for _, g := range grad {
			if g > maxGrad {
				maxGrad = g
			}
		}
		if maxGrad == 0 {
			maxGrad = 1
		}

		threshold := 0.72 - 0.42*density
		baseLow := color.RGBA{8, 18, 35, 255}
		baseHigh := color.RGBA{50, 110, 160, 255}
		crackEdge := color.RGBA{90, 150, 210, 255}
		crackCore := color.RGBA{215, 240, 255, 255}

		f.pixels = make([]color.RGBA, w*h)
		for i, g := range grad {
			base := clamp01(height[i]*0.9 + 0.05)
			edgeStrength := clamp01((g/maxGrad - threshold) / (1 - threshold))
			crack := math.Pow(edgeStrength, 0.75)
			glowStrength := math.Pow(edgeStrength, 1.2) * (0.35 + 0.65*glow)

			cool := lerpRGBA(baseLow, baseHigh, base)
			ice := lerpRGBA(crackEdge, crackCore, 0.35+0.4*glow)

			col := blendRGBA(cool, ice, crack)
			col = blendRGBA(col, crackCore, clamp01(glowStrength))
			f.pixels[i] = col
		}
	})
}

func lerpRGBA(a, b color.RGBA, t float64) color.RGBA {
	t = clamp01(t)
	return color.RGBA{
		R: uint8(math.Round(float64(a.R) + (float64(b.R)-float64(a.R))*t)),
		G: uint8(math.Round(float64(a.G) + (float64(b.G)-float64(a.G))*t)),
		B: uint8(math.Round(float64(a.B) + (float64(b.B)-float64(a.B))*t)),
		A: 255,
	}
}

func blendRGBA(a, b color.RGBA, t float64) color.RGBA {
	return lerpRGBA(a, b, clamp01(t))
}

func sampleClamp(buf []float64, w, h, x, y int) float64 {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= w {
		x = w - 1
	}
	if y >= h {
		y = h - 1
	}
	return buf[y*w+x]
}

func blur3x3(src []float64, w, h int) []float64 {
	dst := make([]float64, len(src))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sum := 0.0
			count := 0
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					sum += sampleClamp(src, w, h, x+kx, y+ky)
					count++
				}
			}
			dst[y*w+x] = sum / float64(count)
		}
	}
	return dst
}
