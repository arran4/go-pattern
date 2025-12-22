package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure BooleanImage implements the image.Image interface.
var _ image.Image = (*BooleanImage)(nil)

type BooleanOpType int

const (
	OpAnd BooleanOpType = iota
	OpOr
	OpXor
	OpNot
)

// ColorPredicate converts a color to a fuzzy value between 0.0 and 1.0.
type ColorPredicate func(c color.Color) float64

// BooleanImage represents a boolean or fuzzy logic operation on input images.
type BooleanImage struct {
	Null
	Op BooleanOpType
	Inputs []image.Image
	Predicate ColorPredicate
	TrueColor
	FalseColor
}

func (bi *BooleanImage) At(x, y int) color.Color {
	if len(bi.Inputs) == 0 {
		return bi.FalseColor.FalseColor
	}

	var val float64

	switch bi.Op {
	case OpAnd:
		val = 1.0
		for _, input := range bi.Inputs {
			if input == nil {
				continue
			}
			v := bi.Predicate(input.At(x, y))
			if v < val {
				val = v
			}
		}
	case OpOr:
		val = 0.0
		for _, input := range bi.Inputs {
			if input == nil {
				continue
			}
			v := bi.Predicate(input.At(x, y))
			if v > val {
				val = v
			}
		}
	case OpXor:
		// XOR for fuzzy logic: |a - b|
		// For multiple inputs, it's cumulative? Xor(a, b, c) = Xor(Xor(a, b), c)
		// Let's assume binary or sequential.
		val = 0.0
		for i, input := range bi.Inputs {
			if input == nil {
				continue
			}
			v := bi.Predicate(input.At(x, y))
			if i == 0 {
				val = v
			} else {
				val = math.Abs(val - v)
			}
		}
	case OpNot:
		if len(bi.Inputs) > 0 && bi.Inputs[0] != nil {
			val = 1.0 - bi.Predicate(bi.Inputs[0].At(x, y))
		}
	}

	// Interpolate between FalseColor and TrueColor based on val
	return interpolateColor(bi.FalseColor.FalseColor, bi.TrueColor.TrueColor, val)
}

func interpolateColor(c0, c1 color.Color, t float64) color.Color {
	if t <= 0 {
		return c0
	}
	if t >= 1 {
		return c1
	}

	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()

	r := float64(r0) + t*(float64(r1)-float64(r0))
	g := float64(g0) + t*(float64(g1)-float64(g0))
	b := float64(b0) + t*(float64(b1)-float64(b0))
	a := float64(a0) + t*(float64(a1)-float64(a0))

	// RGBA() returns 16-bit values (0-0xffff).
	// We need to return a color.Color. The standard library doesn't have a generic 64-bit color.
	// We can use color.RGBA64.

	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

// Common predicates

// AlphaThreshold returns 1.0 if alpha >= threshold, else 0.0
func PredicateAlphaThreshold(threshold uint8) ColorPredicate {
	t := uint32(threshold)
	return func(c color.Color) float64 {
		_, _, _, a := c.RGBA()
		// a is 0-0xffff. We want to compare with threshold 0-0xff.
		if (a >> 8) >= t {
			return 1.0
		}
		return 0.0
	}
}

// RedAbove returns 1.0 if red >= threshold, else 0.0
func PredicateRedAbove(threshold uint8) ColorPredicate {
	t := uint32(threshold)
	return func(c color.Color) float64 {
		r, _, _, _ := c.RGBA()
		if (r >> 8) >= t {
			return 1.0
		}
		return 0.0
	}
}

// AverageGrayAbove returns 1.0 if average of RGB >= threshold, else 0.0
func PredicateAverageGrayAbove(threshold uint8) ColorPredicate {
	t := uint32(threshold)
	return func(c color.Color) float64 {
		r, g, b, _ := c.RGBA()
		avg := (r + g + b) / 3
		if (avg >> 8) >= t {
			return 1.0
		}
		return 0.0
	}
}

// FuzzyAlpha returns the alpha value as a float 0-1
func PredicateFuzzyAlpha() ColorPredicate {
	return func(c color.Color) float64 {
		_, _, _, a := c.RGBA()
		return float64(a) / 65535.0
	}
}

// FuzzyRed returns the red value as a float 0-1
func PredicateFuzzyRed() ColorPredicate {
	return func(c color.Color) float64 {
		r, _, _, _ := c.RGBA()
		return float64(r) / 65535.0
	}
}

// Default predicate
func DefaultPredicate(c color.Color) float64 {
	return PredicateFuzzyAlpha()(c)
}


// SetPredicate sets the predicate for the boolean operation.
type hasPredicate interface {
	SetPredicate(ColorPredicate)
}

func (bi *BooleanImage) SetPredicate(p ColorPredicate) {
	bi.Predicate = p
}

func SetPredicate(p ColorPredicate) func(any) {
	return func(i any) {
		if h, ok := i.(hasPredicate); ok {
			h.SetPredicate(p)
		}
	}
}

// Constructors

func NewAnd(inputs []image.Image, ops ...func(any)) image.Image {
	p := &BooleanImage{
		Null: Null{bounds: image.Rect(0, 0, 255, 255)},
		Op: OpAnd,
		Inputs: inputs,
		Predicate: DefaultPredicate,
	}
	p.TrueColor.TrueColor = color.White
	p.FalseColor.FalseColor = color.Black

	for _, op := range ops {
		op(p)
	}
	return p
}

func NewOr(inputs []image.Image, ops ...func(any)) image.Image {
	p := &BooleanImage{
		Null: Null{bounds: image.Rect(0, 0, 255, 255)},
		Op: OpOr,
		Inputs: inputs,
		Predicate: DefaultPredicate,
	}
	p.TrueColor.TrueColor = color.White
	p.FalseColor.FalseColor = color.Black

	for _, op := range ops {
		op(p)
	}
	return p
}

func NewXor(inputs []image.Image, ops ...func(any)) image.Image {
	p := &BooleanImage{
		Null: Null{bounds: image.Rect(0, 0, 255, 255)},
		Op: OpXor,
		Inputs: inputs,
		Predicate: DefaultPredicate,
	}
	p.TrueColor.TrueColor = color.White
	p.FalseColor.FalseColor = color.Black

	for _, op := range ops {
		op(p)
	}
	return p
}

func NewNot(input image.Image, ops ...func(any)) image.Image {
	p := &BooleanImage{
		Null: Null{bounds: image.Rect(0, 0, 255, 255)},
		Op: OpNot,
		Inputs: []image.Image{input},
		Predicate: DefaultPredicate,
	}
	p.TrueColor.TrueColor = color.White
	p.FalseColor.FalseColor = color.Black

	for _, op := range ops {
		op(p)
	}
	return p
}

// Demo variants

func NewDemoAnd(ops ...func(any)) image.Image {
	// Demo needs some inputs. We can create some default lines.
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	return NewAnd([]image.Image{h, v}, ops...)
}

func NewDemoOr(ops ...func(any)) image.Image {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	return NewOr([]image.Image{h, v}, ops...)
}

func NewDemoXor(ops ...func(any)) image.Image {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	return NewXor([]image.Image{h, v}, ops...)
}

func NewDemoNot(ops ...func(any)) image.Image {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	return NewNot(h, ops...)
}
