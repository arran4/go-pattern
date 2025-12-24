package pattern

import (
	"image"
	"image/color"
)

func init() {
	GlobalGenerators["Shojo"] = GenerateShojo
	GlobalReferences["Shojo"] = GenerateShojoReferences
}

// GenerateShojo generates a Shojo Sparkles pattern.
func GenerateShojo(rect image.Rectangle) image.Image {
	return NewShojo(func(i any) {
		if p, ok := i.(interface{ SetBounds(image.Rectangle) }); ok {
			p.SetBounds(rect)
		}
	})
}

// GenerateShojoReferences generates reference images for the Shojo Sparkles pattern.
func GenerateShojoReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{}, []string{}
}

// NewDemoShojo produces a demo variant for readme.md.
func NewDemoShojo(ops ...func(any)) image.Image {
	return NewShojo(ops...)
}

// ExampleNewShojo_pink demonstrates a pink variant.
func ExampleNewShojo_pink() image.Image {
	return NewShojo(
		SetSpaceColor(color.RGBA{20, 0, 10, 255}), // Dark red/brown bg
		SetFillColor(color.RGBA{255, 200, 220, 255}), // Pink sparkles
	)
}

// ExampleNewShojo_blue demonstrates a blue variant.
func ExampleNewShojo_blue() image.Image {
	return NewShojo(
		SetSpaceColor(color.RGBA{0, 0, 40, 255}), // Dark blue bg
		SetFillColor(color.RGBA{200, 220, 255, 255}), // Blueish sparkles
	)
}
