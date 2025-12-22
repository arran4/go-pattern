package pattern

import (
	"image"
	"image/color"
)

// Range defines a range with optional low and high bounds.
// If Low or High is nil, it represents an unbounded range in that direction.
type Range struct {
	Low, High *int
}

// Bounds defines the boundaries of an object using Ranges.
type Bounds struct {
	Left, Right, Top, Bottom *Range
}

// Bounded is the interface that bounded objects should implement.
// Note: The method is named PatternBounds to avoid conflict with image.Image.Bounds().
type Bounded interface {
	PatternBounds() Bounds
}

// Ensure Grid implements image.Image
var _ image.Image = (*Grid)(nil)

type Grid struct {
	Null
	rows        map[int]map[int]image.Image
	cols        int
	rowsCount   int
	cellWidths  []int
	rowHeights  []int
	fixedWidth  int
	fixedHeight int
}

func (g *Grid) ColorModel() color.Model {
	return color.RGBAModel
}

func (g *Grid) Bounds() image.Rectangle {
	return g.bounds
}

func (g *Grid) At(x, y int) color.Color {
	// Simple lookup logic (to be refined)
	// We need to map x, y to a specific cell

	// If layout is not calculated, we might need to do it.
	// But let's assume bounds and cell sizes are calculated at creation or SetBounds.

	// Find which column x belongs to
	colIdx := -1
	currentX := g.bounds.Min.X
	for i, w := range g.cellWidths {
		if x >= currentX && x < currentX+w {
			colIdx = i
			break
		}
		currentX += w
	}
	if colIdx == -1 {
		return color.RGBA{}
	}

	// Find which row y belongs to
	rowIdx := -1
	currentY := g.bounds.Min.Y
	for i, h := range g.rowHeights {
		if y >= currentY && y < currentY+h {
			rowIdx = i
			break
		}
		currentY += h
	}
	if rowIdx == -1 {
		return color.RGBA{}
	}

	if row, ok := g.rows[rowIdx]; ok {
		if img, ok := row[colIdx]; ok {
			// Calculate local coordinates for the image
			// We need to know where the cell starts.
			cellX := g.bounds.Min.X
			for i := 0; i < colIdx; i++ {
				cellX += g.cellWidths[i]
			}
			cellY := g.bounds.Min.Y
			for i := 0; i < rowIdx; i++ {
				cellY += g.rowHeights[i]
			}

			// Map (x, y) to image's coordinate space.
			// The image inside the cell is drawn at (cellX, cellY) in the grid's space.
			// But the image itself might have its own bounds (e.g. Min.X != 0).
			// If we assume the image is placed at (cellX, cellY), we usually translate.
			// Or do we assume the image fills the cell?
			// The user said "figures out if contents is bound, and intelligently tries to balance".

			// Let's assume we translate grid (x,y) to image local (lx, ly).
			// If we align top-left of image to top-left of cell:
			// lx = img.Bounds().Min.X + (x - cellX)
			// ly = img.Bounds().Min.Y + (y - cellY)

			lx := img.Bounds().Min.X + (x - cellX)
			ly := img.Bounds().Min.Y + (y - cellY)

			if (image.Point{lx, ly}.In(img.Bounds())) {
				return img.At(lx, ly)
			}
		}
	}

	return color.RGBA{}
}

// Option types

type rowOp struct {
	cells []any
}

type colOp struct {
	cells []any
}

type cellOp struct {
	content any
}

type cellPosOp struct {
	x, y    int
	content any
}

type gridSizeOp struct {
	cols, rows int
}

type fixedSizeOp struct {
	w, h int
}

// Helpers

func Row(cells ...any) any {
	return rowOp{cells: cells}
}

func Column(cells ...any) any {
	return colOp{cells: cells}
}

func Cell(content any) any {
	return cellOp{content: content}
}

func CellPos(x, y int, content any) any {
	return cellPosOp{x: x, y: y, content: content}
}

func GridSize(cols, rows int) any {
	return gridSizeOp{cols: cols, rows: rows}
}

func FixedSize(w, h int) any {
	return fixedSizeOp{w: w, h: h}
}

func NewGrid(ops ...any) image.Image {
	g := &Grid{
		Null: Null{
			bounds: image.Rect(0, 0, 0, 0),
		},
		rows: make(map[int]map[int]image.Image),
	}

	// Process ops
	// First pass: gather cells and constraints
	var currentRow, currentCol int

	for _, op := range ops {
		switch o := op.(type) {
		case rowOp:
			// Add cells to current row
			// If we were in column mode, this might be tricky, but let's assume standard flow
			currentCol = 0
			for _, c := range o.cells {
				g.addCell(currentRow, currentCol, c)
				currentCol++
			}
			currentRow++
		case colOp:
			// Add cells to current column?
			// Usually "Column" means a vertical stack.
			// If we mix Row and Column, it gets complex.
			// Let's assume Column(...) adds a column at current position?
			// Or maybe NewGrid(Column(...), Column(...)) means Grid of columns.

			// If we are at (0,0) and receive Column, we fill (0,0), (0,1), (0,2)...
			// Then next op starts at (1,0)?

			// Determine which column we are adding to.
			// If we use Column(...), we append a new column.

			targetCol := g.cols

			r := 0 // Start from top
			for _, c := range o.cells {
				g.addCell(r, targetCol, c) // Add to new column
				r++
			}
			// addCell updates g.cols if targetCol+1 > g.cols.
			// So after the loop, g.cols should be targetCol+1.
			// We don't need to manually increment g.cols unless addCell didn't do it (e.g. empty column).
			if len(o.cells) == 0 {
				g.cols = targetCol + 1
			}

		case cellOp:
			g.addCell(currentRow, currentCol, o.content)
			currentCol++
		case cellPosOp:
			g.addCell(o.y, o.x, o.content)
		case gridSizeOp:
			if g.cols < o.cols {
				g.cols = o.cols
			}
			if g.rowsCount < o.rows {
				g.rowsCount = o.rows
			}
		case fixedSizeOp:
			g.fixedWidth = o.w
			g.fixedHeight = o.h
		case func(any):
			o(g)
		}
	}

	g.layout()
	return g
}

func (g *Grid) addCell(row, col int, content any) {
	// Unwrap content if it's a CellOp
	if c, ok := content.(cellOp); ok {
		content = c.content
	}

	img, ok := content.(image.Image)
	if !ok {
		// Try to handle non-image content if needed, or skip
		return
	}

	if g.rows[row] == nil {
		g.rows[row] = make(map[int]image.Image)
	}
	g.rows[row][col] = img

	if row+1 > g.rowsCount {
		g.rowsCount = row + 1
	}
	if col+1 > g.cols {
		g.cols = col + 1
	}
}

func (g *Grid) layout() {
	// Determine cell sizes
	// 1. Calculate max width for each column and max height for each row based on content bounds
	// This only works for bounded contents.

	colWidths := make([]int, g.cols)
	rowHeights := make([]int, g.rowsCount)

	for r := 0; r < g.rowsCount; r++ {
		for c := 0; c < g.cols; c++ {
			if img, ok := g.rows[r][c]; ok {
				// Check for Intrinsic Bounded interface first (User's Bounded)
				if b, ok := img.(Bounded); ok {
					pb := b.PatternBounds()
					// Check if Right and Left are set (Bounded width)
					if pb.Right != nil && pb.Right.High != nil && pb.Left != nil && pb.Left.Low != nil {
						w := *pb.Right.High - *pb.Left.Low // Simplified
						if w > colWidths[c] {
							colWidths[c] = w
						}
					} else {
						// Fallback to image.Image.Bounds()
						w := img.Bounds().Dx()
						if w > colWidths[c] {
							colWidths[c] = w
						}
					}

					// Similar for height
					if pb.Bottom != nil && pb.Bottom.High != nil && pb.Top != nil && pb.Top.Low != nil {
						h := *pb.Bottom.High - *pb.Top.Low
						if h > rowHeights[r] {
							rowHeights[r] = h
						}
					} else {
						h := img.Bounds().Dy()
						if h > rowHeights[r] {
							rowHeights[r] = h
						}
					}

				} else {
					// Standard image.Image
					w := img.Bounds().Dx()
					h := img.Bounds().Dy()

					if w > colWidths[c] {
						colWidths[c] = w
					}
					if h > rowHeights[r] {
						rowHeights[r] = h
					}
				}
			}
		}
	}

	g.cellWidths = colWidths
	g.rowHeights = rowHeights

	totalW := 0
	for _, w := range colWidths {
		totalW += w
	}
	totalH := 0
	for _, h := range rowHeights {
		totalH += h
	}

	// If FixedSize is set, we might need to adjust or center?
	// For now, let's set the bounds to the calculated total size
	if g.fixedWidth > 0 && g.fixedHeight > 0 {
		g.bounds = image.Rect(0, 0, g.fixedWidth, g.fixedHeight)
		// If content is smaller, we might want to distribute space?
		// "Table balancing formula that is a simplified one webbrowsers use"
		// If content exceeds fixed size, we clip?
	} else {
		g.bounds = image.Rect(0, 0, totalW, totalH)
	}
}
