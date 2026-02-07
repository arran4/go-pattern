package pattern

import (
	"image"
	"testing"
)

func BenchmarkGrid_At(b *testing.B) {
	// Setup a large grid manually to avoid overhead of NewGrid with many ops
	cols := 2000
	rows := 2000
	cellW := 10
	cellH := 10

	g := &Grid{
		bounds:     image.Rect(0, 0, cols*cellW, rows*cellH),
		rows:       make(map[int]map[int]image.Image),
		cols:       cols,
		rowsCount:  rows,
		cellWidths: make([]int, cols),
		rowHeights: make([]int, rows),
	}

	// Compute offsets manually for benchmark since layout() isn't called
	g.colOffsets = make([]int, cols+1)
	current := 0
	for i := range g.cellWidths {
		g.cellWidths[i] = cellW
		g.colOffsets[i] = current
		current += cellW
	}
	g.colOffsets[cols] = current

	g.rowOffsets = make([]int, rows+1)
	current = 0
	for i := range g.rowHeights {
		g.rowHeights[i] = cellH
		g.rowOffsets[i] = current
		current += cellH
	}
	g.rowOffsets[rows] = current

	// We don't populate g.rows because At() performs the search BEFORE looking up the image.
	// The search is what we are optimizing.

	// Test point near the end to hit worst-case scenario for linear search
	x := cols*cellW - 5
	y := rows*cellH - 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.At(x, y)
	}
}
