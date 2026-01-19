package pattern

import (
	"image"
	"image/color"
)

var (
	NormalMapOutputFilename        = "normal_map.png"
	NormalMap_sphereOutputFilename = "normal_map_sphere.png"
)

func ExampleNewNormalMap() image.Image {
	// Create a height map using Perlin noise
	noise := NewNoise(
		NoiseSeed(123),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        123,
			Octaves:     4,
			Persistence: 0.5,
			Lacunarity:  2.0,
			Frequency:   0.05,
		}),
	)

	// Convert to normal map with strength 5.0
	return NewNormalMap(noise, NormalMapStrength(5.0))
}

func ExampleNewNormalMap_sphere() image.Image {
	// A simple sphere gradient to show curvature normals
	grad := NewRadialGradient(
		GradientCenter(0.5, 0.5),
		SetStartColor(color.White),
		SetEndColor(color.Black),
	)

	// Increase strength significantly to visualize the curve on a smooth gradient
	return NewNormalMap(grad, NormalMapStrength(30.0))
}

// Generators need to match func(image.Rectangle) image.Image signature
func GenerateNormalMap(rect image.Rectangle) image.Image {
	return ExampleNewNormalMap()
}

func GenerateNormalMap_sphere(rect image.Rectangle) image.Image {
	return ExampleNewNormalMap_sphere()
}

func init() {
	GlobalGenerators["NormalMap"] = GenerateNormalMap
	GlobalGenerators["NormalMap_sphere"] = GenerateNormalMap_sphere
}
