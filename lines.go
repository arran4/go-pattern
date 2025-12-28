package pattern

import (
	"image"
	"image/color"
)

// Ensure HorizontalLine implements the image.Image interface.
var _ image.Image = (*HorizontalLine)(nil)

// HorizontalLine is a pattern that draws horizontal lines.
type HorizontalLine struct {
	Null
	SpaceSize
	LineSize
	LineColor
	SpaceColor
	LineImageSource
	Phase
}

func (p *HorizontalLine) At(x, y int) color.Color {
	ls := p.LineSize.LineSize
	ss := p.SpaceSize.SpaceSize
	period := ls + ss
	if period == 0 {
		return p.LineColor.LineColor
	}

	// Apply phase offset
	offsetY := y - int(p.Phase.Phase)

	// Handle negative coordinates correctly for modulo
	mod := offsetY % period
	if mod < 0 {
		mod += period
	}

	if mod < ls {
		if p.LineImageSource.LineImageSource != nil {
			return p.LineImageSource.LineImageSource.At(x, y)
		}
		return p.LineColor.LineColor
	}

	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
	}
	return color.RGBA{}
}

// NewHorizontalLine creates a new HorizontalLine pattern.
func NewHorizontalLine(ops ...func(any)) image.Image {
	p := &HorizontalLine{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.LineSize.LineSize = 1
	p.SpaceSize.SpaceSize = 1
	p.LineColor.LineColor = color.Black
	// SpaceColor defaults to nil (transparent)

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoHorizontalLine produces a demo variant for readme.md pre-populated values
func NewDemoHorizontalLine(ops ...func(any)) image.Image {
	return NewHorizontalLine(ops...)
}

// Ensure VerticalLine implements the image.Image interface.
var _ image.Image = (*VerticalLine)(nil)

// VerticalLine is a pattern that draws vertical lines.
type VerticalLine struct {
	Null
	SpaceSize
	LineSize
	LineColor
	SpaceColor
	LineImageSource
	Phase
}

func (p *VerticalLine) At(x, y int) color.Color {
	ls := p.LineSize.LineSize
	ss := p.SpaceSize.SpaceSize
	period := ls + ss
	if period == 0 {
		return p.LineColor.LineColor
	}

	// Apply phase offset
	offsetX := x - int(p.Phase.Phase)

	// Handle negative coordinates correctly for modulo
	mod := offsetX % period
	if mod < 0 {
		mod += period
	}

	if mod < ls {
		if p.LineImageSource.LineImageSource != nil {
			return p.LineImageSource.LineImageSource.At(x, y)
		}
		return p.LineColor.LineColor
	}

	if p.SpaceColor.SpaceColor != nil {
		return p.SpaceColor.SpaceColor
	}
	return color.RGBA{}
}

// NewVerticalLine creates a new VerticalLine pattern.
func NewVerticalLine(ops ...func(any)) image.Image {
	p := &VerticalLine{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.LineSize.LineSize = 1
	p.SpaceSize.SpaceSize = 1
	p.LineColor.LineColor = color.Black
	// SpaceColor defaults to nil (transparent)

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoVerticalLine produces a demo variant for readme.md pre-populated values
func NewDemoVerticalLine(ops ...func(any)) image.Image {
	return NewVerticalLine(ops...)
}
