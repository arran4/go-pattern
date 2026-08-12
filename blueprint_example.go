package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	BlueprintOutputFilename = "blueprint.png"
	BlueprintZoomLevels     = []int{}
)

const BlueprintBaseLabel = "Blueprint"

func init() {
	RegisterGenerator("Blueprint", func(bounds image.Rectangle) image.Image {
		return GenerateBlueprint(bounds)
	})
	RegisterReferences("Blueprint", GenerateBlueprintReferences)
}

// ExampleNewBlueprint renders a technical blueprint grid and saves it to blueprint.png.
func ExampleNewBlueprint() {
	img := NewBlueprint(
		SetBounds(image.Rect(0, 0, 300, 300)),
		SetBlueprintCellSize(10),
		SetBlueprintMajorCellSize(50),
		SetBlueprintLineWidth(1),
		SetBlueprintMajorLineWidth(2),
		SetBlueprintBackgroundColor(color.RGBA{R: 21, G: 67, B: 122, A: 255}),
		SetBlueprintMajorLineColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		SetBlueprintMinorLineColor(color.RGBA{R: 124, G: 167, B: 214, A: 255}),
	)

	f, err := os.Create(BlueprintOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			panic(cerr)
		}
	}()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

// GenerateBlueprint returns a preconfigured blueprint pattern for the registry and CLI.
func GenerateBlueprint(b image.Rectangle) image.Image {
	if b.Dx() == 0 || b.Dy() == 0 {
		b = image.Rect(0, 0, 300, 300)
	}

	return NewBlueprint(
		SetBounds(b),
		SetBlueprintCellSize(10),
		SetBlueprintMajorCellSize(50),
		SetBlueprintLineWidth(1),
		SetBlueprintMajorLineWidth(2),
		SetBlueprintBackgroundColor(color.RGBA{R: 21, G: 67, B: 122, A: 255}),
		SetBlueprintMajorLineColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		SetBlueprintMinorLineColor(color.RGBA{R: 124, G: 167, B: 214, A: 255}),
	)
}

func GenerateBlueprintReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	refs := map[string]func(image.Rectangle) image.Image{
		"Blueprint": GenerateBlueprint,
	}
	return refs, []string{"Blueprint"}
}
