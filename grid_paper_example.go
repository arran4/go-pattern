package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	GridPaperOutputFilename = "grid_paper.png"
	GridPaperZoomLevels     = []int{}
)

const GridPaperBaseLabel = "GridPaper"

func init() {
	RegisterGenerator("GridPaper", func(bounds image.Rectangle) image.Image {
		return GenerateGridPaper(bounds)
	})
	RegisterReferences("GridPaper", GenerateGridPaperReferences)
}

// ExampleNewGridPaper renders a grid paper pattern and saves it to grid_paper.png.
func ExampleNewGridPaper() {
	img := NewGridPaper(
		SetBounds(image.Rect(0, 0, 300, 300)),
		SetGridPaperCellSize(10),
		SetGridPaperMajorCellSize(50),
		SetGridPaperLineWidth(1),
		SetGridPaperMajorLineWidth(2),
		SetGridPaperBackgroundColor(color.RGBA{R: 245, G: 245, B: 245, A: 255}),
		SetGridPaperMajorLineColor(color.RGBA{R: 80, G: 120, B: 180, A: 255}),
		SetGridPaperMinorLineColor(color.RGBA{R: 140, G: 180, B: 230, A: 255}),
	)

	f, err := os.Create(GridPaperOutputFilename)
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

// GenerateGridPaper returns a preconfigured grid paper pattern for the registry and CLI.
func GenerateGridPaper(b image.Rectangle) image.Image {
	if b.Dx() == 0 || b.Dy() == 0 {
		b = image.Rect(0, 0, 300, 300)
	}

	return NewGridPaper(
		SetBounds(b),
		SetGridPaperCellSize(10),
		SetGridPaperMajorCellSize(50),
		SetGridPaperLineWidth(1),
		SetGridPaperMajorLineWidth(2),
		SetGridPaperBackgroundColor(color.RGBA{R: 245, G: 245, B: 245, A: 255}),
		SetGridPaperMajorLineColor(color.RGBA{R: 80, G: 120, B: 180, A: 255}),
		SetGridPaperMinorLineColor(color.RGBA{R: 140, G: 180, B: 230, A: 255}),
	)
}

func GenerateGridPaperReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	refs := map[string]func(image.Rectangle) image.Image{
		"GridPaper": GenerateGridPaper,
	}
	return refs, []string{"GridPaper"}
}
