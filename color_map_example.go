package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var ColorMapOutputFilename = "colormap.png"
var ColorMapZoomLevels = []int{}

const ColorMapOrder = 25

// ColorMap Pattern
// Maps the luminance of a source pattern to a color gradient (ramp).
// This is useful for creating textures like grass, dirt, clouds, or heatmaps.
func ExampleNewColorMap() {
	// 1. Create a Noise source (Perlin Noise with FBM)
	noise := NewNoise(
		NoiseSeed(42), // Fixed seed for reproducible documentation
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        42,
			Octaves:     4,
			Persistence: 0.5,
			Lacunarity:  2.0,
			Frequency:   0.1,
		}),
	)

	// 2. Map the noise to a "Grass" color ramp
	grass := NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{0, 50, 0, 255}},     // Deep shadow green
		ColorStop{Position: 0.4, Color: color.RGBA{10, 100, 10, 255}},  // Mid green
		ColorStop{Position: 0.7, Color: color.RGBA{50, 150, 30, 255}},  // Light green
		ColorStop{Position: 1.0, Color: color.RGBA{100, 140, 60, 255}}, // Dried tip
	)

	f, err := os.Create(ColorMapOutputFilename)
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

func GenerateColorMap(b image.Rectangle) image.Image {
	noise := NewNoise(
		SetBounds(b),
		NoiseSeed(42),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        42,
			Octaves:     4,
			Persistence: 0.5,
			Lacunarity:  2.0,
			Frequency:   0.1,
		}),
	)
	return NewColorMap(noise,
		ColorStop{Position: 0.0, Color: color.RGBA{0, 50, 0, 255}},
		ColorStop{Position: 0.4, Color: color.RGBA{10, 100, 10, 255}},
		ColorStop{Position: 0.7, Color: color.RGBA{50, 150, 30, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{100, 140, 60, 255}},
	)
}

func GenerateColorMapReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Dirt": func(b image.Rectangle) image.Image {
			noise := NewNoise(
				SetBounds(b),
				NoiseSeed(100),
				SetNoiseAlgorithm(&PerlinNoise{
					Seed:        100,
					Octaves:     6,
					Persistence: 0.6,
					Lacunarity:  2.0,
					Frequency:   0.05,
				}),
			)
			return NewColorMap(noise,
				ColorStop{Position: 0.0, Color: color.RGBA{40, 30, 20, 255}},   // Dark Brown
				ColorStop{Position: 0.5, Color: color.RGBA{80, 60, 40, 255}},   // Mid Brown
				ColorStop{Position: 0.8, Color: color.RGBA{100, 80, 60, 255}},  // Light Brown
				ColorStop{Position: 1.0, Color: color.RGBA{120, 100, 80, 255}}, // Dusty
			)
		},
		"Clouds": func(b image.Rectangle) image.Image {
			noise := NewNoise(
				SetBounds(b),
				NoiseSeed(200),
				SetNoiseAlgorithm(&PerlinNoise{
					Seed:        200,
					Octaves:     3,
					Persistence: 0.5,
					Lacunarity:  2.0,
					Frequency:   0.015,
				}),
			)
			return NewColorMap(noise,
				ColorStop{Position: 0.0, Color: color.RGBA{135, 206, 235, 255}}, // Sky Blue
				ColorStop{Position: 0.5, Color: color.RGBA{135, 206, 235, 255}}, // Sky Blue
				ColorStop{Position: 0.7, Color: color.RGBA{240, 240, 255, 255}}, // Whiteish
				ColorStop{Position: 1.0, Color: color.RGBA{255, 255, 255, 255}}, // White
			)
		},
		"Heatmap": func(b image.Rectangle) image.Image {
			// Using standard gradient noise for smooth transition
			noise := NewNoise(
				SetBounds(b),
				NoiseSeed(300),
				SetNoiseAlgorithm(&PerlinNoise{
					Seed:      300,
					Frequency: 0.02,
				}),
			)
			return NewColorMap(noise,
				ColorStop{Position: 0.0, Color: color.RGBA{0, 0, 255, 255}},   // Blue
				ColorStop{Position: 0.4, Color: color.RGBA{0, 255, 255, 255}}, // Cyan
				ColorStop{Position: 0.6, Color: color.RGBA{0, 255, 0, 255}},   // Green
				ColorStop{Position: 0.8, Color: color.RGBA{255, 255, 0, 255}}, // Yellow
				ColorStop{Position: 1.0, Color: color.RGBA{255, 0, 0, 255}},   // Red
			)
		},
	}, []string{"Dirt", "Clouds", "Heatmap"}
}

func init() {
	RegisterGenerator("ColorMap", GenerateColorMap)
	RegisterReferences("ColorMap", GenerateColorMapReferences)
}
