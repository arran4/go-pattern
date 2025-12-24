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

// Center configures the center point of a pattern.
type Center struct {
	CenterX, CenterY int
}

func (c *Center) SetCenter(x, y int) {
	c.CenterX = x
	c.CenterY = y
}

type hasCenter interface {
	SetCenter(int, int)
}

// SetCenter creates an option to set the center.
func SetCenter(x, y int) func(any) {
	return func(i any) {
		if h, ok := i.(hasCenter); ok {
			h.SetCenter(x, y)
		}
	}
}

// SpaceImageSource configures an image source for spaces in a pattern.
type SpaceImageSource struct {
	SpaceImageSource image.Image
}

func (s *SpaceImageSource) SetSpaceImageSource(v image.Image) {
	s.SpaceImageSource = v
}

type hasSpaceImageSource interface {
	SetSpaceImageSource(image.Image)
}

// SetSpaceImageSource creates an option to set the space image source.
func SetSpaceImageSource(v image.Image) func(any) {
	return func(i any) {
		if h, ok := i.(hasSpaceImageSource); ok {
			h.SetSpaceImageSource(v)
		}
	}
}

// MinRadius configures the minimum radius.
type MinRadius struct {
	MinRadius float64
}

func (s *MinRadius) SetMinRadius(v float64) {
	s.MinRadius = v
}

type hasMinRadius interface {
	SetMinRadius(float64)
}

// SetMinRadius creates an option to set the minimum radius.
func SetMinRadius(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasMinRadius); ok {
			h.SetMinRadius(v)
		}
	}
}

// MaxRadius configures the maximum radius.
type MaxRadius struct {
	MaxRadius float64
}

func (s *MaxRadius) SetMaxRadius(v float64) {
	s.MaxRadius = v
}

type hasMaxRadius interface {
	SetMaxRadius(float64)
}

// SetMaxRadius creates an option to set the maximum radius.
func SetMaxRadius(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasMaxRadius); ok {
			h.SetMaxRadius(v)
		}
	}
}

// Density configures the density of a pattern.
type Density struct {
	Density float64
}

func (s *Density) SetDensity(v float64) {
	s.Density = v
}

type hasDensity interface {
	SetDensity(float64)
}

// SetDensity creates an option to set the density.
func SetDensity(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasDensity); ok {
			h.SetDensity(v)
		}
	}
}

// Phase configures the phase/offset of a pattern.
type Phase struct {
	Phase float64
}

func (s *Phase) SetPhase(v float64) {
	s.Phase = v
}

type hasPhase interface {
	SetPhase(float64)
}

// SetPhase creates an option to set the phase.
func SetPhase(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasPhase); ok {
			h.SetPhase(v)
		}
	}
}

// Radius configures the radius of circles/dots in a pattern.
type Radius struct {
	Radius int
}

func (s *Radius) SetRadius(v int) {
	s.Radius = v
}

type hasRadius interface {
	SetRadius(int)
}

// SetRadius creates an option to set the radius.
func SetRadius(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasRadius); ok {
			h.SetRadius(v)
		}
	}
}

// Spacing configures the spacing/periodicity in a pattern.
type Spacing struct {
	Spacing int
}

func (s *Spacing) SetSpacing(v int) {
	s.Spacing = v
}

type hasSpacing interface {
	SetSpacing(int)
}

// SetSpacing creates an option to set the spacing.
func SetSpacing(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasSpacing); ok {
			h.SetSpacing(v)
		}
	}
}

// FillColor configures the fill color in a pattern (e.g. for dots).
type FillColor struct {
	FillColor color.Color
}

func (s *FillColor) SetFillColor(v color.Color) {
	s.FillColor = v
}

type hasFillColor interface {
	SetFillColor(color.Color)
}

// SetFillColor creates an option to set the fill color.
func SetFillColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(hasFillColor); ok {
			h.SetFillColor(v)
		}
	}
}

// FillImageSource configures an image source for fill in a pattern.
type FillImageSource struct {
	FillImageSource image.Image
}

func (s *FillImageSource) SetFillImageSource(v image.Image) {
	s.FillImageSource = v
}

type hasFillImageSource interface {
	SetFillImageSource(image.Image)
}

// SetFillImageSource creates an option to set the fill image source.
func SetFillImageSource(v image.Image) func(any) {
	return func(i any) {
		if h, ok := i.(hasFillImageSource); ok {
			h.SetFillImageSource(v)
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

// Angle configures an angle option.
type Angle struct {
	Angle float64
}

func (s *Angle) SetAngle(v float64) {
	s.Angle = v
}

type hasAngle interface {
	SetAngle(float64)
}

// SetAngle creates an option to set the angle.
func SetAngle(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasAngle); ok {
			h.SetAngle(v)
		}
	}
}

// Angles configures a list of angles in degrees.
type Angles struct {
	Angles []float64
}

func (s *Angles) SetAngles(v []float64) {
	s.Angles = v
}

type hasAngles interface {
	SetAngles([]float64)
}

// SetAngles creates an option to set the angles.
func SetAngles(v ...float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasAngles); ok {
			h.SetAngles(v)
		}
	}
}
