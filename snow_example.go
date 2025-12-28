package pattern

import (
	"image"
	"image/color"
)

var (
	SnowOutputFilename = "snow.png"
	Snow_tracksOutputFilename = "snow_tracks.png"
)

func ExampleNewSnow() image.Image {
	// 1. Soft base noise (drifts) - Bright white/grey
	drifts := NewNoise(
		NoiseSeed(505),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:      505,
			Frequency: 0.01,
		}),
	)

	snowColor := NewColorMap(drifts,
		ColorStop{Position: 0.0, Color: color.RGBA{240, 240, 250, 255}}, // Slight blue-grey shadow
		ColorStop{Position: 1.0, Color: color.White},
	)

	// 2. Sparkle: Use white/blue dots.
	// We can use Scatter pattern to place small bright dots.
	// But let's fix the noise approach.
	// High frequency noise, thresholded.
	sparkleNoise := NewNoise(
		NoiseSeed(606),
		SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.8}),
	)

	// We want sparkles to be White/Blue on Transparent background.
	// Then Overlay or Screen them.
	// If background is Transparent, Screen (1-(1-A)*(1-B)) of Snow(A) and Transparent(B=0) -> A.
	// So sparkles need to be Additive.
	// Or we can just use Mix/Over.

	sparkles := NewColorMap(sparkleNoise,
		ColorStop{Position: 0.0, Color: color.Transparent},
		ColorStop{Position: 0.9, Color: color.Transparent},
		ColorStop{Position: 0.92, Color: color.RGBA{200, 220, 255, 255}}, // Blue tint
		ColorStop{Position: 1.0, Color: color.White},
	)

	// Use BlendNormal (Over) for sparkles
	return NewBlend(snowColor, sparkles, BlendNormal)
}

func ExampleNewSnow_tracks() image.Image {
	snow := ExampleNewSnow()

	// 3. Compression Tracks: Blueish/Grey depression.
	tracks := NewCrossHatch(
		SetLineColor(color.RGBA{200, 210, 230, 255}), // Icy blue/grey
		SetSpaceColor(color.White), // Neutral for Multiply
		SetLineSize(15),
		SetSpaceSize(80),
		SetAngles(25, 35), // Overlapping tracks
	)

	// Distort tracks
	distort := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.05}))
	organicTracks := NewWarp(tracks, WarpDistortion(distort), WarpScale(8.0))

	// Multiply tracks onto snow
	// LineColor (Blueish) * Snow (White) -> Blueish.
	// SpaceColor (White) * Snow (White) -> White.
	return NewBlend(snow, organicTracks, BlendMultiply)
}

func GenerateSnow(rect image.Rectangle) image.Image {
	return ExampleNewSnow()
}

func GenerateSnow_tracks(rect image.Rectangle) image.Image {
	return ExampleNewSnow_tracks()
}

func init() {
	GlobalGenerators["Snow"] = GenerateSnow
	GlobalGenerators["Snow_tracks"] = GenerateSnow_tracks
}
