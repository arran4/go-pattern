package pattern

import (
	"image"
	"image/color"
	"math"
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

	if h == 0 || s == 0 {
		return color.RGBA{}
	}

	v := dy / h
	u := dx/s + v/2.0

	// Check if point is outside the main triangle
	// Logical Triangle is 0 <= u <= v <= 1 ?
	if v < 0 || v > 1 || u < 0 || u > v {
		if p.SpaceColor.SpaceColor != nil {
			return p.SpaceColor.SpaceColor
		}
		return color.RGBA{}
	}

	// Iterative Sierpinski Check
	// Determine recursion depth based on triangle size in pixels.
	// We want to stop when the sub-triangle is approximately 1 pixel in size.
	// Initial size is s. At depth k, size is s / 2^k.
	// We stop when s / 2^k < 0.5 (sub-pixel).
	// 2^k > 2s => k > log2(2s) = log2(s) + 1.
	// We add a small buffer for sharpness.

	var maxIterations int
	if s > 1 {
		maxIterations = int(math.Log2(s)) + 2
	} else {
		maxIterations = 1
	}

	// Cap iterations to prevent excessive computation for huge images
	if maxIterations > 20 {
		maxIterations = 20
	}
	// And ensure at least some iterations
	if maxIterations < 1 {
		maxIterations = 1
	}

	for i := 0; i < maxIterations; i++ {
		// Check for Hole
		// Hole is defined by: v >= 0.5 AND u <= 0.5 AND (v - u) < 0.5

		isBelowMid := v >= 0.5
		isLeftMid := u < 0.5

		if isBelowMid && isLeftMid {
			// Possibly in hole or BL triangle.
			// Hole if v - u < 0.5.
			// Vertices of Hole: (0, 0.5), (0.5, 1), (0.5, 0.5).

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
