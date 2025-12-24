package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestSierpinski_At(t *testing.T) {
	// Create a 100x100 Sierpinski triangle
	s := NewSierpinski(
		SetBounds(image.Rect(0, 0, 100, 100)),
		SetFillColor(color.Black),
		SetSpaceColor(color.White),
	)

	// Test bounds
	if s.Bounds().Dx() != 100 || s.Bounds().Dy() != 100 {
		t.Errorf("Bounds size mismatch: %v", s.Bounds())
	}

	// Calculate center
	// cx := 50
	// cy := 50
	// Height of triangle is roughly 100 * sqrt(3)/2 = 86.6
	// h := 86.6
	// Top vertex at (50, 50 - 43.3) = (50, 6.7)

	// Check a point on the LEFT EDGE.
	// Left edge corresponds to logical x=0.
	// Condition (y & 0) == 0 is always true.
	// So Left Edge should be solid.

	// Left Edge equation: x = cx - (y - topY) / sqrt(3)
	// At y=50.
	// topY = 6.7.
	// y - topY = 43.3.
	// (y - topY) / sqrt(3) = 43.3 / 1.732 = 25.
	// x = 50 - 25 = 25.

	// Sample (25, 50).
	if s.At(25, 50) != color.Black {
		t.Errorf("Left edge point (25,50) should be Black, got %v", s.At(25, 50))
	}

	// Verify Outside behavior.
	// (0, 0) should be outside (White).
	if s.At(0, 0) != color.White {
		t.Errorf("Corner (0,0) should be White, got %v", s.At(0, 0))
	}
	// (0, 50). x=0. left edge x=25.
	// 0 < 25. So outside on the left.
	if s.At(0, 50) != color.White {
		t.Errorf("Point (0,50) should be White, got %v", s.At(0, 50))
	}

	// Verify a point in the main central hole.
	// Center hole centroid: (50, 64.4).
	// Sample (50, 65).
	// lx approx ly/2.
	// If ly is large (near 2^29), and centered vertically?
	// The main hole corresponds to the middle inverted triangle in standard Sierpinski.
	// In Pascal terms, this is the even numbers in the middle.
	// E.g. Row 2^K to 2^(K+1)-1 has odd numbers on edges but evens in middle?
	// Row 3: 1 1 1 1. (Filled)
	// Row 2: 1 0 1. (Hole at x=1).
	// Mapping maps top half of image to top half of Pascal triangle?
	// S is large. So we are at large row numbers.
	// The structure is self-similar.
	// The main central hole corresponds to the region where bits overlap.
	// (y & x) != x ? No, (y & x) == x is Fill.
	// So Hole is (y & x) != x.
	// Central hole:
	// In standard Sierpinski construction (removal of middle triangle):
	// Top (0,0), BL(0,1), BR(1,1). (Right triangle coords).
	// Middle triangle vertices: (0, 0.5), (0.5, 0.5), (0.5, 1)? No.
	// Midpoints: (0, 0.5), (0.5, 1), (0.5, 0.5) ?
	// In my mapping:
	// Top (0,0). BL (0,S). BR (S,S).
	// Midpoints:
	// M_Top_BL = (0, S/2).
	// M_Top_BR = (S/2, S/2).
	// M_BL_BR = (S/2, S).
	// The hole is the triangle bounded by these midpoints.
	// Triangle ((0, S/2), (S/2, S/2), (S/2, S)).
	// Check centroid of this hole:
	// Average x = (0 + S/2 + S/2)/3 = S/3.
	// Average y = (S/2 + S/2 + S)/3 = 2S/3.
	// Point (S/3, 2S/3).
	// Condition (iy & ix) == ix.
	// iy = 2S/3. ix = S/3.
	// iy = 2 * ix.
	// (2x & x) == x ?
	// 2x shifts bits left.
	// If x has any bits, 2x usually has different bits.
	// e.g. x=1 (01). 2x=2 (10). 01 & 10 = 00 != 01.
	// So (iy & ix) == 0 != ix.
	// So it is Space (Hole).

	// Calculate screen coordinates for Logical (S/3, 2S/3).
	// ly = 2S/3 = scale * (2 * dy / sqrt3).
	// S = scale * s.
	// 2/3 * scale * s = scale * 2 * dy / sqrt3.
	// s/3 = dy / sqrt3.
	// dy = s * sqrt3 / 3 = s * 0.577.
	// Recall h = s * sqrt3 / 2.
	// dy = h * 2 / 3.
	// Screen y = topY + 2/3 * h.
	// topY = 6.7. h = 86.6.
	// y = 6.7 + 57.7 = 64.4.
	// Screen x?
	// lx = S/3 = scale * (dx + dy/sqrt3).
	// s/3 = dx + dy/sqrt3.
	// s/3 = dx + s/3.
	// dx = 0.
	// So Screen x = cx = 50.
	// Screen y = 64.4.

	// So point (50, 65) is inside the hole.
	if s.At(50, 65) != color.White {
		t.Errorf("Center hole pixel (50, 65) should be White (Space), got %v", s.At(50, 65))
	}
}
