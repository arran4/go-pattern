package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var CamouflageOutputFilename = "camouflage.png"

const CamouflageBaseLabel = "Camouflage"

func ExampleNewCamouflage() {
	img := GenerateCamouflage(image.Rect(0, 0, 150, 150))
	f, err := os.Create(CamouflageOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}

func GenerateCamouflage(b image.Rectangle) image.Image {
	// Standard Woodland Camo
	// Large blotches of Green, Brown, Black, Tan.

	// 1. Base Noise: Low Frequency Perlin
	noise := NewNoise(
		SetBounds(b),
		NoiseSeed(333),
		SetNoiseAlgorithm(&PerlinNoise{
			Frequency:   0.02,
			Octaves:     4,
			Persistence: 0.7, // Rough edges
		}),
	)

	// 2. Quantize: We want hard edges between colors.
	// Map grayscale to 4 discrete bands.

	// Define Colors
	black := color.RGBA{20, 20, 20, 255}
	brown := color.RGBA{101, 67, 33, 255}
	green := color.RGBA{85, 107, 47, 255} // Olive Drab
	tan := color.RGBA{210, 180, 140, 255}

	return NewColorMap(noise,
		ColorStop{0.0, black},
		ColorStop{0.3, black},
		ColorStop{0.31, brown},
		ColorStop{0.5, brown},
		ColorStop{0.51, green},
		ColorStop{0.7, green},
		ColorStop{0.71, tan},
		ColorStop{1.0, tan},
	)
}

func init() {
	RegisterGenerator(CamouflageBaseLabel, GenerateCamouflage)
	RegisterReferences(CamouflageBaseLabel, func() (map[string]func(image.Rectangle) image.Image, []string) {
		return map[string]func(image.Rectangle) image.Image{}, []string{}
	})
}
