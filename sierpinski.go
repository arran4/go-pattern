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

	// Centered triangle coordinates relative to cx, cy:
	// Top Vertex: (0, -h/2)
	// Bottom Left: (-s/2, h/2)
	// Bottom Right: (s/2, h/2)

	// Translate pixel (x,y) to be relative to the Top Vertex (which will map to 0,0).
	// Top Vertex in image coordinates:
	topX := cx
	topY := cy - h/2

	dx := float64(x) - topX
	dy := float64(y) - topY

	// We map the equilateral triangle to the Pascal Triangle logical space.
	// Vertices map as follows:
	// Top (0, 0) -> Logical (0, 0)
	// Bottom Left (-s/2, h) -> Logical (0, S) (Left Edge of Pascal triangle)
	// Bottom Right (s/2, h) -> Logical (S, S) (Hypotenuse of Pascal triangle)

	// Note: We use a fixed large power of 2 for S to ensure high precision for the bitwise check.
	// S = 2^30
	const S = 1 << 30

	// Transform derivation:
	// lx = (S/s) * dx + (S/2h) * dy = (S/s) * (dx + dy/sqrt(3))
	// ly = (S/h) * dy = (S/s) * (2 * dy / sqrt(3))

	const invSqrt3 = 0.57735026919
	scale := float64(S) / s

	lx := scale * (dx + dy*invSqrt3)
	ly := scale * (2 * dy * invSqrt3)

	// Convert to integers for bitwise check.
	// We verify if the point is inside the triangle first.
	// In logical space (Pascal), the triangle is bounded by lx >= 0, ly >= lx, ly <= S.
	// Wait, x goes from 0 to y. So 0 <= lx <= ly.
	// And ly goes from 0 to S.

	ix := int(lx)
	iy := int(ly)

	// Check bounds
	if ix < 0 || iy > S || ix > iy {
		// Outside the triangle
		if p.SpaceColor.SpaceColor != nil {
			return p.SpaceColor.SpaceColor
		}
		return color.RGBA{}
	}

	// Check Sierpiński condition (Pascal's Triangle Mod 2)
	// Condition: binomial(y, x) is odd iff (y & x) == x
	// Here lx corresponds to 'k' (column), ly corresponds to 'n' (row).
	if (iy & ix) == ix {
		if p.FillColor.FillColor != nil {
			return p.FillColor.FillColor
		}
		return color.RGBA{}
	}

	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
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
