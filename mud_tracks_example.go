package pattern

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var MudTracksOutputFilename = "mud_tracks.png"

const MudTracksBaseLabel = "MudTracks"

// MudTracksOption configures the appearance of the mud tracks example/generator.
type MudTracksOption func(*mudTracksConfig)

type mudTracksConfig struct {
	bandWidth   int
	trackWobble float64
	mudDarkness float64
	bounds      image.Rectangle
}

// SetMudTracksBandWidth sets the width of each compacted band.
func SetMudTracksBandWidth(width int) MudTracksOption {
	return func(cfg *mudTracksConfig) {
		cfg.bandWidth = width
	}
}

// SetMudTracksWobble controls how much the bands meander vertically.
func SetMudTracksWobble(amount float64) MudTracksOption {
	return func(cfg *mudTracksConfig) {
		cfg.trackWobble = amount
	}
}

// SetMudTracksDarkness darkens the compacted mud and embedded pebbles (0 = light, 1 = dark).
func SetMudTracksDarkness(value float64) MudTracksOption {
	return func(cfg *mudTracksConfig) {
		cfg.mudDarkness = value
	}
}

// ExampleNewMudTracks lays down multiple compacted bands with embedded pebble noise.
func ExampleNewMudTracks() {
	img := buildMudTracks(nil)

	f, err := os.Create(MudTracksOutputFilename)
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

// GenerateMudTracks builds the pattern for registry-driven generation.
func GenerateMudTracks(b image.Rectangle) image.Image {
	cfg := defaultMudTracksConfig()
	cfg.bounds = b
	return buildMudTracks(cfg)
}

func defaultMudTracksConfig(opts ...MudTracksOption) *mudTracksConfig {
	cfg := &mudTracksConfig{
		bandWidth:   36,
		trackWobble: 12.0,
		mudDarkness: 0.35,
		bounds:      image.Rect(0, 0, 255, 255),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func buildMudTracks(cfg *mudTracksConfig) image.Image {
	if cfg == nil {
		cfg = defaultMudTracksConfig()
	}

	if cfg.bandWidth < 2 {
		cfg.bandWidth = 2
	}
	if cfg.trackWobble < 0 {
		cfg.trackWobble = 0
	}
	cfg.mudDarkness = clamp01(cfg.mudDarkness)

	baseHeight := NewNoise(
		SetBounds(cfg.bounds),
		NoiseSeed(2045),
		SetNoiseAlgorithm(&PerlinNoise{Seed: 2045, Frequency: 0.02, Octaves: 4, Persistence: 0.55}),
	)

	baseMud := NewColorMap(baseHeight, mudPalette(cfg.mudDarkness)...)

	compactedHeight := &yCompactedImage{
		Image: NewNoise(
			SetBounds(cfg.bounds),
			NoiseSeed(9090),
			SetNoiseAlgorithm(&PerlinNoise{Seed: 9090, Frequency: 0.025, Octaves: 5, Persistence: 0.6}),
		),
		factor: 0.6,
	}

	trackMud := NewColorMap(compactedHeight, trackPalette(cfg.mudDarkness)...)

	pebbleNoise := NewNoise(
		SetBounds(cfg.bounds),
		NoiseSeed(7777),
		SetNoiseAlgorithm(&PerlinNoise{Seed: 7777, Frequency: 0.12, Octaves: 3, Persistence: 0.5}),
	)
	pebbleLayer := NewColorMap(pebbleNoise,
		ColorStop{Position: 0.0, Color: color.RGBA{0, 0, 0, 0}},
		ColorStop{Position: 0.55, Color: color.RGBA{0, 0, 0, 0}},
		ColorStop{Position: 0.7, Color: color.RGBA{105, 95, 85, 90}},
		ColorStop{Position: 1.0, Color: color.RGBA{150, 140, 125, 160}},
	)

	trackMask := NewHorizontalLine(
		SetBounds(cfg.bounds),
		SetLineSize(cfg.bandWidth),
		SetSpaceSize(cfg.bandWidth+cfg.bandWidth/2),
		SetLineColor(color.White),
		SetSpaceColor(color.Transparent),
	)

	wobbleNoise := NewNoise(
		SetBounds(cfg.bounds),
		NoiseSeed(9911),
		SetNoiseAlgorithm(&PerlinNoise{Seed: 9911, Frequency: 0.01, Octaves: 2, Persistence: 0.6}),
	)

	wobbledMask := NewWarp(
		trackMask,
		WarpDistortionY(wobbleNoise),
		WarpScale(cfg.trackWobble),
		WarpDistortionScale(0.85),
	)

	trackLayer := NewBitwiseAnd([]image.Image{trackMud, wobbledMask})
	trackPebbles := NewBitwiseAnd([]image.Image{pebbleLayer, wobbledMask})
	compactedWithPebbles := NewBlend(trackLayer, trackPebbles, BlendOverlay)

	return NewBlend(baseMud, compactedWithPebbles, BlendNormal)
}

func mudPalette(darkness float64) []ColorStop {
	factor := 1.0 - 0.35*darkness
	return []ColorStop{
		{Position: 0.0, Color: darken(color.RGBA{70, 50, 35, 255}, factor)},
		{Position: 0.45, Color: darken(color.RGBA{120, 95, 70, 255}, factor)},
		{Position: 1.0, Color: darken(color.RGBA{165, 140, 105, 255}, factor)},
	}
}

func trackPalette(darkness float64) []ColorStop {
	factor := 0.85 * (1.0 - 0.35*darkness)
	return []ColorStop{
		{Position: 0.0, Color: darken(color.RGBA{55, 40, 28, 255}, factor)},
		{Position: 0.45, Color: darken(color.RGBA{100, 78, 58, 255}, factor)},
		{Position: 1.0, Color: darken(color.RGBA{140, 115, 86, 255}, factor)},
	}
}

type yCompactedImage struct {
	image.Image
	factor float64
}

func (c *yCompactedImage) At(x, y int) color.Color {
	b := c.Bounds()
	yy := int(math.Round((float64(y-b.Min.Y) * c.factor))) + b.Min.Y
	if yy < b.Min.Y {
		yy = b.Min.Y
	}
	if yy >= b.Max.Y {
		yy = b.Max.Y - 1
	}
	return c.Image.At(x, yy)
}

func darken(c color.RGBA, factor float64) color.RGBA {
	clamp := func(v float64) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v)
	}

	return color.RGBA{
		R: clamp(float64(c.R) * factor),
		G: clamp(float64(c.G) * factor),
		B: clamp(float64(c.B) * factor),
		A: c.A,
	}
}

func clamp01(v float64) float64 {
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return v
	}
}

func init() {
	RegisterGenerator(MudTracksBaseLabel, GenerateMudTracks)
}
