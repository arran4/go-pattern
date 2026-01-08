package pattern

import (
	"image"
	"image/png"
	"os"
)

var WindRidgesOutputFilename = "wind.png"

const WindRidgesBaseLabel = "WindRidges"

// ExampleNewWindRidges writes a wind-swept noise PNG showcasing parameterized streaks.
func ExampleNewWindRidges() {
	img := GenerateWindRidges(image.Rect(0, 0, 200, 200))

	// Output:

	f, err := os.Create(WindRidgesOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}

// GenerateWindRidges builds a directional streaked noise with soft shadow ridges.
func GenerateWindRidges(b image.Rectangle) image.Image {
	// Base white noise
	noise := NewNoise(
		SetBounds(b),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        2024,
			Frequency:   0.08,
			Octaves:     3,
			Persistence: 0.55,
		}),
	)

	return NewWindRidges(
		SetBounds(b),
		SetWindNoise(noise),
		SetWindAngle(28.0),
		SetStreakLength(22),
		SetWindContrast(1.3),
	)
}

// GenerateWindRidgesReferences returns variations focusing on the key tunables.
func GenerateWindRidgesReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Diagonal": func(b image.Rectangle) image.Image {
			return GenerateWindRidges(b)
		},
		"Crosswind": func(b image.Rectangle) image.Image {
			noise := NewNoise(SetBounds(b), SetNoiseAlgorithm(&PerlinNoise{
				Seed:        909,
				Frequency:   0.06,
				Octaves:     4,
				Persistence: 0.6,
			}))
			return NewWindRidges(
				SetBounds(b),
				SetWindNoise(noise),
				SetWindAngle(75.0),
				SetStreakLength(26),
				SetWindContrast(1.1),
			)
		},
		"HighContrast": func(b image.Rectangle) image.Image {
			return NewWindRidges(
				SetBounds(b),
				SetWindNoise(NewNoise(SetBounds(b), SetNoiseAlgorithm(&HashNoise{Seed: 777}))),
				SetWindAngle(10.0),
				SetStreakLength(16),
				SetWindContrast(1.6),
			)
		},
	}, []string{"Diagonal", "Crosswind", "HighContrast"}
}

func init() {
	RegisterGenerator(WindRidgesBaseLabel, GenerateWindRidges)
	RegisterReferences(WindRidgesBaseLabel, GenerateWindRidgesReferences)
}
