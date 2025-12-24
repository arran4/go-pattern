package pattern

import (
	"image"
	"image/color"
)

// SpaceSize configures the size of spaces in a pattern.
type SpaceSize struct {
	SpaceSize int
}

func (s *SpaceSize) SetSpaceSize(v int) {
	s.SpaceSize = v
}

type hasSpaceSize interface {
	SetSpaceSize(int)
}

// SetSpaceSize creates an option to set the space size.
func SetSpaceSize(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasSpaceSize); ok {
			h.SetSpaceSize(v)
		}
	}
}

// LineSize configures the thickness of lines in a pattern.
type LineSize struct {
	LineSize int
}

func (s *LineSize) SetLineSize(v int) {
	s.LineSize = v
}

type hasLineSize interface {
	SetLineSize(int)
}

// SetLineSize creates an option to set the line size.
func SetLineSize(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasLineSize); ok {
			h.SetLineSize(v)
		}
	}
}

// LineColor configures the color of lines in a pattern.
// Default should be black, but that's up to the consumer to initialize if not set.
type LineColor struct {
	LineColor color.Color
}

func (s *LineColor) SetLineColor(v color.Color) {
	s.LineColor = v
}

type hasLineColor interface {
	SetLineColor(color.Color)
}

// SetLineColor creates an option to set the line color.
func SetLineColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasLineColor); ok {
			h.SetLineColor(v)
		}
	}
}

// SpaceColor configures the color of spaces in a pattern.
type SpaceColor struct {
	SpaceColor color.Color
}

func (s *SpaceColor) SetSpaceColor(v color.Color) {
	s.SpaceColor = v
}

type hasSpaceColor interface {
	SetSpaceColor(color.Color)
}

// SetSpaceColor creates an option to set the space color.
func SetSpaceColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasSpaceColor); ok {
			h.SetSpaceColor(v)
		}
	}
}

// StartColor configures the start color for a gradient.
type StartColor struct {
	StartColor color.Color
}

func (s *StartColor) SetStartColor(v color.Color) {
	s.StartColor = v
}

type hasStartColor interface {
	SetStartColor(color.Color)
}

// SetStartColor creates an option to set the start color.
func SetStartColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasStartColor); ok {
			h.SetStartColor(v)
		}
	}
}

// EndColor configures the end color for a gradient.
type EndColor struct {
	EndColor color.Color
}

func (s *EndColor) SetEndColor(v color.Color) {
	s.EndColor = v
}

type hasEndColor interface {
	SetEndColor(color.Color)
}

// SetEndColor creates an option to set the end color.
func SetEndColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasEndColor); ok {
			h.SetEndColor(v)
		}
	}
}

// LineImageSource configures an image source for lines in a pattern.
type LineImageSource struct {
	LineImageSource image.Image
}

func (s *LineImageSource) SetLineImageSource(v image.Image) {
	s.LineImageSource = v
}

type hasLineImageSource interface {
	SetLineImageSource(image.Image)
}

// SetLineImageSource creates an option to set the line image source.
func SetLineImageSource(v image.Image) func(any) {
	return func(i any) {
		if h, ok := i.(hasLineImageSource); ok {
			h.SetLineImageSource(v)
		}
	}
}

// TrueColor configures the color used for "true" values in boolean/fuzzy operations.
type TrueColor struct {
	TrueColor color.Color
}

func (s *TrueColor) SetTrueColor(v color.Color) {
	s.TrueColor = v
}

type hasTrueColor interface {
	SetTrueColor(color.Color)
}

// SetTrueColor creates an option to set the "true" color.
func SetTrueColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasTrueColor); ok {
			h.SetTrueColor(v)
		}
	}
}

// FalseColor configures the color used for "false" values in boolean/fuzzy operations.
type FalseColor struct {
	FalseColor color.Color
}

func (s *FalseColor) SetFalseColor(v color.Color) {
	s.FalseColor = v
}

type hasFalseColor interface {
	SetFalseColor(color.Color)
}

// SetFalseColor creates an option to set the "false" color.
func SetFalseColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasFalseColor); ok {
			h.SetFalseColor(v)
		}
	}
}
