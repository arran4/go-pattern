package pattern

import (
	"image"
	"image/color"
	"testing"
)

// mockBounded implements Bounded and image.Image
type mockBounded struct {
	bounds Bounds
	rect   image.Rectangle
}

func (m *mockBounded) ColorModel() color.Model { return color.RGBAModel }
func (m *mockBounded) Bounds() image.Rectangle { return m.rect }
func (m *mockBounded) At(x, y int) color.Color { return color.RGBA{255, 0, 0, 255} }
func (m *mockBounded) PatternBounds() Bounds   { return m.bounds }

func TestGrid_Layout_Basic(t *testing.T) {
	// Create a simple 2x2 grid with 10x10 images
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	g := NewGrid(
		Row(Cell(img), Cell(img)),
		Row(Cell(img), Cell(img)),
	)

	if g.Bounds().Dx() != 20 {
		t.Errorf("Expected width 20, got %d", g.Bounds().Dx())
	}
	if g.Bounds().Dy() != 20 {
		t.Errorf("Expected height 20, got %d", g.Bounds().Dy())
	}
}

func TestGrid_Layout_Bounded(t *testing.T) {
	// Create a mock bounded object that reports 100x100 via PatternBounds,
	// but 10x10 via image.Bounds()

	zero := 0
	hundred := 100

	mb := &mockBounded{
		rect: image.Rect(0, 0, 10, 10),
		bounds: Bounds{
			Left:   &Range{Low: &zero, High: &zero},
			Right:  &Range{Low: &hundred, High: &hundred},
			Top:    &Range{Low: &zero, High: &zero},
			Bottom: &Range{Low: &hundred, High: &hundred},
		},
	}

	g := NewGrid(
		Row(Cell(mb)),
	)

	if g.Bounds().Dx() != 100 {
		t.Errorf("Expected width 100, got %d", g.Bounds().Dx())
	}
	if g.Bounds().Dy() != 100 {
		t.Errorf("Expected height 100, got %d", g.Bounds().Dy())
	}
}

func TestGrid_FixedSize(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	g := NewGrid(
		FixedSize(50, 50),
		Row(Cell(img)),
	)

	if g.Bounds().Dx() != 50 {
		t.Errorf("Expected width 50, got %d", g.Bounds().Dx())
	}
	if g.Bounds().Dy() != 50 {
		t.Errorf("Expected height 50, got %d", g.Bounds().Dy())
	}
}

func TestGrid_Column(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	// Two columns side by side
	g := NewGrid(
		Column(Cell(img), Cell(img)), // Col 1: 2 images stacked -> 10w, 20h
		Column(Cell(img)),            // Col 2: 1 image -> 10w, 10h
	)

	// Total width should be 20, total height should be max(20, 10) = 20
	if g.Bounds().Dx() != 20 {
		t.Errorf("Expected width 20, got %d", g.Bounds().Dx())
	}
	if g.Bounds().Dy() != 20 {
		t.Errorf("Expected height 20, got %d", g.Bounds().Dy())
	}
}

func TestGrid_GridSize(t *testing.T) {
	// Create a 5x5 grid but only populate one cell
	// It should be 0 size if no content, but since we don't know content size,
	// GridSize mainly affects the number of cells reserved.
	// However, if we don't have content, the size is 0.
	// But if we have content at (0,0) and GridSize(5,5), the grid should be at least (0,0) to (5*w, 5*h)?
	// My implementation only uses existing content to determine row/col sizes.
	// If a row/col is empty, its size is 0.

	// So GridSize(5,5) with one cell at (0,0) (10x10) will result in:
	// Col 0: 10w. Cols 1-4: 0w. Total w: 10.
	// Row 0: 10h. Rows 1-4: 0h. Total h: 10.

	// Wait, if I supply GridSize, does it mean empty cells have some default size?
	// The prompt says: "NewGrid(GridSize(2,2), CellPos(0,1,i1), CellPos(1,0,i3))".
	// This implies GridSize sets up the structure.

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	g := NewGrid(
		GridSize(5, 5),
		CellPos(0, 0, img),
	)

	// My current implementation: cols=5, rows=5.
	// cellWidths[0] = 10. cellWidths[1..4] = 0.
	// rowHeights[0] = 10. rowHeights[1..4] = 0.
	// Bounds: 10x10.

	// This seems correct for "flexible" layout where empty columns collapse.
	// Unless "Table balancing formula" implies uniform distribution?
	// "intelligently tries to balance the internals using table balancing formula that is a simplified one webbrowsers use"

	if g.Bounds().Dx() != 10 {
		t.Errorf("Expected width 10, got %d", g.Bounds().Dx())
	}
}
