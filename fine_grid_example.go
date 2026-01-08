package pattern

import (
	"image"
	"image/png"
	"os"
)

var (
	FineGridOutputFilename = "fine_grid.png"
	FineGridZoomLevels     = []int{}
)

const FineGridBaseLabel = "FineGrid"

func init() {
	RegisterGenerator("FineGrid", func(bounds image.Rectangle) image.Image {
		return GenerateFineGrid(bounds)
	})
	RegisterReferences("FineGrid", GenerateFineGridReferences)
}

// ExampleNewFineGrid renders a neon grid with glow and saves it to fine_grid.png.
func ExampleNewFineGrid() {
	img := NewFineGrid(
		SetBounds(image.Rect(0, 0, 640, 640)),
		SetFineGridCellSize(12),
		SetFineGridGlowRadius(3.5),
		SetFineGridHue(205),
		SetFineGridAberration(1),
		SetFineGridGlowStrength(0.9),
		SetFineGridLineStrength(1.4),
		SetFineGridBackgroundFade(0.0),
	)

	f, err := os.Create(FineGridOutputFilename)
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

// GenerateFineGrid returns a preconfigured fine grid for the registry and CLI.
func GenerateFineGrid(b image.Rectangle) image.Image {
	if b.Dx() == 0 || b.Dy() == 0 {
		b = image.Rect(0, 0, 640, 640)
	}

	return NewFineGrid(
		SetBounds(b),
		SetFineGridCellSize(12),
		SetFineGridGlowRadius(3.5),
		SetFineGridHue(205),
		SetFineGridAberration(1),
		SetFineGridGlowStrength(0.9),
		SetFineGridLineStrength(1.4),
		SetFineGridBackgroundFade(0.05),
	)
}

// GenerateFineGridWarm variant showcases hue customization.
func GenerateFineGridWarm(b image.Rectangle) image.Image {
	if b.Dx() == 0 || b.Dy() == 0 {
		b = image.Rect(0, 0, 640, 640)
	}

	return NewFineGrid(
		SetBounds(b),
		SetFineGridCellSize(10),
		SetFineGridGlowRadius(4.0),
		SetFineGridHue(35),
		SetFineGridAberration(1),
		SetFineGridGlowStrength(0.9),
		SetFineGridLineStrength(1.2),
		SetFineGridBackgroundFade(0.08),
	)
}

// GenerateFineGridMagenta variant showcases a denser grid and different hue.
func GenerateFineGridMagenta(b image.Rectangle) image.Image {
	if b.Dx() == 0 || b.Dy() == 0 {
		b = image.Rect(0, 0, 640, 640)
	}

	return NewFineGrid(
		SetBounds(b),
		SetFineGridCellSize(8),
		SetFineGridGlowRadius(3.0),
		SetFineGridHue(305),
		SetFineGridAberration(1),
		SetFineGridGlowStrength(0.85),
		SetFineGridLineStrength(1.1),
		SetFineGridBackgroundFade(0.05),
	)
}

func GenerateFineGridReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	refs := map[string]func(image.Rectangle) image.Image{
		"FineGrid": GenerateFineGrid,
		"Warm":     GenerateFineGridWarm,
		"Magenta":  GenerateFineGridMagenta,
	}
	return refs, []string{"FineGrid", "Warm", "Magenta"}
}
