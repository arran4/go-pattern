package pattern

import (
	"image"
	"image/color"
)

var (
	BrickOutputFilename = "brick.png"
	BrickZoomLevels     = []int{}
	BrickBaseLabel      = "Brick"
	BrickOrder          = 20
)

var (
	Brick_texturesOutputFilename = "brick_textures.png"
	Brick_texturesZoomLevels     = []int{}
	Brick_texturesBaseLabel      = "Brick Textures"
)

var (
	Brick_stoneOutputFilename = "brick_stone.png"
	Brick_stoneZoomLevels     = []int{}
	Brick_stoneBaseLabel      = "Stone Wall"
)

func init() {
	RegisterGenerator("Brick", func(bounds image.Rectangle) image.Image {
		return NewBrick(SetBounds(bounds))
	})
	RegisterReferences("Brick", GenerateBrickReferences)

	RegisterGenerator("Brick_textures", func(bounds image.Rectangle) image.Image {
		return ExampleNewBrick_textures()
	})
	RegisterGenerator("Brick_stone", func(bounds image.Rectangle) image.Image {
		return ExampleNewBrick_stone()
	})
}

func GenerateBrickReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	refs := make(map[string]func(image.Rectangle) image.Image)
	var names []string

	names = append(names, "Basic", "Textures", "Stone")

	refs["Basic"] = func(r image.Rectangle) image.Image {
		return NewBrick(
			SetBrickSize(40, 20),
			SetMortarSize(4),
		)
	}

	refs["Textures"] = func(r image.Rectangle) image.Image {
		return ExampleNewBrick_textures()
	}

	refs["Stone"] = func(r image.Rectangle) image.Image {
		return ExampleNewBrick_stone()
	}

	return refs, names
}

// ExampleNewBrick creates a basic brick pattern.
// Output:
func ExampleNewBrick() image.Image {
	return NewBrick(
		SetBrickSize(50, 20),
		SetMortarSize(4),
	)
}

// ExampleNewBrick_textures demonstrates using different textures for bricks and mortar.
func ExampleNewBrick_textures() image.Image {
	// Bricks with variations
	// Create 3 variations of brick textures using Noise
	var bricks []image.Image
	for i := 0; i < 3; i++ {
		// Noise with different seeds to ensure different texture per variant
		noise := NewNoise(SetNoiseAlgorithm(&PerlinNoise{
			Seed:      int64(i*100 + 1),
			Frequency: 0.1,
		}))

		// Tint the noise red/brown
		colored := NewColorMap(noise,
			ColorStop{0.0, color.RGBA{100, 30, 30, 255}},
			ColorStop{1.0, color.RGBA{180, 60, 50, 255}},
		)
		bricks = append(bricks, colored)
	}

	// Mortar texture: grey noise
	mortar := NewColorMap(
		NewNoise(SetNoiseAlgorithm(&PerlinNoise{
			Seed:      999,
			Frequency: 0.5,
		})),
		ColorStop{0.0, color.RGBA{180, 180, 180, 255}},
		ColorStop{1.0, color.RGBA{220, 220, 220, 255}},
	)

	return NewBrick(
		SetBrickSize(60, 25),
		SetMortarSize(3),
		SetBrickImages(bricks...),
		SetMortarImage(mortar),
	)
}

// ExampleNewBrick_stone demonstrates a stone-like wall using grey colors and size variations.
func ExampleNewBrick_stone() image.Image {
	// Create "Stone" textures
	var stones []image.Image
	for i := 0; i < 4; i++ {
		noise := NewNoise(SetNoiseAlgorithm(&PerlinNoise{
			Seed:      int64(i*50 + 123),
			Frequency: 0.2,
		}))
		// Grey/Blueish stone colors
		colored := NewColorMap(noise,
			ColorStop{0.0, color.RGBA{80, 80, 90, 255}},
			ColorStop{0.6, color.RGBA{120, 120, 130, 255}},
			ColorStop{1.0, color.RGBA{160, 160, 170, 255}},
		)
		stones = append(stones, colored)
	}

	mortar := NewUniform(color.RGBA{50, 50, 50, 255})

	// Larger bricks/stones
	return NewBrick(
		SetBrickSize(40, 30),
		SetMortarSize(6),
		SetBrickImages(stones...),
		SetMortarImage(mortar),
		SetBrickOffset(0.3), // Non-standard offset
	)
}

// Helper for uniform color image
func NewUniform(c color.Color) image.Image {
	return &image.Uniform{C: c}
}
