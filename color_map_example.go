package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var ColorMapOutputFilename = "colormap.png"
var ColorMapZoomLevels = []int{}

// Base label for the main example
const ColorMapBaseLabel = "Grass"

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
		"Fire": func(b image.Rectangle) image.Image {
			noise := NewNoise(
				SetBounds(b),
				NoiseSeed(400),
				SetNoiseAlgorithm(&PerlinNoise{
					Seed:        400,
					Octaves:     5,
					Persistence: 0.6,
					Lacunarity:  2.2,
					Frequency:   0.08,
				}),
			)
			return NewColorMap(noise,
				ColorStop{Position: 0.0, Color: color.RGBA{0, 0, 0, 255}},       // Black soot
				ColorStop{Position: 0.3, Color: color.RGBA{180, 20, 0, 255}},    // Deep Red
				ColorStop{Position: 0.6, Color: color.RGBA{255, 140, 0, 255}},   // Orange
				ColorStop{Position: 0.9, Color: color.RGBA{255, 255, 0, 255}},   // Yellow
				ColorStop{Position: 1.0, Color: color.RGBA{255, 255, 255, 255}}, // White hot
			)
		},
		"Water": func(b image.Rectangle) image.Image {
			noise := NewNoise(
				SetBounds(b),
				NoiseSeed(500),
				SetNoiseAlgorithm(&PerlinNoise{
					Seed:        500,
					Octaves:     2,
					Persistence: 0.5,
					Lacunarity:  2.0,
					Frequency:   0.03,
				}),
			)
			return NewColorMap(noise,
				ColorStop{Position: 0.0, Color: color.RGBA{0, 0, 50, 255}},      // Deep Blue
				ColorStop{Position: 0.5, Color: color.RGBA{0, 50, 150, 255}},    // Mid Blue
				ColorStop{Position: 0.8, Color: color.RGBA{0, 150, 200, 255}},   // Shallow
				ColorStop{Position: 1.0, Color: color.RGBA{200, 240, 255, 255}}, // Foam/Highlights
			)
		},
		"Rust": func(b image.Rectangle) image.Image {
			noise := NewNoise(
				SetBounds(b),
				NoiseSeed(600),
				SetNoiseAlgorithm(&PerlinNoise{
					Seed:        600,
					Octaves:     6,
					Persistence: 0.7,
					Lacunarity:  2.5,
					Frequency:   0.15, // High frequency for grainy look
				}),
			)
			return NewColorMap(noise,
				ColorStop{Position: 0.0, Color: color.RGBA{50, 20, 10, 255}},    // Dark corroded metal
				ColorStop{Position: 0.4, Color: color.RGBA{120, 40, 10, 255}},   // Deep Rust
				ColorStop{Position: 0.7, Color: color.RGBA{180, 90, 20, 255}},   // Bright Rust
				ColorStop{Position: 0.9, Color: color.RGBA{150, 150, 150, 255}}, // Exposed metal
				ColorStop{Position: 1.0, Color: color.RGBA{200, 200, 200, 255}}, // Highlights
			)
		},
		"Terrain": func(b image.Rectangle) image.Image {
			noise := NewNoise(
				SetBounds(b),
				NoiseSeed(700),
				SetNoiseAlgorithm(&PerlinNoise{
					Seed:        700,
					Octaves:     4,
					Persistence: 0.5,
					Lacunarity:  2.0,
					Frequency:   0.02,
				}),
			)
			return NewColorMap(noise,
				ColorStop{Position: 0.0, Color: color.RGBA{0, 0, 150, 255}},      // Water (Deep)
				ColorStop{Position: 0.45, Color: color.RGBA{0, 100, 200, 255}},   // Water (Shallow)
				ColorStop{Position: 0.50, Color: color.RGBA{240, 230, 140, 255}}, // Sand
				ColorStop{Position: 0.60, Color: color.RGBA{50, 150, 50, 255}},   // Grass
				ColorStop{Position: 0.80, Color: color.RGBA{100, 100, 100, 255}}, // Rock
				ColorStop{Position: 0.95, Color: color.RGBA{255, 255, 255, 255}}, // Snow
			)
		},
		"Zebra": func(b image.Rectangle) image.Image {
			// Sharp transitions for stripe-like pattern
			noise := NewNoise(
				SetBounds(b),
				NoiseSeed(800),
				SetNoiseAlgorithm(&PerlinNoise{
					Seed:      800,
					Frequency: 0.1,
					Octaves:   1,
				}),
			)
			return NewColorMap(noise,
				ColorStop{Position: 0.45, Color: color.Black},
				ColorStop{Position: 0.55, Color: color.White},
			)
		},
	}, []string{"Dirt", "Clouds", "Heatmap", "Fire", "Water", "Rust", "Terrain", "Zebra"}
}

func init() {
	RegisterGenerator("ColorMap", GenerateColorMap)
	RegisterReferences("ColorMap", GenerateColorMapReferences)
}
