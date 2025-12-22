package pattern

import (
	"image"
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
	}, []string{"Cols", "Fixed", "Bounded"}
}

func ExampleNewGrid(ops ...func(any)) image.Image {
	// Example 1: Simple 2x2 grid with Gophers
	gopher := NewGopher()

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
	gopher := NewGopher()

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
	gopher := NewGopher()

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
	gopher := NewGopher()

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
