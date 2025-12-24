package pattern

import (
	"image"
	"image/color"
)

// Ensure Sierpinski implements the image.Image interface.
var _ image.Image = (*Sierpinski)(nil)

// Sierpinski is a pattern that draws a Sierpiński triangle fitting within its bounds.
// The triangle is equilateral and centered.
type Sierpinski struct {
	Null
	FillColor
	SpaceColor
}

func (p *Sierpinski) At(x, y int) color.Color {
	b := p.Bounds()
	width := float64(b.Dx())
	height := float64(b.Dy())

	// Calculate geometric properties of the maximal equilateral triangle fitting in bounds.
	// Equilateral triangle aspect ratio: height = side * sqrt(3) / 2
	const sin60 = 0.86602540378 // sqrt(3) / 2

	var s, h float64
	if height/width > sin60 {
		// Width is the limiting factor
		s = width
		h = s * sin60
	} else {
		// Height is the limiting factor
		h = height
		s = h / sin60
	}

	// Center of the bounding box
	cx := float64(b.Min.X) + width/2
	cy := float64(b.Min.Y) + height/2

	// Translate pixel (x,y) to be relative to the Top Vertex (which will map to 0,0).
	// Top Vertex in image coordinates:
	topX := cx
	topY := cy - h/2

	dx := float64(x) - topX
	dy := float64(y) - topY

	// We map the equilateral triangle to a normalized Logical Right Triangle (0,0)-(0,1)-(1,1).
	// Transformation logic:
	// We want to map:
	// Top Vertex (0, 0) -> (0, 0)
	// Bottom Left (-s/2, h) -> (0, 1)  (u=0, v=1)
	// Bottom Right (s/2, h) -> (1, 1)  (u=1, v=1)

	// Derivation:
	// v = dy / h
	// u = (dx + s/2 * v) / s  ?
	// Check BL: dy=h -> v=1. dx=-s/2. u = (-s/2 + s/2)/s = 0. Correct.
	// Check BR: dy=h -> v=1. dx=s/2. u = (s/2 + s/2)/s = 1. Correct.
	// Check Top: dy=0 -> v=0. dx=0. u = 0. Correct.

	// Wait, at top u depends on v?
	// dx ranges from -s/2 * (y/h) to s/2 * (y/h) inside the triangle?
	// Left edge: x = topX - (y-topY)/sqrt3. dx = -dy/sqrt3 = -dy * (s/2h).
	// u = (-dy*s/2h + s/2 * dy/h) / s = (-s/2 + s/2) * (dy/h) / s = 0. Correct.

	// Formula:
	// v = dy / h
	// u = (dx + dy / sqrt(3)) / s + 0.5 * (dy/h)? No.
	// u = dx/s + 0.5*dy/h + something?
	// Let's stick to u = (dx + s/2 * v) / s = dx/s + v/2.

	if h == 0 || s == 0 {
		return color.RGBA{}
	}

	v := dy / h
	u := dx/s + v/2.0

	// Check if point is outside the main triangle
	// Logical Triangle is 0 <= u <= v <= 1 ?
	// Check BL (0,1). u=0, v=1. u <= v holds.
	// Check BR (1,1). u=1, v=1. u <= v holds.
	// Check Top (0,0). u=0, v=0.
	// Check Mid-Base (0,h). dx=0. v=1. u=0.5. 0.5 <= 1.
	// Check Left Edge (u=0). 0 <= v <= 1.
	// Check Right Edge (u=v).
	// Right Edge equation: dx = dy/sqrt3.
	// u = (dy/sqrt3)/s + (dy/h)/2 = (dy/s)*(1/sqrt3) + 0.5*(dy/h).
	// h = s * sqrt3 / 2 => s = 2h/sqrt3.
	// u = (dy / (2h/sqrt3)) * (1/sqrt3) + 0.5 * v
	// u = (dy/h) * (sqrt3/2 * 1/sqrt3) + 0.5 * v
	// u = v * 0.5 + 0.5 * v = v.
	// So u = v is the Right Edge.

	// Bounding conditions: v >= 0, v <= 1, u >= 0, u <= v.
	// Note: u <= v is implicitly handled by the iterative check usually, but good to check.

	if v < 0 || v > 1 || u < 0 || u > v {
		if p.SpaceColor.SpaceColor != nil {
			return p.SpaceColor.SpaceColor
		}
		return color.RGBA{}
	}

	// Iterative Sierpinski Check
	// Recursively check if point is in the central hole.
	// Triangle vertices at step k are defined in normalized coords.
	// But we can just transform coordinate space.
	// Initial state: Point (u, v) in Right Triangle (0,0)-(0,1)-(1,1).
	// Hole is triangle with vertices (0, 0.5), (0.5, 0.5), (0.5, 1).
	// Hole Condition:
	// v >= 0.5  (Below horizontal mid-line)
	// u <= 0.5  (Left of vertical mid-line)
	// v - u <= 0.5 (Above diagonal hole edge)

	// Transformations for sub-triangles:
	// Top (T): v < 0.5. Map (0,0)->(0,0), (0,0.5)->(0,1), (0.5,0.5)->(1,1).
	//   u' = 2u, v' = 2v.
	// Bottom Left (BL): u < 0.5, v >= 0.5 (but not hole).
	//   Map (0, 0.5)->(0,0).
	//   u' = 2u, v' = 2(v - 0.5).
	// Bottom Right (BR): u >= 0.5.
	//   Map (0.5, 0.5)->(0,0).
	//   u' = 2(u - 0.5), v' = 2(v - 0.5).

	const MaxIterations = 25

	for i := 0; i < MaxIterations; i++ {
		// Check for Hole
		// Hole is defined by: v >= 0.5 AND u <= 0.5 AND (v - u) <= 0.5
		// Note: v-u corresponds to distance from diagonal u=v.
		// Diagonal of square is v - u = 0.
		// Diagonal of Hole is v - u = 0.5.

		isBelowMid := v >= 0.5
		isLeftMid := u < 0.5

		if isBelowMid && isLeftMid {
			// Possibly in hole or BL triangle.
			// Hole if v - u <= 0.5?
			// Vertices of Hole: (0, 0.5), (0.5, 1), (0.5, 0.5).
			// Check point (0.25, 0.75). v=0.75, u=0.25. v-u=0.5. Boundary.
			// Check (0.25, 0.6). v=0.6, u=0.25. v-u=0.35 <= 0.5. Inside Hole.
			// Check (0.1, 0.9). v=0.9, u=0.1. v-u=0.8 > 0.5. BL Triangle.

			if v-u < 0.5 {
				// Inside Hole
				if p.SpaceColor.SpaceColor != nil {
					return p.SpaceColor.SpaceColor
				}
				return color.RGBA{}
			}

			// In Bottom Left Triangle
			u = 2 * u
			v = 2 * (v - 0.5)
		} else if isBelowMid {
			// Bottom Right Triangle (u >= 0.5)
			u = 2 * (u - 0.5)
			v = 2 * (v - 0.5)
		} else {
			// Top Triangle (v < 0.5)
			u = 2 * u
			v = 2 * v
		}
	}

	// If we survive iterations, we are in the set.
	if p.FillColor.FillColor != nil {
		return p.FillColor.FillColor
	}
	return color.RGBA{}
}

// NewSierpinski creates a new Sierpinski pattern.
func NewSierpinski(ops ...func(any)) image.Image {
	p := &Sierpinski{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.FillColor.FillColor = color.Black
	// SpaceColor defaults to nil (transparent)

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoSierpinski produces a demo variant for readme.md pre-populated values
func NewDemoSierpinski(ops ...func(any)) image.Image {
	return NewSierpinski(ops...)
}
