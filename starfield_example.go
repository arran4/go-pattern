package pattern

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var StarfieldOutputFilename = "starfield.png"

const StarfieldBaseLabel = "Starfield"

func ExampleNewStarfield() {
	img := GenerateStarfield(image.Rect(0, 0, 150, 150))
	f, err := os.Create(StarfieldOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}

func GenerateStarfield(b image.Rectangle) image.Image {
	// 1. Background: Deep Space (Black/Very Dark Blue)
	bg := NewColorMap(
		NewNoise(SetBounds(b), NoiseSeed(999), SetNoiseAlgorithm(&PerlinNoise{Frequency:0.01})),
		ColorStop{0.0, color.RGBA{0, 0, 10, 255}},
		ColorStop{1.0, color.RGBA{10, 0, 20, 255}},
	)

	// 2. Stars: Scatter white dots
	// We use Scatter pattern with a custom generator.

    starGen := func(u, v float64, hash uint64) (color.Color, float64) {
        // Randomize size slightly
        size := 1.0 + float64(hash % 10) / 10.0 // 1.0 to 2.0
        dist := math.Sqrt(u*u + v*v)
        if dist < size {
             return color.White, 0
        }
        return color.Transparent, 0
    }

    starsSmall := NewScatter(
        SetScatterGenerator(starGen),
        SetScatterFrequency(0.05), // 1/0.05 = 20 pixels cell
        SetScatterDensity(0.3),
        func(i any) {
            if p, ok := i.(*Scatter); ok {
                p.Seed = 1
            }
        },
    )

	// Blend Stars onto BG
	// Blend mode Screen or Add
	withStars := NewBlend(bg, starsSmall, BlendScreen)

	// 3. Nebula: Colored Noise
	nebula := NewNoise(
		SetBounds(b),
		NoiseSeed(888),
		SetNoiseAlgorithm(&PerlinNoise{
			Frequency: 0.015,
			Octaves: 3,
		}),
	)
	nebulaColor := NewColorMap(nebula,
		ColorStop{0.0, color.RGBA{0, 0, 0, 0}}, // Transparent
		ColorStop{0.4, color.RGBA{0, 0, 0, 0}},
		ColorStop{0.6, color.RGBA{100, 0, 100, 100}}, // Purple, semi-transparent
		ColorStop{0.8, color.RGBA{0, 0, 150, 150}}, // Blue
		ColorStop{1.0, color.RGBA{0, 100, 200, 180}}, // Bright Blue
	)

	// Blend Nebula
	final := NewBlend(withStars, nebulaColor, BlendScreen)

	return final
}

func init() {
	RegisterGenerator(StarfieldBaseLabel, GenerateStarfield)
	RegisterReferences(StarfieldBaseLabel, func() (map[string]func(image.Rectangle) image.Image, []string) {
		return map[string]func(image.Rectangle) image.Image{}, []string{}
	})
}
