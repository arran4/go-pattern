package pattern

import (
	"fmt"
	"image"
	"image/color"
)

// Main Warp Example
// Variables defined here to ensure findStringVar finds them easily.
var (
	WarpOutputFilename = "warp.png"
	WarpZoomLevels     = []int{}
	WarpBaseLabel      = "Warp"
)

func ExampleNewWarp() {
	// Standard demo: Grid warped by noise
	// We want a visual that clearly shows the warping effect.
	// A checkerboard is good.

	checker := NewChecker(
		color.RGBA{200, 200, 200, 255},
		color.RGBA{50, 50, 50, 255},
	)

	// Distortion noise
	noise := NewNoise(NoiseSeed(99), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.03,
		Octaves:   2,
	}))

	// Apply Warp
	warped := NewWarp(checker,
		WarpDistortion(noise),
		WarpScale(10.0),
	)

	fmt.Println(warped.At(10, 10))
	// Output: {50 50 50 255}
}

// Wood Example

var (
	Warp_woodOutputFilename = "warp_wood.png"
	Warp_woodZoomLevels     = []int{}
	Warp_woodBaseLabel      = "Wood"
)

func ExampleNewWarp_wood() image.Image {
	woodLight := color.RGBA{222, 184, 135, 255}
	woodDark := color.RGBA{139, 69, 19, 255}

	colors := []color.Color{}
	steps := 20
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		r := uint8(float64(woodLight.R)*(1-t) + float64(woodDark.R)*t)
		g := uint8(float64(woodLight.G)*(1-t) + float64(woodDark.G)*t)
		b := uint8(float64(woodLight.B)*(1-t) + float64(woodDark.B)*t)
		colors = append(colors, color.RGBA{r, g, b, 255})
	}
	for i := steps - 1; i >= 0; i-- {
		colors = append(colors, colors[i])
	}

	rings := NewConcentricRings(colors)

	noiseLow := NewNoise(NoiseSeed(123), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.02,
		Octaves:   2,
	}))

	// Apply Warp
	return NewWarp(rings,
		WarpDistortion(noiseLow),
		WarpScale(15.0),
	)
}

// Marble Example

var (
	Warp_marbleOutputFilename = "warp_marble.png"
	Warp_marbleZoomLevels     = []int{}
	Warp_marbleBaseLabel      = "Marble"
)

func ExampleNewWarp_marble() image.Image {
	colors := []color.Color{
		color.RGBA{240, 240, 245, 255},
		color.RGBA{240, 240, 245, 255},
		color.RGBA{240, 240, 245, 255},
		color.RGBA{200, 200, 210, 255},
		color.RGBA{100, 100, 110, 255},
		color.RGBA{200, 200, 210, 255},
	}
	stripes := NewModuloStripe(colors)

	noise := NewNoise(NoiseSeed(456), SetNoiseAlgorithm(&PerlinNoise{
		Frequency:   0.04,
		Octaves:     4,
		Persistence: 0.6,
	}))

	return NewWarp(stripes,
		WarpDistortion(noise),
		WarpScale(30.0),
	)
}

// Clouds Example

var (
	Warp_cloudsOutputFilename = "warp_clouds.png"
	Warp_cloudsZoomLevels     = []int{}
	Warp_cloudsBaseLabel      = "Clouds"
)

func ExampleNewWarp_clouds() image.Image {
	baseNoise := NewNoise(NoiseSeed(777), SetNoiseAlgorithm(&PerlinNoise{
		Frequency:   0.02,
		Octaves:     4,
		Persistence: 0.5,
	}))

	warpNoise := NewNoise(NoiseSeed(888), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.02,
		Octaves:   2,
	}))

	warped := NewWarp(baseNoise,
		WarpDistortion(warpNoise),
		WarpScale(50.0),
	)

	stops := []ColorStop{
		{0.0, color.RGBA{0, 100, 200, 255}},
		{0.4, color.RGBA{100, 150, 255, 255}},
		{0.6, color.RGBA{255, 255, 255, 255}},
		{1.0, color.RGBA{255, 255, 255, 255}},
	}

	return NewColorMap(warped, stops...)
}

// Terrain Example

var (
	Warp_terrainOutputFilename = "warp_terrain.png"
	Warp_terrainZoomLevels     = []int{}
	Warp_terrainBaseLabel      = "Terrain"
)

func ExampleNewWarp_terrain() image.Image {
	fbm := func(seed int64) image.Image {
		return NewNoise(NoiseSeed(seed), SetNoiseAlgorithm(&PerlinNoise{
			Frequency:   0.015,
			Octaves:     6,
			Persistence: 0.5,
			Lacunarity:  2.0,
		}))
	}

	base := fbm(101)

	warp := NewNoise(NoiseSeed(202), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.01,
		Octaves:   2,
	}))

	warped := NewWarp(base,
		WarpDistortion(warp),
		WarpScale(80.0),
	)

	stops := []ColorStop{
		{0.0, color.RGBA{0, 0, 150, 255}},
		{0.2, color.RGBA{0, 50, 200, 255}},
		{0.22, color.RGBA{240, 230, 140, 255}},
		{0.3, color.RGBA{34, 139, 34, 255}},
		{0.6, color.RGBA{107, 142, 35, 255}},
		{0.8, color.RGBA{139, 69, 19, 255}},
		{0.9, color.RGBA{100, 100, 100, 255}},
		{0.98, color.RGBA{255, 250, 250, 255}},
	}

	return NewColorMap(warped, stops...)
}

func init() {
	GlobalGenerators[WarpBaseLabel] = GenerateWarp
	GlobalReferences[WarpBaseLabel] = GenerateWarpReferences

	// Register variants using the exact suffix ID expected by bootstrap (Warp_wood, etc.)
	GlobalGenerators["Warp_wood"] = GenerateWarp_wood
	GlobalReferences["Warp_wood"] = GenerateWarpReferences_Empty

	GlobalGenerators["Warp_marble"] = GenerateWarp_marble
	GlobalReferences["Warp_marble"] = GenerateWarpReferences_Empty

	GlobalGenerators["Warp_clouds"] = GenerateWarp_clouds
	GlobalReferences["Warp_clouds"] = GenerateWarpReferences_Empty

	GlobalGenerators["Warp_terrain"] = GenerateWarp_terrain
	GlobalReferences["Warp_terrain"] = GenerateWarpReferences_Empty
}

// Generator Wrappers

func GenerateWarp(rect image.Rectangle) image.Image {
	checker := NewChecker(
		color.RGBA{200, 200, 200, 255},
		color.RGBA{50, 50, 50, 255},
	)

	noise := NewNoise(NoiseSeed(99), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.03,
		Octaves:   2,
	}))

	return NewWarp(checker,
		WarpDistortion(noise),
		WarpScale(10.0),
	)
}

func GenerateWarp_wood(rect image.Rectangle) image.Image {
	return ExampleNewWarp_wood()
}

func GenerateWarp_marble(rect image.Rectangle) image.Image {
	return ExampleNewWarp_marble()
}

func GenerateWarp_clouds(rect image.Rectangle) image.Image {
	return ExampleNewWarp_clouds()
}

func GenerateWarp_terrain(rect image.Rectangle) image.Image {
	return ExampleNewWarp_terrain()
}

func GenerateWarpReferences_Empty() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{}, []string{}
}

// GenerateWarpReferences registers the examples for documentation generation (for the main Warp pattern references if needed).
func GenerateWarpReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{}, []string{}
}
