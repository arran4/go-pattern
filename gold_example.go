package pattern

import (
	"image"
	"image/color"
)

var (
	GoldOutputFilename = "gold.png"
)

func ExampleNewGold() image.Image {
	// 1. High frequency noise for grain
	noise := NewNoise(
		NoiseSeed(777),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        777,
			Frequency:   0.1,
			Octaves:     3,
			Persistence: 0.5,
		}),
	)

	// 2. Anisotropic scaling for brushed look (horizontal grain)
	brushed := NewScale(noise, ScaleX(10.0), ScaleY(1.0))

	// 3. Map to gold gradient
	// Dark: Brown/Bronze
	// Mid: Gold
	// Light: Pale Gold
	goldTex := NewColorMap(brushed,
		ColorStop{Position: 0.0, Color: color.RGBA{100, 70, 20, 255}},
		ColorStop{Position: 0.5, Color: color.RGBA{210, 170, 50, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{255, 230, 150, 255}},
	)

	return goldTex
}

func GenerateGold(rect image.Rectangle) image.Image {
	return ExampleNewGold()
}

func init() {
	RegisterGenerator("Gold", GenerateGold)
}
