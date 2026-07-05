package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var LavaOutputFilename = "lava.png"

const LavaBaseLabel = "Lava"

func ExampleNewLava() {
	// This function is for the testable example and documentation.
	// It creates the file directly.
	img := GenerateLava(image.Rect(0, 0, 150, 150))
	f, err := os.Create(LavaOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}

func GenerateLava(b image.Rectangle) image.Image {
	// 1. Base Turbulence: Perlin Noise
	baseNoise := NewNoise(
		SetBounds(b),
		NoiseSeed(1234),
		SetNoiseAlgorithm(&PerlinNoise{
			Frequency:   0.03,
			Octaves:     4,
			Persistence: 0.6,
		}),
	)

	// 2. Distortion: Another Perlin Noise
	distortion := NewNoise(
		SetBounds(b),
		NoiseSeed(5678),
		SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.02,
			Octaves:   2,
		}),
	)

	// 3. Warp the base noise
	warped := NewWarp(baseNoise,
		WarpDistortion(distortion),
		WarpScale(30.0), // Large distortion for fluid look
	)

	// 4. Color Mapping: Heat Gradient
	// 0.0 -> Black/Dark Red (Crust)
	// 0.3 -> Red
	// 0.6 -> Orange
	// 0.8 -> Yellow
	// 1.0 -> White (Hot)
	return NewColorMap(warped,
		ColorStop{Position: 0.0, Color: color.RGBA{20, 0, 0, 255}},
		ColorStop{Position: 0.2, Color: color.RGBA{80, 0, 0, 255}},
		ColorStop{Position: 0.4, Color: color.RGBA{200, 20, 0, 255}},
		ColorStop{Position: 0.7, Color: color.RGBA{255, 140, 0, 255}},
		ColorStop{Position: 0.9, Color: color.RGBA{255, 255, 0, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{255, 255, 200, 255}},
	)
}

func init() {
	RegisterGenerator(LavaBaseLabel, GenerateLava)
	RegisterReferences(LavaBaseLabel, func() (map[string]func(image.Rectangle) image.Image, []string) {
		return map[string]func(image.Rectangle) image.Image{}, []string{}
	})
}
