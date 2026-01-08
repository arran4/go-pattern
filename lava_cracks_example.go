package pattern

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var LavaCracksOutputFilename = "lava_cracks.png"

const LavaCracksBaseLabel = "LavaCracks"

// LavaCracks renders cooled volcanic rock with emissive lava held inside inverted cracks.
// The look is controlled by three main parameters:
//   - CrackDensity: frequency of the Worley noise used for the crack network.
//   - LavaGlowRadius: softness of the lava falloff away from the crack core.
//   - RockRoughness: amplitude of the surface noise used to roughen the rock.
type LavaCracks struct {
	Null
	CrackDensity   float64
	LavaGlowRadius float64
	RockRoughness  float64

	crackNoise image.Image
	rockNoise  image.Image
}

// SetCrackDensity adjusts how dense the crack network is.
func SetCrackDensity(freq float64) func(any) {
	return func(i any) {
		if lc, ok := i.(*LavaCracks); ok {
			lc.CrackDensity = freq
			lc.rebuildSources()
		}
	}
}

// SetLavaGlowRadius changes how wide the lava glow spreads from the crack core.
func SetLavaGlowRadius(radius float64) func(any) {
	return func(i any) {
		if lc, ok := i.(*LavaCracks); ok {
			lc.LavaGlowRadius = radius
		}
	}
}

// SetRockRoughness tunes the micro variation of the cooled rock surface.
func SetRockRoughness(roughness float64) func(any) {
	return func(i any) {
		if lc, ok := i.(*LavaCracks); ok {
			lc.RockRoughness = roughness
			lc.rebuildSources()
		}
	}
}

func (lc *LavaCracks) SetBounds(bounds image.Rectangle) {
	lc.Null.SetBounds(bounds)
	lc.rebuildSources()
}

// At blends emissive lava in the crack cores with a cooled, rough rock plate.
func (lc *LavaCracks) At(x, y int) color.Color {
	crackGray := color.GrayModel.Convert(lc.crackNoise.At(x, y)).(color.Gray)
	crackEdge := float64(crackGray.Y) / 255.0

	// Invert the crack mask to hold lava inside the fissures.
	lavaMask := clamp01(1.0 - crackEdge)

	// Soften lava based on the requested glow radius.
	glowTightness := 0.75 + math.Max(lc.LavaGlowRadius, 0.1)*0.2
	glow := math.Pow(lavaMask, 1.0/glowTightness)
	core := math.Pow(glow, 0.65)

	rockGray := color.GrayModel.Convert(lc.rockNoise.At(x, y)).(color.Gray)
	rockSample := (float64(rockGray.Y)/255.0 - 0.5) * 2.0

	roughScale := clamp01(lc.RockRoughness)
	roughMix := clamp01(0.5 + rockSample*roughScale*0.6)

	rockDark := colorVec{R: 40, G: 38, B: 46}
	rockLight := colorVec{R: 86, G: 82, B: 92}
	rockColor := lerpVec(rockDark, rockLight, roughMix)

	// Cool the edges where lava has chilled.
	coolEdge := clamp01(crackEdge * 0.8)
	coolTint := colorVec{R: 34, G: 48, B: 64}
	rockColor = lerpVec(rockColor, coolTint, coolEdge*0.6)

	lavaDeep := colorVec{R: 120, G: 24, B: 12}
	lavaHot := colorVec{R: 255, G: 168, B: 70}
	lavaCore := colorVec{R: 255, G: 236, B: 200}

	lavaColor := lerpVec(lavaDeep, lavaHot, glow)
	lavaColor = lerpVec(lavaColor, lavaCore, core)

	blend := clamp01(glow * 1.05)
	finalColor := lerpVec(rockColor, lavaColor, blend)

	return color.RGBA{
		R: clampToUint8(finalColor.R),
		G: clampToUint8(finalColor.G),
		B: clampToUint8(finalColor.B),
		A: 255,
	}
}

// NewLavaCracks creates the lava crack pattern with optional parameter overrides.
func NewLavaCracks(ops ...func(any)) image.Image {
	lc := &LavaCracks{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		CrackDensity:   0.035,
		LavaGlowRadius: 6.0,
		RockRoughness:  0.6,
	}

	for _, op := range ops {
		op(lc)
	}

	lc.rebuildSources()
	return lc
}

func (lc *LavaCracks) rebuildSources() {
	freq := lc.CrackDensity
	if freq <= 0 {
		freq = 0.035
	}

	lc.crackNoise = NewWorleyNoise(
		SetBounds(lc.bounds),
		SetFrequency(freq),
		SetWorleyOutput(OutputF2MinusF1),
		SetWorleyMetric(MetricEuclidean),
		WithSeed(777),
	)

	roughness := lc.RockRoughness
	if roughness <= 0 {
		roughness = 0.6
	}
	rockFreq := 0.02 + 0.08*roughness

	lc.rockNoise = NewNoise(
		SetBounds(lc.bounds),
		SetNoiseAlgorithm(&PerlinNoise{
			Frequency:   rockFreq,
			Octaves:     4,
			Persistence: 0.55,
		}),
		NoiseSeed(4242),
	)
}

// GenerateLavaCracks renders the pattern into the provided bounds using default styling.
func GenerateLavaCracks(b image.Rectangle) image.Image {
	return NewLavaCracks(
		SetBounds(b),
		SetCrackDensity(0.04),
		SetLavaGlowRadius(7.0),
		SetRockRoughness(0.7),
	)
}

// ExampleNewLavaCracks produces an example PNG demonstrating the lava cracks material.
func ExampleNewLavaCracks() {
	img := GenerateLavaCracks(image.Rect(0, 0, 200, 200))
	f, err := os.Create(LavaCracksOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}

func init() {
	RegisterGenerator(LavaCracksBaseLabel, GenerateLavaCracks)
	RegisterReferences(LavaCracksBaseLabel, func() (map[string]func(image.Rectangle) image.Image, []string) {
		return map[string]func(image.Rectangle) image.Image{}, []string{}
	})
}

type colorVec struct {
	R, G, B float64
}

func lerpVec(a, b colorVec, t float64) colorVec {
	t = clamp01(t)
	return colorVec{
		R: a.R + (b.R-a.R)*t,
		G: a.G + (b.G-a.G)*t,
		B: a.B + (b.B-a.B)*t,
	}
}

func clampToUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(math.Round(v))
}
