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
	Palette []color.Color
}

func (p *EightByEight) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *EightByEight) Bounds() image.Rectangle {
	return p.bounds
}

func (p *EightByEight) SetPalette(colors []color.Color) {
	p.Palette = colors
}

func (p *EightByEight) At(x, y int) color.Color {
	xpRaw := x % 8

	useMultiColor := len(p.Palette) > 2
	var fgIdx, bgIdx int
	if useMultiColor {
		bgIdx = p.Mode % len(p.Palette)
		fgIdx = (p.Mode / len(p.Palette)) % len(p.Palette)
	}

	// The original logic checks if mode >= xpRaw.
	// In multi-color mode, this check is bypassed (implied || useMultiColor).
	if useMultiColor || p.Mode >= xpRaw {
		xpAbs := xpRaw
		if xpAbs < 0 {
			xpAbs = -xpAbs
		}

		// Calculate south value for this column (xpAbs).
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
			if useMultiColor {
				return p.Palette[fgIdx]
			}
			if p.LineColor.LineColor != nil {
				return p.LineColor.LineColor
			}
			if len(p.Palette) > 1 {
				return p.Palette[1]
			}
			return color.Black
		}
	}

	if useMultiColor {
		return p.Palette[bgIdx]
	}
	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
	}
	if len(p.Palette) > 0 {
		return p.Palette[0]
	}
	return color.White
}

// NewEightByEight creates a new EightByEight pattern.
func NewEightByEight(mode int, ops ...func(any)) image.Image {
	p := &EightByEight{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Mode: mode,
	}
	// Note: We do NOT set default LineColor/SpaceColor here to allow Palette fallbacks to work.
	// They will be nil by default, triggering the fallback logic in At().

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
