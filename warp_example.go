package pattern

import (
	"fmt"
	"image"
	"image/color"
)

func ExampleNewWarp_wood() {
	// Base pattern: Concentric Rings (mimicking tree rings)
	// We want dense rings.
	// Since ConcentricRings relies on integer coords, we might need to zoom or scale to get fine rings if we just use alternating black/white.
	// Better: Use a Gradient over the rings? Or just many colors?
	// Let's use a palette of wood colors.

	woodLight := color.RGBA{222, 184, 135, 255} // Burlywood
	woodDark := color.RGBA{139, 69, 19, 255}    // SaddleBrown

	// Create a gradient-like palette for the rings
	// Alternating bands of light and dark
	colors := []color.Color{}
	for i := 0; i < 10; i++ {
		// Soft gradient from light to dark
		colors = append(colors, woodLight)
		colors = append(colors, woodDark)
	}

	rings := NewConcentricRings(colors)

	// Distortion noise
	// We want the rings to be wobbly.
	noise := NewNoise(NoiseSeed(1), SetNoiseAlgorithm(&PerlinNoise{
		Frequency:   0.05,
		Octaves:     3,
		Persistence: 0.5,
		Lacunarity:  2.0,
		Seed:        1,
	}))

	// Apply Warp
	// We want significant distortion to look like organic growth
	wood := NewWarp(rings,
		WarpDistortion(noise),
		WarpScale(20.0), // Magnitude of wobble
	)

	// To make it look more like a plank, we might want to stretch the noise or the rings?
	// But let's start with a cross-section.

	fmt.Println(wood.At(10, 10))
	// Output: {222 184 135 255}
}

func ExampleNewWarp_marble() {
	// Marble: Linear gradient or Sine waves distorted by turbulence.
	// Let's use LinearGradient for veins.

	// Base: White background
	// We need veins.
	// Maybe a Checker pattern where one color is White and the other is Light Gray?
	// Or Sine waves? We don't have a "SineWave" pattern exposed directly except maybe via math ops or Stripe?
	// ModuloStripe might work.

	colors := []color.Color{
		color.RGBA{250, 250, 250, 255}, // White
		color.RGBA{200, 200, 200, 255}, // Gray
		color.RGBA{50, 50, 50, 255},    // Dark Vein
	}

	// Stripes
	stripes := NewModuloStripe(colors) // Defaults to vertical stripes

	// High frequency turbulence
	noiseX := NewNoise(NoiseSeed(2), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.02,
		Octaves: 5,
	}))
	noiseY := NewNoise(NoiseSeed(3), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.02,
		Octaves: 5,
	}))

	marble := NewWarp(stripes,
		WarpDistortionX(noiseX),
		WarpDistortionY(noiseY),
		WarpXScale(40.0),
		WarpYScale(40.0),
	)

	fmt.Println(marble.At(10, 10))
	// Output: {200 200 200 255}
}

const WarpBaseLabel = "Warp"

func init() {
	GlobalGenerators[WarpBaseLabel] = GenerateWarp
	GlobalReferences[WarpBaseLabel] = GenerateWarpReferences
}

func GenerateWarp(rect image.Rectangle) image.Image {
	// Standard demo: Grid warped by noise
	// CellPos requires 3 args: x, y, content.
	// But in NewGrid, CellPos is a helper function defined in grid.go: func CellPos(x, y int, content any) any

	// I used NewGrid(CellPos(image.Rect(0,0,50,50)), FixedSize(50,50)) which is wrong.
	// CellPos is for placing a cell at a position.

	// I want a grid of cells. But `checker` is simpler.

	checker := NewChecker(
		color.RGBA{200, 200, 200, 255},
		color.RGBA{50, 50, 50, 255},
	)

	noise := NewNoise(NoiseSeed(99), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.03,
		Octaves: 2,
	}))

	return NewWarp(checker,
		WarpDistortion(noise),
		WarpScale(10.0),
	)
}

// GenerateWarpReferences registers the examples for documentation generation.
func GenerateWarpReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	refs := make(map[string]func(image.Rectangle) image.Image)

	refs["Wood"] = func(rect image.Rectangle) image.Image {
		woodLight := color.RGBA{222, 184, 135, 255}
		woodDark := color.RGBA{139, 69, 19, 255}

		// Interpolate colors for smoother rings
		colors := []color.Color{}
		steps := 20
		for i := 0; i < steps; i++ {
			t := float64(i) / float64(steps-1)
			// Simple linear interpolation
			r := uint8(float64(woodLight.R)*(1-t) + float64(woodDark.R)*t)
			g := uint8(float64(woodLight.G)*(1-t) + float64(woodDark.G)*t)
			b := uint8(float64(woodLight.B)*(1-t) + float64(woodDark.B)*t)
			colors = append(colors, color.RGBA{r, g, b, 255})
		}
		// And back to light for smooth cycling
		for i := steps - 1; i >= 0; i-- {
			colors = append(colors, colors[i])
		}

		rings := NewConcentricRings(colors)

		// Strong, low frequency noise for shape
		noiseLow := NewNoise(NoiseSeed(123), SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.02,
			Octaves: 2,
		}))

		// Warping the rings
		return NewWarp(rings,
			WarpDistortion(noiseLow),
			WarpScale(15.0),
		)
	}

	refs["Marble"] = func(rect image.Rectangle) image.Image {
		// Marble logic
		colors := []color.Color{
			color.RGBA{240, 240, 245, 255}, // White-ish
			color.RGBA{240, 240, 245, 255},
			color.RGBA{240, 240, 245, 255},
			color.RGBA{200, 200, 210, 255}, // Light Gray
			color.RGBA{100, 100, 110, 255}, // Dark Vein
			color.RGBA{200, 200, 210, 255},
		}
		stripes := NewModuloStripe(colors)

		noise := NewNoise(NoiseSeed(456), SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.04,
			Octaves: 4,
			Persistence: 0.6,
		}))

		return NewWarp(stripes,
			WarpDistortion(noise),
			WarpScale(30.0), // Strong warping
		)
	}

	refs["Clouds"] = func(rect image.Rectangle) image.Image {
		// Warped Noise for Wispy Clouds

		baseNoise := NewNoise(NoiseSeed(777), SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.02,
			Octaves: 4,
			Persistence: 0.5,
		}))

		warpNoise := NewNoise(NoiseSeed(888), SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.02,
			Octaves: 2,
		}))

		warped := NewWarp(baseNoise,
			WarpDistortion(warpNoise),
			WarpScale(50.0),
		)

		// Map grayscale to Sky Blue -> White
		stops := []ColorStop{
			{0.0, color.RGBA{0, 100, 200, 255}}, // Blue Sky
			{0.4, color.RGBA{100, 150, 255, 255}}, // Lighter Blue
			{0.6, color.RGBA{255, 255, 255, 255}}, // Clouds
			{1.0, color.RGBA{255, 255, 255, 255}},
		}

		return NewColorMap(warped, stops...)
	}

	refs["Terrain"] = func(rect image.Rectangle) image.Image {
		// Domain Warped Terrain (Heightmap visualized)
		// FBM warped by FBM

		fbm := func(seed int64) image.Image {
			return NewNoise(NoiseSeed(seed), SetNoiseAlgorithm(&PerlinNoise{
				Frequency: 0.015,
				Octaves: 6,
				Persistence: 0.5,
				Lacunarity: 2.0,
			}))
		}

		base := fbm(101)

		// Warp with lower frequency
		warp := NewNoise(NoiseSeed(202), SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.01,
			Octaves: 2,
		}))

		warped := NewWarp(base,
			WarpDistortion(warp),
			WarpScale(80.0),
		)

		// Earth-like coloring
		stops := []ColorStop{
			{0.0, color.RGBA{0, 0, 150, 255}},     // Deep Ocean
			{0.2, color.RGBA{0, 50, 200, 255}},    // Ocean
			{0.22, color.RGBA{240, 230, 140, 255}}, // Sand
			{0.3, color.RGBA{34, 139, 34, 255}},   // Grass
			{0.6, color.RGBA{107, 142, 35, 255}},  // Forest
			{0.8, color.RGBA{139, 69, 19, 255}},   // Mountain
			{0.9, color.RGBA{100, 100, 100, 255}}, // Rock
			{0.98, color.RGBA{255, 250, 250, 255}}, // Snow
		}

		return NewColorMap(warped, stops...)
	}

	return refs, []string{"Wood", "Marble", "Clouds", "Terrain"}
}

var (
	WarpZoomLevels = []int{}
	WarpOutputFilename = "warp.png"
)
