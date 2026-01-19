package pattern

import (
	"image"
	"image/color"
)

// Main Clouds Example
var (
	CloudsOutputFilename = "clouds.png"
	CloudsZoomLevels     = []int{}
	CloudsBaseLabel      = "Clouds"
)

// ExampleNewClouds generates a generic cloud pattern.
func ExampleNewClouds() image.Image {
	return ExampleNewClouds_cumulus()
}

// -----------------------------------------------------------------------------
// Cumulus Clouds
// -----------------------------------------------------------------------------

var (
	Clouds_cumulusOutputFilename = "clouds_cumulus.png"
	Clouds_cumulusZoomLevels     = []int{}
	Clouds_cumulusBaseLabel      = "Cumulus"
)

// ExampleNewClouds_cumulus generates fluffy, white cumulus clouds on a blue sky.
// It uses Perlin noise with a specific color map that has a sharp transition
// from blue to white to simulate the defined edges of cumulus clouds.
func ExampleNewClouds_cumulus() image.Image {
	// 1. Base shape: Low frequency noise to define the cloud blobs
	noise := NewNoise(NoiseSeed(42), SetNoiseAlgorithm(&PerlinNoise{
		Frequency:   0.015,
		Octaves:     4,
		Persistence: 0.5,
		Lacunarity:  2.0,
	}))

	// 2. Color Map: Sky Blue -> White
	// We use a steep ramp around 0.5-0.6 to create distinct cloud shapes
	// rather than a smooth fog.
	return NewColorMap(noise,
		ColorStop{0.0, color.RGBA{100, 180, 255, 255}},  // Blue Sky
		ColorStop{0.4, color.RGBA{130, 200, 255, 255}},  // Light Sky
		ColorStop{0.55, color.RGBA{245, 245, 255, 255}}, // Cloud Edge (White-ish)
		ColorStop{0.7, color.RGBA{255, 255, 255, 255}},  // Cloud Body
		ColorStop{1.0, color.RGBA{230, 230, 240, 255}},  // Cloud Shadow/Density
	)
}

// -----------------------------------------------------------------------------
// Cirrus Clouds
// -----------------------------------------------------------------------------

var (
	Clouds_cirrusOutputFilename = "clouds_cirrus.png"
	Clouds_cirrusZoomLevels     = []int{}
	Clouds_cirrusBaseLabel      = "Cirrus"
)

// ExampleNewClouds_cirrus generates wispy, high-altitude cirrus clouds.
func ExampleNewClouds_cirrus() image.Image {
	// High frequency noise with high persistence to simulate wisps
	wispyNoise := NewNoise(NoiseSeed(103), SetNoiseAlgorithm(&PerlinNoise{
		Frequency:   0.05,
		Octaves:     6,
		Persistence: 0.7,
	}))

	return NewColorMap(wispyNoise,
		ColorStop{0.0, color.RGBA{20, 50, 150, 255}},   // Dark Blue Sky
		ColorStop{0.6, color.RGBA{50, 100, 200, 255}},  // Blue
		ColorStop{0.7, color.RGBA{150, 200, 255, 100}}, // Faint wisp
		ColorStop{1.0, color.RGBA{255, 255, 255, 200}}, // Bright wisp
	)
}

// -----------------------------------------------------------------------------
// Stormy Clouds
// -----------------------------------------------------------------------------

var (
	Clouds_stormOutputFilename = "clouds_storm.png"
	Clouds_stormZoomLevels     = []int{}
	Clouds_stormBaseLabel      = "Storm"
)

// ExampleNewClouds_storm generates dark, turbulent storm clouds.
// It blends multiple layers of noise to create depth and complexity.
func ExampleNewClouds_storm() image.Image {
	// Layer 1: Large, brooding shapes
	base := NewNoise(NoiseSeed(666), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.01,
		Octaves:   3,
	}))

	// Layer 2: Detailed turbulence
	detail := NewNoise(NoiseSeed(777), SetNoiseAlgorithm(&PerlinNoise{
		Frequency:   0.04,
		Octaves:     5,
		Persistence: 0.6,
	}))

	// Blend them: Overlay adds contrast
	blended := NewBlend(base, detail, BlendOverlay)

	// Map to stormy colors
	return NewColorMap(blended,
		ColorStop{0.0, color.RGBA{20, 20, 25, 255}},    // Darkest Grey
		ColorStop{0.4, color.RGBA{50, 50, 60, 255}},    // Dark Grey
		ColorStop{0.6, color.RGBA{80, 80, 90, 255}},    // Mid Grey
		ColorStop{0.8, color.RGBA{120, 120, 130, 255}}, // Light Grey highlights
		ColorStop{1.0, color.RGBA{160, 160, 170, 255}}, // Brightest peaks
	)
}

// -----------------------------------------------------------------------------
// Sunset Clouds
// -----------------------------------------------------------------------------

var (
	Clouds_sunsetOutputFilename = "clouds_sunset.png"
	Clouds_sunsetZoomLevels     = []int{}
	Clouds_sunsetBaseLabel      = "Sunset"
)

// ExampleNewClouds_sunset generates clouds illuminated by a setting sun.
// It uses a linear gradient for the sky background and Perlin noise for the clouds,
// blending them to simulate under-lighting.
func ExampleNewClouds_sunset() image.Image {
	// 1. Sky Gradient (Orange to Purple)
	sky := NewLinearGradient(
		SetStartColor(color.RGBA{255, 100, 50, 255}), // Orange/Red Horizon
		SetEndColor(color.RGBA{50, 20, 100, 255}),    // Purple/Blue Zenith
		GradientVertical(),
	)

	// 2. Cloud Shapes
	clouds := NewNoise(NoiseSeed(888), SetNoiseAlgorithm(&PerlinNoise{
		Frequency: 0.012,
		Octaves:   4,
	}))

	// Map cloud noise to alpha/color
	// We want the clouds to be dark at the bottom (shadow) and pink/gold at the edges
	cloudColor := NewColorMap(clouds,
		ColorStop{0.0, color.Black},                    // No clouds
		ColorStop{0.4, color.Black},                    // No clouds
		ColorStop{0.5, color.RGBA{80, 40, 60, 255}},    // Dark cloud base
		ColorStop{0.7, color.RGBA{200, 100, 80, 255}},  // Orange/Pink mid
		ColorStop{1.0, color.RGBA{255, 200, 100, 255}}, // Gold highlights
	)

	// 3. Composite Clouds over Sky using Screen blend mode for a glowing effect
	return NewBlend(sky, cloudColor, BlendScreen)
}

func init() {
	GlobalGenerators[CloudsBaseLabel] = GenerateClouds
	GlobalReferences[CloudsBaseLabel] = GenerateCloudsReferences

	GlobalGenerators["Clouds_cumulus"] = GenerateClouds_cumulus
	GlobalReferences["Clouds_cumulus"] = GenerateCloudsReferences_Empty

	GlobalGenerators["Clouds_cirrus"] = GenerateClouds_cirrus
	GlobalReferences["Clouds_cirrus"] = GenerateCloudsReferences_Empty

	GlobalGenerators["Clouds_storm"] = GenerateClouds_storm
	GlobalReferences["Clouds_storm"] = GenerateCloudsReferences_Empty

	GlobalGenerators["Clouds_sunset"] = GenerateClouds_sunset
	GlobalReferences["Clouds_sunset"] = GenerateCloudsReferences_Empty
}

func GenerateClouds(rect image.Rectangle) image.Image {
	return ExampleNewClouds()
}

func GenerateClouds_cumulus(rect image.Rectangle) image.Image {
	return ExampleNewClouds_cumulus()
}

func GenerateClouds_cirrus(rect image.Rectangle) image.Image {
	return ExampleNewClouds_cirrus()
}

func GenerateClouds_storm(rect image.Rectangle) image.Image {
	return ExampleNewClouds_storm()
}

func GenerateClouds_sunset(rect image.Rectangle) image.Image {
	return ExampleNewClouds_sunset()
}

func GenerateCloudsReferences_Empty() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{}, []string{}
}

func GenerateCloudsReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Cumulus": GenerateClouds_cumulus,
		"Cirrus":  GenerateClouds_cirrus,
		"Storm":   GenerateClouds_storm,
		"Sunset":  GenerateClouds_sunset,
	}, []string{"Cumulus", "Cirrus", "Storm", "Sunset"}
}
