package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var GrassCloseOutputFilename = "grass_close.png"
var GrassCloseZoomLevels = []int{}

const GrassCloseBaseLabel = "Grass Close"

// Grass Close Example
// Demonstrates a procedural grass texture using the GrassClose pattern composed with Noise.
func ExampleNewGrassClose() {
	// 1. Background: Dirt
	dirt := NewColorMap(
		NewNoise(SetFrequency(0.05), NoiseSeed(1)),
		ColorStop{0.0, color.RGBA{40, 30, 20, 255}},
		ColorStop{1.0, color.RGBA{80, 60, 40, 255}},
	)

	// 2. Wind map (Perlin noise)
	wind := NewNoise(
		SetFrequency(0.01),
		NoiseSeed(2),
		SetNoiseAlgorithm(&PerlinNoise{Seed: 2, Octaves: 2, Persistence: 0.5}),
	)

	// 3. Density map (Worley noise for clumping)
	density := NewWorleyNoise(
		SetFrequency(0.02),
		SetSeed(3),
	)

	// 4. Grass Layer
	grass := NewGrassClose(
		SetBladeHeight(35),
		SetBladeWidth(5),
		SetFillColor(color.RGBA{20, 160, 30, 255}),
		SetWindSource(wind),
		SetDensitySource(density),
		// Background source
		func(p any) {
			if g, ok := p.(*GrassClose); ok {
				g.Source = dirt
			}
		},
	)

	f, err := os.Create(GrassCloseOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, grass); err != nil {
		panic(err)
	}
}

func GenerateGrassClose(b image.Rectangle) image.Image {
	dirt := NewColorMap(
		NewNoise(SetBounds(b), SetFrequency(0.05), NoiseSeed(1)),
		ColorStop{0.0, color.RGBA{40, 30, 20, 255}},
		ColorStop{1.0, color.RGBA{80, 60, 40, 255}},
	)
	wind := NewNoise(
		SetBounds(b),
		SetFrequency(0.01),
		NoiseSeed(2),
	)
	density := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.02),
		SetSeed(3),
	)

	// Helper to set source since it's not a standard option yet
	setSource := func(src image.Image) func(any) {
		return func(p any) {
			if g, ok := p.(*GrassClose); ok {
				g.Source = src
			}
		}
	}

	return NewGrassClose(
		SetBounds(b),
		SetBladeHeight(30),
		SetBladeWidth(4),
		SetFillColor(color.RGBA{20, 160, 30, 255}),
		SetWindSource(wind),
		SetDensitySource(density),
		setSource(dirt),
	)
}

func GenerateGrassCloseReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Lush": func(b image.Rectangle) image.Image {
			// Deep green background
			bg := NewColorMap(
				NewNoise(SetBounds(b), SetFrequency(0.05), NoiseSeed(10)),
				ColorStop{0.0, color.RGBA{0, 40, 0, 255}},
				ColorStop{1.0, color.RGBA{0, 80, 20, 255}},
			)
			// Layer 1: Dark, short undergrowth
			layer1 := NewGrassClose(
				SetBounds(b),
				SetBladeHeight(15),
				SetBladeWidth(3),
				SetFillColor(color.RGBA{0, 100, 0, 255}),
				func(p any) { p.(*GrassClose).Source = bg }, // Source
				func(p any) { p.(*GrassClose).Seed = 11 },
			)
			// Layer 2: Longer, lighter grass
			layer2 := NewGrassClose(
				SetBounds(b),
				SetBladeHeight(30),
				SetBladeWidth(4),
				SetFillColor(color.RGBA{50, 180, 40, 255}),
				func(p any) { p.(*GrassClose).Source = layer1 },
				func(p any) { p.(*GrassClose).Seed = 12 },
				SetWindSource(NewNoise(SetBounds(b), SetFrequency(0.02), NoiseSeed(13))),
			)
			return layer2
		},
		"Dry": func(b image.Rectangle) image.Image {
			bg := NewColorMap(
				NewNoise(SetBounds(b), SetFrequency(0.05), NoiseSeed(20)),
				ColorStop{0.0, color.RGBA{60, 50, 30, 255}},
				ColorStop{1.0, color.RGBA{100, 90, 60, 255}},
			)
			return NewGrassClose(
				SetBounds(b),
				SetBladeHeight(25),
				SetBladeWidth(3),
				SetFillColor(color.RGBA{160, 140, 80, 255}), // Yellowish
				func(p any) { p.(*GrassClose).Source = bg },
				SetWindSource(NewNoise(SetBounds(b), SetFrequency(0.02), NoiseSeed(21))),
			)
		},
		"Distant Clumps": func(b image.Rectangle) image.Image {
			// Ground
			bg := NewColorMap(
				NewNoise(SetBounds(b), SetFrequency(0.1), NoiseSeed(30)),
				ColorStop{0.0, color.RGBA{50, 40, 30, 255}},
				ColorStop{1.0, color.RGBA{70, 60, 50, 255}},
			)
			// Density map: Worley Noise clumps
			// High contrast to make distinct patches
			density := NewColorMap(
				NewWorleyNoise(SetBounds(b), SetFrequency(0.05), SetSeed(31), SetWorleyOutput(OutputF1)),
				ColorStop{0.0, color.White},                    // Center of cell -> Dense grass
				ColorStop{0.3, color.RGBA{128, 128, 128, 255}}, // Edge of clump -> Sparse
				ColorStop{0.4, color.Black},                    // Outside -> No grass
			)

			return NewGrassClose(
				SetBounds(b),
				SetBladeHeight(8), // Small blades
				SetBladeWidth(2),
				SetFillColor(color.RGBA{30, 100, 30, 255}),
				func(p any) { p.(*GrassClose).Source = bg },
				SetDensitySource(density),
				SetWindSource(NewNoise(SetBounds(b), SetFrequency(0.01), NoiseSeed(32))),
				func(p any) { p.(*GrassClose).Seed = 33 },
			)
		},
	}, []string{"Lush", "Dry", "Distant Clumps"}
}

func init() {
	RegisterGenerator("GrassClose", GenerateGrassClose)
	RegisterReferences("GrassClose", GenerateGrassCloseReferences)
}
