package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Gradient implementations implement the image.Image interface.
var _ image.Image = (*LinearGradient)(nil)
var _ image.Image = (*RadialGradient)(nil)
var _ image.Image = (*ConicGradient)(nil)

// LinearGradient represents a linear color gradient.
type LinearGradient struct {
	Null
	StartColor
	EndColor
	Vertical bool
}

// At returns the color at (x, y).
func (g *LinearGradient) At(x, y int) color.Color {
	b := g.Bounds()
	if b.Empty() {
		return color.RGBA{}
	}

	var t float64
	if g.Vertical {
		if b.Dy() <= 1 {
			t = 0
		} else {
			t = float64(y-b.Min.Y) / float64(b.Dy()-1)
		}
	} else {
		if b.Dx() <= 1 {
			t = 0
		} else {
			t = float64(x-b.Min.X) / float64(b.Dx()-1)
		}
	}

	return lerpColor(g.StartColor.StartColor, g.EndColor.EndColor, t)
}

// NewLinearGradient creates a new LinearGradient pattern.
func NewLinearGradient(ops ...func(any)) image.Image {
	g := &LinearGradient{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	g.StartColor.StartColor = color.White
	g.EndColor.EndColor = color.Black

	for _, op := range ops {
		op(g)
	}
	return g
}

// Option to set direction to vertical
func GradientVertical() func(any) {
	return func(i any) {
		if g, ok := i.(*LinearGradient); ok {
			g.Vertical = true
		}
	}
}

// RadialGradient represents a radial color gradient.
type RadialGradient struct {
	Null
	StartColor
	EndColor
	// Center can be added if needed, defaulting to bounds center
}

// At returns the color at (x, y).
func (g *RadialGradient) At(x, y int) color.Color {
	b := g.Bounds()
	if b.Empty() {
		return color.RGBA{}
	}

	cx := float64(b.Min.X + b.Dx()/2)
	cy := float64(b.Min.Y + b.Dy()/2)

	dx := float64(x) - cx
	dy := float64(y) - cy

	dist := math.Sqrt(dx*dx + dy*dy)

	// Max distance is from center to corner (or side?)
	// Usually radial gradient goes to the furthest corner or closest side.
	// CSS radial-gradient defaults to furthest-corner.
	// Let's use half of the smallest dimension (circle fits in box) or distance to corner.
	// Distance to corner:
	maxDist := math.Sqrt(float64(b.Dx()*b.Dx() + b.Dy()*b.Dy())) / 2.0

	t := dist / maxDist

	return lerpColor(g.StartColor.StartColor, g.EndColor.EndColor, t)
}

// NewRadialGradient creates a new RadialGradient pattern.
func NewRadialGradient(ops ...func(any)) image.Image {
	g := &RadialGradient{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	g.StartColor.StartColor = color.White
	g.EndColor.EndColor = color.Black

	for _, op := range ops {
		op(g)
	}
	return g
}

// ConicGradient represents a conic (angular) color gradient.
type ConicGradient struct {
	Null
	StartColor
	EndColor
}

// At returns the color at (x, y).
func (g *ConicGradient) At(x, y int) color.Color {
	b := g.Bounds()
	if b.Empty() {
		return color.RGBA{}
	}

	cx := float64(b.Min.X + b.Dx()/2)
	cy := float64(b.Min.Y + b.Dy()/2)

	dx := float64(x) - cx
	dy := float64(y) - cy

	// Atan2 returns -Pi to Pi
	angle := math.Atan2(dy, dx)

	// Normalize to 0..1
	// Atan2: 0 is right (positive X), Pi/2 is down (positive Y), -Pi/2 is up.
	// We want 0..1.
	// Let's say we start at top (-Pi/2) or right (0).
	// Standard usually starts at top (12 o'clock).
	// atan2(y, x).
	// If we want top to be 0:
	// top: dy=-1, dx=0 -> atan2 = -Pi/2.
	// right: dy=0, dx=1 -> atan2 = 0.
	// bottom: dy=1, dx=0 -> atan2 = Pi/2.
	// left: dy=0, dx=-1 -> atan2 = Pi.

	// Map -Pi..Pi to 0..1
	// t = (angle + Pi) / (2 * Pi) maps (-Pi -> 0, Pi -> 1).
	// This makes -Pi (Left) start.

	// If we want to rotate start, we can add offset.
	// For now, simple mapping.

	t := (angle + math.Pi) / (2 * math.Pi)

	return lerpColor(g.StartColor.StartColor, g.EndColor.EndColor, t)
}

// NewConicGradient creates a new ConicGradient pattern.
func NewConicGradient(ops ...func(any)) image.Image {
	g := &ConicGradient{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	g.StartColor.StartColor = color.White
	g.EndColor.EndColor = color.Black

	for _, op := range ops {
		op(g)
	}
	return g
}

// lerpColor interpolates between two colors.
func lerpColor(c1, c2 color.Color, t float64) color.Color {
	if t <= 0 {
		return c1
	}
	if t >= 1 {
		return c2
	}

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	// RGBA() returns alpha-premultiplied values in range [0, 0xffff].
	// Interpolate linearly.

	r := uint16(float64(r1) + t*(float64(r2)-float64(r1)))
	g := uint16(float64(g1) + t*(float64(g2)-float64(g1)))
	b := uint16(float64(b1) + t*(float64(b2)-float64(b1)))
	a := uint16(float64(a1) + t*(float64(a2)-float64(a1)))

	return color.RGBA64{R: r, G: g, B: b, A: a}
}
