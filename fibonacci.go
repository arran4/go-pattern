package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Fibonacci implements the image.Image interface.
var _ image.Image = (*Fibonacci)(nil)

// Fibonacci is a pattern that draws a Fibonacci (Golden) spiral.
// It uses the logarithmic spiral equation r = a * e^(b * theta) with b = 2*ln(Phi)/pi.
// It supports LineSize, LineColor, and SpaceColor.
type Fibonacci struct {
	Null
	LineSize
	LineColor
	SpaceColor
}

const (
	// Phi is the Golden Ratio.
	Phi = 1.618033988749895
	// GoldenSpiralB is the growth factor for the Golden Spiral (2 * ln(Phi) / pi).
	GoldenSpiralB = 0.30634896253
)

func (p *Fibonacci) At(x, y int) color.Color {
	b := p.Bounds()

	// Center coordinate.
	// Using float for precision.
	cx := float64(b.Min.X+b.Max.X) / 2.0
	cy := float64(b.Min.Y+b.Max.Y) / 2.0

	// Offset from center.
	// Note: y axis increases downwards in images.
	// Standard math assumes y increases upwards.
	// However, spiral shape is symmetric under reflection except for chirality.
	// Let's stick to image coordinates for simplicity.
	dx := float64(x) - cx
	dy := float64(y) - cy

	r := math.Hypot(dx, dy)
	theta := math.Atan2(dy, dx) // Returns (-Pi, Pi]

	// Logarithmic Spiral: r = a * e^(b * theta)
	// We want to check if pixel is close to ANY arm of the spiral.
	// The arms correspond to theta + 2*pi*k.
	// r = a * e^(b * (theta + 2*pi*k))
	// ln(r) = ln(a) + b * theta + b * 2*pi*k
	// (ln(r) - ln(a)) / b - theta = 2*pi*k

	// We set a = 1.0 implicitly for the base scale.
	// Scale adjustments can be done via Zoom or by the user provided bounds scaling if we supported it.
	// For now, fixed a=1.
	a := 1.0

	// Handle center singularity
	if r < 1e-6 {
		// At the exact center.
		// If line size is non-zero, center is part of the line.
		if p.LineSize.LineSize > 0 {
			return p.LineColor.LineColor
		}
		if p.SpaceColor.SpaceColor != nil {
			return p.SpaceColor.SpaceColor
		}
		return color.RGBA{}
	}

	val := (math.Log(r/a))/GoldenSpiralB - theta

	// val should be close to 2*pi*k
	period := 2 * math.Pi

	// Normalize val to [0, 2*pi) to find the fractional part relative to period?
	// No, we just want distance to nearest multiple of 2*pi.

	k := math.Round(val / period)
	diff := val - k*period

	// Angular distance
	angularDist := math.Abs(diff)

	// Approximate Euclidean distance to the curve.
	// The normal vector direction has slope -1/b relative to the radial vector.
	// The pitch angle alpha satisfies tan(alpha) = 1/b.
	// The distance perpendicular to the spiral arm is approximately:
	// dist = r * angularDist * sin(alpha)
	// sin(alpha) = 1 / sqrt(1 + b^2) is INCORRECT.
	// tan(alpha) = b (growth rate). No, r' = b*r.
	// Let's rely on gradient magnitude.
	// f(r, theta) = ln(r) - b*theta. Gradient del f = (1/r, -b/r).
	// |del f| = sqrt(1+b^2)/r.
	// dist = |f_val| / |del f| = |angular_dist * b| / (sqrt(1+b^2)/r) = r * angular_dist * b / sqrt(1+b^2).

	factor := GoldenSpiralB / math.Sqrt(1+GoldenSpiralB*GoldenSpiralB)
	dist := r * angularDist * factor

	if dist <= float64(p.LineSize.LineSize)/2.0 {
		return p.LineColor.LineColor
	}

	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
	}
	return color.RGBA{}
}

// NewFibonacci creates a new Fibonacci pattern.
func NewFibonacci(ops ...func(any)) image.Image {
	p := &Fibonacci{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.LineSize.LineSize = 1
	p.LineColor.LineColor = color.Black
	// SpaceColor defaults to nil (transparent)

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoFibonacci produces a demo variant for readme.md pre-populated values
func NewDemoFibonacci(ops ...func(any)) image.Image {
	return NewFibonacci(ops...)
}
