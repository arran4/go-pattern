package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var DamascusOutputFilename = "damascus.png"

const DamascusBaseLabel = "Damascus"

func ExampleNewDamascus() {
	img := GenerateDamascus(image.Rect(0, 0, 150, 150))
	f, err := os.Create(DamascusOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}

func GenerateDamascus(b image.Rectangle) image.Image {
	// Damascus Steel Effect
	// High frequency ridge noise distorted by low freq noise

	// 1. Base Pattern: Sine waves or very regular noise
	// Using simple Perlin with high frequency
	base := NewNoise(
		SetBounds(b),
		NoiseSeed(111),
		SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.1, // High freq for bands
			Octaves:   1,
		}),
	)

	// 2. Strong Distortion
	dist := NewNoise(
		SetBounds(b),
		NoiseSeed(222),
		SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.015,
			Octaves:   2,
		}),
	)

	warped := NewWarp(base,
		WarpDistortion(dist),
		WarpScale(60.0), // Strong warping
	)

	// 3. Map to Metal bands
	// We want repeating bands of dark/light grey

	dark := color.RGBA{60, 60, 65, 255}
	light := color.RGBA{180, 180, 190, 255}

	// Manually set alternating bands
	stops := []ColorStop{
		{0.0, dark}, {0.1, light}, {0.2, dark}, {0.3, light},
		{0.4, dark}, {0.5, light}, {0.6, dark}, {0.7, light},
		{0.8, dark}, {0.9, light}, {1.0, dark},
	}

	return NewColorMap(warped, stops...)
}

func init() {
	RegisterGenerator(DamascusBaseLabel, GenerateDamascus)
	RegisterReferences(DamascusBaseLabel, func() (map[string]func(image.Rectangle) image.Image, []string) {
		return map[string]func(image.Rectangle) image.Image{}, []string{}
	})
}
