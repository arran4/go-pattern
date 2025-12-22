package pattern

import (
	"image"
	"image/color"
)

var (
	GridUnboundedOutputFilename = "grid_unbounded.png"
	GridUnboundedZoomLevels     = []int{}
	GridUnboundedOrder          = 7
	GridUnboundedBaseLabel      = "Unbounded"
)

func init() {
	RegisterGenerator("GridUnbounded", func(bounds image.Rectangle) image.Image {
		return ExampleNewGridUnbounded(SetBounds(bounds))
	})
	RegisterReferences("GridUnbounded", BootstrapGridUnboundedReferences)
}

func BootstrapGridUnboundedReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Flexible": func(bounds image.Rectangle) image.Image {
			return NewDemoGridFlexible(SetBounds(bounds))
		},
	}, []string{"Flexible"}
}

// UnboundedPattern mimics a pattern with no bounds (nil ranges).
type unboundedPattern struct {
	image.Image
}

func (u *unboundedPattern) PatternBounds() Bounds {
	return Bounds{
		Left:   nil,
		Right:  nil,
		Top:    nil,
		Bottom: nil,
	}
}

func ExampleNewGridUnbounded(ops ...func(any)) image.Image {
	// 300x100 Grid
	// Col 0: Bounded (100x100)
	// Col 1: Unbounded (Should take remaining 200px)

	// bounded := NewChecker(color.Black, color.White) // Checkers default to 255x255 but here we want fixed?
	// Actually NewChecker returns default bounds.
	// Let's use NewCrop or just standard bounds behavior.
	// But `layout()` uses `image.Bounds()` if not `Bounded`.

	// Let's create a bounded Mock that is 100x100.
	hundred := 100
	zero := 0

	b := &boundedGopher{
		Image: NewScale(NewGopher(), ScaleSize(100, 100)),
		bounds: Bounds{
			Left:   &Range{Low: &zero, High: &zero},
			Right:  &Range{Low: &hundred, High: &hundred},
			Top:    &Range{Low: &zero, High: &zero},
			Bottom: &Range{Low: &hundred, High: &hundred},
		},
	}

	// Unbounded pattern: e.g. a generic Tile or Checker that we want to fill space.
	// NewChecker returns 255x255.
	// Let's wrap it in an unbounded structure.
	u := &unboundedPattern{
		Image: NewChecker(color.RGBA{200, 0, 0, 255}, color.White),
	}

	args := []any{
		FixedSize(300, 100),
		Row(Cell(b), Cell(u)),
	}
	for _, op := range ops {
		args = append(args, op)
	}

	return NewGrid(args...)
}


func NewDemoGridFlexible(ops ...func(any)) image.Image {
	// Demonstrate Row flexibility
	// Fixed Height 200.
	// Row 0: Bounded 50h.
	// Row 1: Unbounded (Should take 150h).

	fifty := 50
	zero := 0
	b := &boundedGopher{
		Image: NewScale(NewGopher(), ScaleSize(50, 50)),
		bounds: Bounds{
			Left:   &Range{Low: &zero, High: &zero},
			Right:  &Range{Low: &fifty, High: &fifty},
			Top:    &Range{Low: &zero, High: &zero},
			Bottom: &Range{Low: &fifty, High: &fifty},
		},
	}

	u := &unboundedPattern{
		Image: NewChecker(color.RGBA{0, 200, 0, 255}, color.White),
	}

	args := []any{
		FixedSize(100, 200),
		Row(Cell(b)), // Row 0
		Row(Cell(u)), // Row 1
	}
	for _, op := range ops {
		args = append(args, op)
	}

	return NewGrid(args...)
}
