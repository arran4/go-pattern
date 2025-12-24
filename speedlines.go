package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure SpeedLines implements the image.Image interface.
var _ image.Image = (*SpeedLines)(nil)

// SpeedLinesType defines the type of speed lines (Radial or Linear).
type SpeedLinesType int

const (
	SpeedLinesRadial SpeedLinesType = iota
	SpeedLinesLinear
)

// SpeedLines pattern generates manga-style speed lines.
type SpeedLines struct {
	Null
	Center
	LineColor
	SpaceColor
	SpaceImageSource
	LineImageSource
	MinRadius
	MaxRadius
	Density
	Phase
	Type SpeedLinesType
}

func (p *SpeedLines) At(x, y int) color.Color {
	b := p.Bounds()
	if !b.Empty() {
		// If we are outside bounds, return transparent? Null handles it if embedded correctly?
		// image.Image At() usually handles bounds checking, but we might be called by something else.
		// Standard image At returns 0 for out of bounds.
		if !(image.Point{x, y}).In(b) {
			return color.RGBA{}
		}
	}

	cx := float64(p.CenterX)
	cy := float64(p.CenterY)

	// Default center if not set (0,0 is valid, but usually we want center of image)
	// But `Center` defaults to 0,0.
	// If the user didn't set center, it's 0,0.
	// Let's check if we should default to middle.
	// But Center struct just holds ints. We don't know if it was set.
	// We'll assume the user sets it, or it is top-left.
	// Wait, for radial burst, 0,0 is weird if image is 100x100.
	// I'll stick to 0,0 default as per struct zero value.

	dx := float64(x) - cx
	dy := float64(y) - cy

	var angle float64
	var dist float64

	if p.Type == SpeedLinesRadial {
		dist = math.Sqrt(dx*dx + dy*dy)
		angle = math.Atan2(dy, dx)
	} else {
		// Linear: Angle determines direction.
		// Use Phase as angle? Or add Angle field?
		// For now, let's say Linear is Horizontal.
		// Or utilize `p.Phase` as angle? No, Phase is for animation.
		// Let's assume Linear means "Horizontal" (moving left/right) for now, or use rotation pattern.
		// But prompt said "radial or linear".
		// I'll implement "Linear" as rays coming from one side (infinity).
		// Basically parallel lines.
		// Let's assume vertical lines (rain) or horizontal.
		// I'll use `angle = atan2(dy, dx)` for Radial.
		// For Linear, we just use one coordinate.
		// Let's assume Linear is horizontal stripes (like speed).
		// We can rotate it using `Rotate`.
		dist = dx // Distance along the line
		angle = dy / 100.0 // "Angle" maps to the cross dimension.
	}

	// Algorithm for Lines:
	// We want random lines based on angle.
	// We map angle to a seed.

	// Normalize angle to 0..1 (or larger range for resolution)
	// Atan2 is -Pi to Pi.
	normalizedAngle := angle
	if p.Type == SpeedLinesRadial {
		normalizedAngle = (angle + math.Pi) / (2 * math.Pi)
	}

	// Density scales the "frequency" of lines.
	// We want discrete "sectors".
	// Multiply by Density (e.g. 100 lines).
	// We also want them to be irregular.

	// Use a 1D noise function on the angle.
	// `val = Noise(normalizedAngle * Density)`

	// To get sharp lines, we threshold the noise.
	// But we also want varying lengths.

	// Let's use a hash function that returns a random float 0..1 for a given integer "sector".
	// But lines can have width.

	// Simpler approach:
	// Continuous noise.
	// If Noise(angle * density) > Threshold, draw line.
	// Length of line:
	// minR = MinRadius + Noise2(angle * density) * (MaxRadius - MinRadius)
	// If dist > minR, draw.

	// We need a stateless noise function. I'll use a simple hash helper.

	seed := int64(p.Phase.Phase * 1000) // Phase affects seed/offset

	// High frequency noise for line presence
	n1 := noise1D(normalizedAngle * p.Density.Density, seed)

	// If n1 is high enough, we have a line.
	// Manga lines are usually black on white (or ink on paper).
	// "Overlay" means ink on transparent.

	// Threshold. Let's say 0.5.
	// But we want thin lines. Maybe > 0.8?
	// Or configurable.
	// Let's hardcode threshold for now or use `LineSize` somehow?
	// No, let's use a fixed threshold and rely on Density.

	if n1 > 0.6 {
		// This angle has a line.
		// Check radial distance.

		// Determine start distance for this line.
		// Vary it using a second noise channel (offset seed).
		n2 := noise1D(normalizedAngle * p.Density.Density, seed + 12345)

		// n2 is 0..1.
		// effectiveMinRadius = MinRadius + n2 * (MaxRadius - MinRadius)
		// Usually max radius is where lines *end*. But for speed lines, they often go to infinity.
		// The "MinRadius" is the "safe zone" radius.
		// So the line starts at `MinRadius + Variation`.

		// If MaxRadius is not set (0), assume infinity/image bound.
		// Actually, let's interpret MaxRadius as "Scatter amount" for the start point?
		// Or "Length"?
		// Standard: Inner Radius (safe zone) + Variance.

		effectiveStart := p.MinRadius.MinRadius
		if p.MaxRadius.MaxRadius > p.MinRadius.MinRadius {
             effectiveStart += n2 * (p.MaxRadius.MaxRadius - p.MinRadius.MinRadius)
		}

		if dist >= effectiveStart {
			// Draw line
			if p.LineImageSource.LineImageSource != nil {
				return p.LineImageSource.LineImageSource.At(x, y)
			}
			return p.LineColor.LineColor
		}
	}

	// Space
	if p.SpaceImageSource.SpaceImageSource != nil {
		return p.SpaceImageSource.SpaceImageSource.At(x, y)
	}
	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
	}
	return color.RGBA{} // Transparent
}

// noise1D returns a value 0..1
func noise1D(x float64, seed int64) float64 {
	// Simple hash of coordinate
	ix := int64(math.Floor(x))

	// Linear interpolation between integer points to smooth it?
	// Or just raw hash for jaggedness?
	// Speed lines are jagged.
	// But if we want width, we need blocks.
	// Let's just hash the integer part.

	h1 := hashInt(ix + seed)
	// h2 := hashInt(ix + 1 + seed)

	// If we just return h1, we get blocks of width 1/Density.
	// This creates constant width sectors.
	// This is fine.

	return h1
}

func hashInt(x int64) float64 {
	z := uint64(x)
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	z = z ^ (z >> 31)
	// Normalize to 0..1
	return float64(z) / float64(math.MaxUint64)
}

// NewSpeedLines creates a new SpeedLines pattern.
func NewSpeedLines(ops ...func(any)) image.Image {
	p := &SpeedLines{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Type: SpeedLinesRadial,
	}
	// Defaults
	p.LineColor.LineColor = color.Black
	p.SpaceColor.SpaceColor = nil // Transparent
	p.Density.Density = 100
	p.MinRadius.MinRadius = 50
	p.MaxRadius.MaxRadius = 100 // Variance

	for _, op := range ops {
		op(p)
	}
	return p
}

// Option for setting Linear Type
func SpeedLinesLinearType() func(any) {
	return func(i any) {
		if p, ok := i.(*SpeedLines); ok {
			p.Type = SpeedLinesLinear
		}
	}
}
