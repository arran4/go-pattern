package pattern

import (
	"image"
	"image/color"
)

var (
	WaterOutputFilename         = "water.png"
	Water_surfaceOutputFilename = "water_surface.png"
)

func ExampleNewWater() image.Image {
	// 1. Base Noise: Simplex/Perlin noise (FBM)
	baseNoise := NewNoise(
		NoiseSeed(1),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        1,
			Octaves:     6,
			Persistence: 0.5,
			Lacunarity:  2.0,
			Frequency:   0.03,
		}),
	)

	// 2. Flow Maps: We simulate flow by warping the noise.
	// We'll use another lower frequency noise as the vector field (x/y displacement).
	// Since Warp takes one image for X and Y, we can use the same noise or different ones.
	// Let's create a "flow" map.
	flowX := NewNoise(
		NoiseSeed(2),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:      2,
			Frequency: 0.01,
		}),
	)

	// Apply warp
	warped := NewWarp(baseNoise,
		WarpDistortionX(flowX),
		WarpDistortionY(flowX), // Using same for simplicity, or could offset.
		WarpScale(20.0),
	)

	// 3. Normal Map: Convert the heightmap to normals
	normals := NewNormalMap(warped, NormalMapStrength(4.0))

	// 4. Colorization: We can use the normal map directly (it looks cool/techy),
	// or we can try to render it. But the prompt asked for "normals + flow maps".
	// Usually water is rendered with reflection/refraction which needs a shader.
	// Here we can output the normal map as the representation of the water surface.
	// Or we can blend it with a blue tint to make it look like water.

	waterBlue := color.RGBA{0, 0, 100, 255}
	waterTint := NewRect(SetFillColor(waterBlue))

	// Blend normals with blue using Overlay or SoftLight
	// Overlay might be too harsh for normals.
	// Let's just return the normal map as it is the "texture" of water surface.
	// Or maybe "Multiply" the blue with the normal map to tint it.

	blended := NewBlend(normals, waterTint, BlendAverage)

	return blended
}

func ExampleNewWater_surface() image.Image {
	// A variation showing just the normal map which is often what is used in game engines.
	baseNoise := NewNoise(
		NoiseSeed(42),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        42,
			Octaves:     5,
			Persistence: 0.6,
			Lacunarity:  2.0,
			Frequency:   0.04,
		}),
	)

	// Strong warp for "choppy" water
	distortion := NewNoise(
		NoiseSeed(100),
		SetNoiseAlgorithm(&PerlinNoise{Seed: 100, Frequency: 0.02}),
	)

	warped := NewWarp(baseNoise, WarpDistortion(distortion), WarpScale(30.0))

	return NewNormalMap(warped, NormalMapStrength(8.0))
}

func GenerateWater(rect image.Rectangle) image.Image {
	return ExampleNewWater()
}

func GenerateWater_surface(rect image.Rectangle) image.Image {
	return ExampleNewWater_surface()
}

func init() {
	GlobalGenerators["Water"] = GenerateWater
	GlobalGenerators["Water_surface"] = GenerateWater_surface
}
