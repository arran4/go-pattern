package pattern

import (
	"image"
	"image/color"
)

var (
	GridOutputFilename = "grid.png"
	GridZoomLevels     = []int{}
	GridOrder          = 5
	GridBaseLabel      = "Grid"
)

func init() {
	RegisterGenerator("Grid", func(bounds image.Rectangle) image.Image {
		return ExampleNewGrid(SetBounds(bounds))
	})
	RegisterReferences("Grid", BootstrapGridReferences)
}

func BootstrapGridReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
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
	}, []string{"Cols", "Fixed", "Bounded", "Advanced"}
}

func ExampleNewGrid(ops ...func(any)) image.Image {
	// Example 1: Simple 2x2 grid with Gophers
	// Shrink the Gopher so it fits better
	gopher := NewScale(NewGopher(), ScaleFactor(0.25))

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
	gopher := NewScale(NewGopher(), ScaleFactor(0.25))

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
	gopher := NewScale(NewGopher(), ScaleFactor(0.25))

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
	gopher := NewScale(NewGopher(), ScaleFactor(0.25))

	// Create a bounded version that claims to be larger
	zero := 0
	fiveHundred := 500

	b := &boundedGopher{
		Image: gopher,
		bounds: Bounds{
			Left:   &Range{Low: &zero, High: &zero},
			Right:  &Range{Low: &fiveHundred, High: &fiveHundred},
			Top:    &Range{Low: &zero, High: &zero},
			Bottom: &Range{Low: &fiveHundred, High: &fiveHundred},
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
	smallGopher := NewScale(gopher, ScaleFactor(0.2))

	// 2. Cropped Gopher (Head only)
	// Gopher is roughly 500x700. Head is at top.
	head := NewCrop(gopher, image.Rect(100, 50, 400, 350))
	smallHead := NewScale(head, ScaleFactor(0.3))

	// 3. Tiled Gopher
	// Tile a small gopher into a 100x100 box
	tinyGopher := NewScale(gopher, ScaleSize(20, 30))
	tiled := NewTile(tinyGopher, image.Rect(0, 0, 100, 100))

	// 4. Text Label
	txt := NewText("Hello Grid", 24, color.Black, nil)

	// 5. Padding with Checker background
	checker := NewChecker(color.RGBA{200, 200, 200, 255}, color.White)
	padded := NewPadding(smallGopher, 20, checker)

	args := []any{
		Row(Cell(txt), Cell(padded)),
		Row(Cell(smallHead), Cell(tiled)),
	}
	for _, op := range ops {
		args = append(args, op)
	}

	return NewGrid(args...)
}
