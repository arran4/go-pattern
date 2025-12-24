package pattern

import (
	"image"
	"image/color"
)

var (
	GridOutputFilename = "grid.png"
	GridZoomLevels     = []int{}
)

const (
	GridOrder     = 5
	GridBaseLabel = "Grid"
)

func init() {
	RegisterGenerator("Grid", func(bounds image.Rectangle) image.Image {
		return ExampleNewGrid(SetBounds(bounds))
	})
	RegisterReferences("Grid", BootstrapGridReferences)
}

func BootstrapGridReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Rows": func(bounds image.Rectangle) image.Image {
			return NewDemoGridRows(SetBounds(bounds))
		},
		"Cols": func(bounds image.Rectangle) image.Image {
			return NewDemoGridColumns(SetBounds(bounds))
		},
		"Fixed": func(bounds image.Rectangle) image.Image {
			return NewDemoGridFixed(SetBounds(bounds))
		},
		"Bounded": func(bounds image.Rectangle) image.Image {
			return NewDemoGridBounded(SetBounds(bounds))
		},
		"Advanced": func(bounds image.Rectangle) image.Image {
			return NewDemoGridAdvanced(SetBounds(bounds))
		},
	}, []string{"Rows", "Cols", "Fixed", "Bounded", "Advanced"}
}

func ExampleNewGrid(ops ...func(any)) image.Image {
	// Example 1: Simple 2x2 grid with Gophers
	// Shrink the Gopher so it fits better
	gopher := NewScale(NewGopher(), ScaleToRatio(0.25))

	args := []any{
		Row(Cell(gopher), Cell(gopher)),
		Row(Cell(gopher), Cell(gopher)),
	}
	for _, op := range ops {
		args = append(args, op)
	}

	// Create a grid with explicit Rows
	return NewGrid(args...)
}

func NewDemoGridColumns(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.25))

	args := []any{
		Column(Cell(gopher), Cell(gopher)),
		Column(Cell(gopher), Cell(gopher)),
	}
	for _, op := range ops {
		args = append(args, op)
	}

	// Create a grid with explicit Columns
	return NewGrid(args...)
}

func NewDemoGridFixed(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.25))

	args := []any{
		FixedSize(400, 400),
		CellPos(0, 0, gopher),
		CellPos(1, 1, gopher),
	}
	for _, op := range ops {
		args = append(args, op)
	}

	return NewGrid(args...)
}

// Mock Bounded object for testing
type boundedGopher struct {
	image.Image
	bounds Bounds
}

func (b *boundedGopher) PatternBounds() Bounds {
	return b.bounds
}

func NewDemoGridBounded(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.25))

	// Create a bounded version that claims to be larger
	zero := 0
	twoHundred := 200

	b := &boundedGopher{
		Image: gopher,
		bounds: Bounds{
			Left:   &Range{Low: &zero, High: &zero},
			Right:  &Range{Low: &twoHundred, High: &twoHundred},
			Top:    &Range{Low: &zero, High: &zero},
			Bottom: &Range{Low: &twoHundred, High: &twoHundred},
		},
	}

	args := []any{
		Row(Cell(b), Cell(gopher)),
	}
	for _, op := range ops {
		args = append(args, op)
	}

	return NewGrid(args...)
}

func NewDemoGridAdvanced(ops ...func(any)) image.Image {
	gopher := NewGopher()
	// 1. Scaled Gopher (Small)
	smallGopher := NewScale(gopher, ScaleToRatio(0.2))

	// 2. Cropped Gopher (Head only)
	// Gopher is roughly 500x700. Head is at top.
	head := NewCrop(gopher, image.Rect(100, 50, 400, 350))
	smallHead := NewScale(head, ScaleToRatio(0.3))

	// 3. Tiled Gopher
	// Tile a small gopher into a 100x100 box
	tinyGopher := NewScale(gopher, ScaleToSize(20, 30))
	tiled := NewTile(tinyGopher, image.Rect(0, 0, 100, 100))

	// 4. Text Label
	// White background for visibility
	txt := NewText("Hello Grid", TextSize(18), TextColorColor(color.Black), TextBackgroundColorColor(color.White))

	// Center the text in a 150x50 box
	centeredTxt := NewCenter(txt, 150, 50, image.NewUniform(color.White))

	// 5. Padding with Checker background
	checker := NewChecker(color.RGBA{200, 200, 200, 255}, color.White)
	padded := NewPadding(smallGopher, PaddingMargin(20), PaddingBackground(checker))

	args := []any{
		Row(Cell(centeredTxt), Cell(padded)),
		Row(Cell(smallHead), Cell(tiled)),
	}
	for _, op := range ops {
		args = append(args, op)
	}

	return NewGrid(args...)
}
