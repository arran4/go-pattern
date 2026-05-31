package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestGrid_FlexibleLayout(t *testing.T) {
	// Test flexible column expansion
	// Grid FixedSize: 200x100
	// Col 0: Bounded 50x100
	// Col 1: Unbounded (0 width initially)
	// Expected: Col 0 = 50, Col 1 = 150.

	zero := 0
	fifty := 50
	hundred := 100

	b := &mockBounded{
		rect: image.Rect(0, 0, 10, 10),
		bounds: Bounds{
			Left:   &Range{Low: &zero, High: &zero},
			Right:  &Range{Low: &fifty, High: &fifty},
			Top:    &Range{Low: &zero, High: &zero},
			Bottom: &Range{Low: &hundred, High: &hundred},
		},
	}

	// Unbounded pattern (returns nil bounds via PatternBounds if checked, or simply not Bounded)
	// Here we use a struct that implements Bounded with nil ranges to ensure it has "0" calculated size.
	u := &mockUnboundedPattern{}

	g := NewGrid(
		FixedSize(200, 100),
		Row(Cell(b), Cell(u)),
	).(*Grid)

	// Check calculated cell widths
	if len(g.cellWidths) != 2 {
		t.Fatalf("Expected 2 columns, got %d", len(g.cellWidths))
	}

	if g.cellWidths[0] != 50 {
		t.Errorf("Expected col 0 width 50, got %d", g.cellWidths[0])
	}

	if g.cellWidths[1] != 150 {
		t.Errorf("Expected col 1 width 150, got %d", g.cellWidths[1])
	}
}

type mockUnboundedPattern struct {
	image.Image
}

func (u *mockUnboundedPattern) PatternBounds() Bounds {
	return Bounds{Left: nil, Right: nil, Top: nil, Bottom: nil}
}

func (u *mockUnboundedPattern) ColorModel() color.Model { return color.RGBAModel }
func (u *mockUnboundedPattern) Bounds() image.Rectangle { return image.Rect(0,0,1,1) }
func (u *mockUnboundedPattern) At(x, y int) color.Color { return color.RGBA{} }

func TestGrid_FlexibleLayout_Rows(t *testing.T) {
	// Test flexible row expansion
	// Grid FixedSize: 100x200
	// Row 0: Bounded 100x50
	// Row 1: Unbounded
	// Expected: Row 0 = 50, Row 1 = 150.

	zero := 0
	fifty := 50
	hundred := 100

	b := &mockBounded{
		rect: image.Rect(0, 0, 10, 10),
		bounds: Bounds{
			Left:   &Range{Low: &zero, High: &zero},
			Right:  &Range{Low: &hundred, High: &hundred},
			Top:    &Range{Low: &zero, High: &zero},
			Bottom: &Range{Low: &fifty, High: &fifty},
		},
	}

	u := &mockUnboundedPattern{}

	g := NewGrid(
		FixedSize(100, 200),
		Row(Cell(b)),
		Row(Cell(u)),
	).(*Grid)

	if len(g.rowHeights) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(g.rowHeights))
	}

	if g.rowHeights[0] != 50 {
		t.Errorf("Expected row 0 height 50, got %d", g.rowHeights[0])
	}

	if g.rowHeights[1] != 150 {
		t.Errorf("Expected row 1 height 150, got %d", g.rowHeights[1])
	}
}
