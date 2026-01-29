package pattern

import (
	"image"
	"image/color"
)

// Ensure EightByEight implements the image.Image interface.
var _ image.Image = (*EightByEight)(nil)

// EightByEight represents a set of 8x8 tiling patterns.
// Based on https://github.com/arran4/eightbyeight
type EightByEight struct {
	Null
	Mode int
	LineColor
	SpaceColor
}

func (p *EightByEight) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *EightByEight) Bounds() image.Rectangle {
	return p.bounds
}

func (p *EightByEight) At(x, y int) color.Color {
	xpRaw := x % 8
	// The original logic checks if mode >= xpRaw.
	// Since xpRaw is in [-7, 7], and mode is typically >= 0.
	// If mode is 0, and xpRaw is 1, check fails -> returns SpaceColor.
	if p.Mode >= xpRaw {
		xpAbs := xpRaw
		if xpAbs < 0 {
			xpAbs = -xpAbs
		}

		// Calculate south value for this column (xpAbs).
		// Original logic initializes south to {1, 0, 0, 0}.
		// Then fills it with base-4 digits of mode.
		// If mode=0, loop doesn't run, south remains {1, 0, 0, 0}.
		idx := xpAbs % 4
		svVal := (p.Mode >> (2 * idx)) & 3
		if idx == 0 && p.Mode == 0 {
			svVal = 1
		}

		// sv = 2^(4 - svVal)
		// svVal is 0..3 -> sv is 16, 8, 4, 2
		sv := 1 << (4 - svVal)

		// dp corresponds to diagonal position
		dp := (y - xpRaw) % 8
		if dp < 0 {
			dp += 8
		}

		if sv > 0 && dp%sv == 0 {
			if p.LineColor.LineColor != nil {
				return p.LineColor.LineColor
			}
			return color.Black // Default should be handled in New, but safe fallback
		}
	}

	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
	}
	return color.White // Default fallback
}

// NewEightByEight creates a new EightByEight pattern.
func NewEightByEight(mode int, ops ...func(any)) image.Image {
	p := &EightByEight{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Mode: mode,
	}
	// Defaults
	p.LineColor.LineColor = color.Black
	p.SpaceColor.SpaceColor = color.White

	for _, op := range ops {
		op(p)
	}
	return p
}

// SetMode sets the mode for the EightByEight pattern.
func SetMode(mode int) func(any) {
	return func(a any) {
		if p, ok := a.(*EightByEight); ok {
			p.Mode = mode
		}
	}
}
