package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	ShojoOutputFilename = "shojo.png"
	ShojoZoomLevels     = []int{}
	ShojoOrder          = 25
)

const ShojoBaseLabel = "Shojo"

func init() {
	GlobalGenerators["Shojo"] = GenerateShojo
	GlobalReferences["Shojo"] = GenerateShojoReferences

	GlobalGenerators["Shojo_pink"] = GenerateShojo_pink
	GlobalReferences["Shojo_pink"] = GenerateShojoReferences

	GlobalGenerators["Shojo_blue"] = GenerateShojo_blue
	GlobalReferences["Shojo_blue"] = GenerateShojoReferences
}

// GenerateShojo generates a Shojo Sparkles pattern.
func GenerateShojo(rect image.Rectangle) image.Image {
	return NewShojo(SetBounds(rect))
}

// GenerateShojoReferences generates reference images for the Shojo Sparkles pattern.
func GenerateShojoReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{}, []string{}
}

// ExampleNewShojo produces a demo variant for readme.md.
func ExampleNewShojo() {
	i := NewShojo()
	f, err := os.Create(ShojoOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

var (
	Shojo_pinkOutputFilename = "shojo_pink.png"
	Shojo_pinkZoomLevels     = []int{}
	Shojo_pinkOrder          = 26
)

const Shojo_pinkBaseLabel = "Pink Variant"

func GenerateShojo_pink(rect image.Rectangle) image.Image {
	return NewShojo(
		SetSpaceColor(color.RGBA{20, 0, 10, 255}),    // Dark red/brown bg
		SetFillColor(color.RGBA{255, 200, 220, 255}), // Pink sparkles
		SetBounds(rect),
	)
}

// ExampleNewShojo_pink demonstrates a pink variant.
func ExampleNewShojo_pink() {
	i := NewShojo(
		SetSpaceColor(color.RGBA{20, 0, 10, 255}),    // Dark red/brown bg
		SetFillColor(color.RGBA{255, 200, 220, 255}), // Pink sparkles
	)
	f, err := os.Create(Shojo_pinkOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

var (
	Shojo_blueOutputFilename = "shojo_blue.png"
	Shojo_blueZoomLevels     = []int{}
	Shojo_blueOrder          = 27
)

const Shojo_blueBaseLabel = "Blue Variant"

func GenerateShojo_blue(rect image.Rectangle) image.Image {
	return NewShojo(
		SetSpaceColor(color.RGBA{0, 0, 40, 255}),     // Dark blue bg
		SetFillColor(color.RGBA{200, 220, 255, 255}), // Blueish sparkles
		SetBounds(rect),
	)
}

// ExampleNewShojo_blue demonstrates a blue variant.
func ExampleNewShojo_blue() {
	i := NewShojo(
		SetSpaceColor(color.RGBA{0, 0, 40, 255}),     // Dark blue bg
		SetFillColor(color.RGBA{200, 220, 255, 255}), // Blueish sparkles
	)
	f, err := os.Create(Shojo_blueOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}
