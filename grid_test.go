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
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	g := NewGrid(
		GridSize(5, 5),
		CellPos(0, 0, img),
	)

	if g.Bounds().Dx() != 10 {
		t.Errorf("Expected width 10, got %d", g.Bounds().Dx())
	}
}

func TestGrid_NonZeroOrigin(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// Set a pixel to verify coordinate mapping: Local (5,5)
	img.Set(5, 5, color.RGBA{255, 0, 0, 255})
	// Transparent everywhere else by default

	// Case 1: Fixed Non-zero explicitly
	g1 := NewGrid(
		CellPos(0, 0, img),
		SetBounds(image.Rect(10, 10, 20, 20)),
		FixedSize(10, 10), // Required to hit the first layout branch
	)

	gb1 := g1.Bounds()
	if gb1.Min.X != 10 || gb1.Min.Y != 10 || gb1.Max.X != 20 || gb1.Max.Y != 20 {
		t.Fatalf("Grid bounds do not respect explicit SetBounds (FixedSize branch): %v", gb1)
	}

	c1 := g1.At(15, 15)
	r1, g1c, b1, a1 := c1.RGBA()
	if r1 != 0xFFFF || g1c != 0 || b1 != 0 || a1 != 0xFFFF {
		t.Errorf("Expected pixel at (15, 15) to map to local (5, 5) and be fully red. Got color: %v", c1)
	}

	// Case 2: Natural Non-zero explicitly
	// If FixedSize isn't called, SetBounds should still inform the Min bounds of the layout
	g2 := NewGrid(
		CellPos(0, 0, img),
		SetBounds(image.Rect(10, 10, 20, 20)),
	)

	gb2 := g2.Bounds()
	// It takes Min from bounds, but Max from totalW/H.
	// Total W/H is 10/10 from the image.
	// So expected is (10,10)-(20,20).
	if gb2.Min.X != 10 || gb2.Min.Y != 10 || gb2.Max.X != 20 || gb2.Max.Y != 20 {
		t.Fatalf("Grid bounds do not respect explicit SetBounds (Natural branch): %v", gb2)
	}

	c2 := g2.At(15, 15)
	r2, g2c, b2, a2 := c2.RGBA()
	if r2 != 0xFFFF || g2c != 0 || b2 != 0 || a2 != 0xFFFF {
		t.Errorf("Expected pixel at (15, 15) to map to local (5, 5) and be fully red. Got color: %v", c2)
	}

	// Case 3: Negative Origin
	g3 := NewGrid(
		CellPos(0, 0, img),
		SetBounds(image.Rect(-10, -10, 0, 0)),
	)

	gb3 := g3.Bounds()
	if gb3.Min.X != -10 || gb3.Min.Y != -10 || gb3.Max.X != 0 || gb3.Max.Y != 0 {
		t.Fatalf("Grid bounds do not respect negative origin explicit SetBounds (Natural branch): %v", gb3)
	}

	c3 := g3.At(-5, -5)
	r3, g3c, b3, a3 := c3.RGBA()
	if r3 != 0xFFFF || g3c != 0 || b3 != 0 || a3 != 0xFFFF {
		t.Errorf("Expected pixel at (-5, -5) to map to local (5, 5) and be fully red. Got color: %v", c3)
	}

	// Case 4: Default Zero-origin behavior unaffected
	g4 := NewGrid(
		CellPos(0, 0, img),
	)

	gb4 := g4.Bounds()
	if gb4.Min.X != 0 || gb4.Min.Y != 0 || gb4.Max.X != 10 || gb4.Max.Y != 10 {
		t.Fatalf("Grid bounds default behavior altered: %v", gb4)
	}
}
