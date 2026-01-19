package pattern

import (
	"fmt"
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
	OpBitwiseAnd
	OpBitwiseOr
	OpBitwiseXor
	OpBitwiseNot
)

type BooleanMode int

const (
	ModeAuto          BooleanMode = iota // Infer mode based on configuration
	ModeFuzzy                            // Predicate -> Float Logic -> Interpolate
	ModeThreshold                        // Predicate -> Threshold (0/1) -> Boolean Logic -> Select Color (Strict True/False)
	ModeComponentWise                    // R=Min(Ra, Rb), etc.
	ModeBitwise                          // R=Ra & Rb, etc.
)

// ColorPredicate converts a color to a fuzzy value between 0.0 and 1.0.
type ColorPredicate func(c color.Color) float64

// BooleanImage represents a boolean or fuzzy logic operation on input images.
type BooleanImage struct {
	Null
	Op        BooleanOpType
	Inputs    []image.Image
	Predicate ColorPredicate
	TrueColor
	FalseColor
	Mode      BooleanMode
	Threshold float64
}

func (bi *BooleanImage) At(x, y int) color.Color {
	if len(bi.Inputs) == 0 {
		if bi.FalseColor.FalseColor != nil {
			return bi.FalseColor.FalseColor
		}
		return color.RGBA{}
	}

	mode := bi.Mode
	if mode == ModeAuto {
		if bi.Op >= OpBitwiseAnd {
			mode = ModeBitwise
		} else if isZeroColor(bi.TrueColor.TrueColor) && isZeroColor(bi.FalseColor.FalseColor) {
			mode = ModeComponentWise
		} else {
			mode = ModeFuzzy
		}
	}

	switch mode {
	case ModeComponentWise:
		return bi.atComponentWise(x, y)
	case ModeBitwise:
		return bi.atBitwise(x, y)
	case ModeThreshold:
		return bi.atThreshold(x, y)
	case ModeFuzzy:
		fallthrough
	default:
		return bi.atFuzzy(x, y)
	}
}

func (bi *BooleanImage) atFuzzy(x, y int) color.Color {
	var val float64

	// Helper to get fuzzy value, handling nil inputs safely
	getVal := func(img image.Image) float64 {
		if img == nil {
			return 0
		}
		return bi.Predicate(img.At(x, y))
	}

	switch bi.Op {
	case OpAnd, OpBitwiseAnd:
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
	case OpOr, OpBitwiseOr:
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
	case OpXor, OpBitwiseXor:
		val = 0.0
		if len(bi.Inputs) >= 1 {
			val = getVal(bi.Inputs[0])
		}
		if len(bi.Inputs) >= 2 {
			v := getVal(bi.Inputs[1])
			val = math.Abs(val - v)
		}
	case OpNot, OpBitwiseNot:
		if len(bi.Inputs) > 0 {
			val = 1.0 - getVal(bi.Inputs[0])
		}
	}

	tc := bi.TrueColor.TrueColor
	fc := bi.FalseColor.FalseColor
	if tc == nil {
		tc = color.White
	}
	if fc == nil {
		fc = color.Black
	}

	return interpolateColor(fc, tc, val)
}

func (bi *BooleanImage) atThreshold(x, y int) color.Color {
	// Logic is similar to Fuzzy, but intermediate values are binarized against Threshold
	// Actually, "Based on threshold ... to determine true / false".
	// Should we threshold the inputs? Or the result?
	// Usually: Input -> Threshold -> Boolean Logic -> Output.

	threshold := bi.Threshold
	if threshold == 0 {
		threshold = 0.5
	}

	getBool := func(img image.Image) bool {
		if img == nil {
			return false
		}
		return bi.Predicate(img.At(x, y)) >= threshold
	}

	var result bool

	switch bi.Op {
	case OpAnd, OpBitwiseAnd:
		result = true
		for _, input := range bi.Inputs {
			if !getBool(input) {
				result = false
				break
			}
		}
	case OpOr, OpBitwiseOr:
		result = false
		for _, input := range bi.Inputs {
			if getBool(input) {
				result = true
				break
			}
		}
	case OpXor, OpBitwiseXor:
		b1 := false
		b2 := false
		if len(bi.Inputs) >= 1 {
			b1 = getBool(bi.Inputs[0])
		}
		if len(bi.Inputs) >= 2 {
			b2 = getBool(bi.Inputs[1])
		}
		result = b1 != b2
	case OpNot, OpBitwiseNot:
		if len(bi.Inputs) > 0 {
			result = !getBool(bi.Inputs[0])
		}
	}

	tc := bi.TrueColor.TrueColor
	fc := bi.FalseColor.FalseColor
	if tc == nil {
		tc = color.White
	}
	if fc == nil {
		fc = color.Black
	}

	if result {
		return tc
	}
	return fc
}

func (bi *BooleanImage) atComponentWise(x, y int) color.Color {
	// Maps OpBitwise... to normal behavior for ComponentWise
	switch bi.Op {
	case OpAnd, OpBitwiseAnd:
		// Component-wise Min
		var minC color.Color
		for i, input := range bi.Inputs {
			if input == nil {
				continue
			}
			c := input.At(x, y)
			if i == 0 || minC == nil {
				minC = c
			} else {
				minC = minColor(minC, c)
			}
		}
		if minC == nil {
			return color.RGBA{}
		}
		return minC
	case OpOr, OpBitwiseOr:
		// Component-wise Max
		var maxC color.Color
		for i, input := range bi.Inputs {
			if input == nil {
				continue
			}
			c := input.At(x, y)
			if i == 0 || maxC == nil {
				maxC = c
			} else {
				maxC = maxColor(maxC, c)
			}
		}
		if maxC == nil {
			return color.RGBA{}
		}
		return maxC
	case OpXor, OpBitwiseXor:
		// Component-wise AbsDiff
		if len(bi.Inputs) < 2 {
			if len(bi.Inputs) == 1 && bi.Inputs[0] != nil {
				return bi.Inputs[0].At(x, y)
			}
			return color.RGBA{}
		}
		c1 := bi.Inputs[0].At(x, y)
		c2 := bi.Inputs[1].At(x, y)
		return absDiffColor(c1, c2)
	case OpNot, OpBitwiseNot:
		if len(bi.Inputs) > 0 && bi.Inputs[0] != nil {
			return invertColor(bi.Inputs[0].At(x, y))
		}
		return color.RGBA{}
	}
	return color.RGBA{}
}

func (bi *BooleanImage) atBitwise(x, y int) color.Color {
	switch bi.Op {
	case OpAnd, OpBitwiseAnd:
		var res color.Color
		for i, input := range bi.Inputs {
			if input == nil {
				continue
			}
			c := input.At(x, y)
			if i == 0 || res == nil {
				res = c
			} else {
				res = bitwiseAndColor(res, c)
			}
		}
		if res == nil {
			return color.RGBA{}
		}
		return res
	case OpOr, OpBitwiseOr:
		var res color.Color
		for i, input := range bi.Inputs {
			if input == nil {
				continue
			}
			c := input.At(x, y)
			if i == 0 || res == nil {
				res = c
			} else {
				res = bitwiseOrColor(res, c)
			}
		}
		if res == nil {
			return color.RGBA{}
		}
		return res
	case OpXor, OpBitwiseXor:
		if len(bi.Inputs) < 2 {
			if len(bi.Inputs) == 1 && bi.Inputs[0] != nil {
				return bi.Inputs[0].At(x, y)
			}
			return color.RGBA{}
		}
		var res color.Color
		for i, input := range bi.Inputs {
			if input == nil {
				continue
			}
			c := input.At(x, y)
			if i == 0 || res == nil {
				res = c
			} else {
				res = bitwiseXorColor(res, c)
			}
		}
		if res == nil {
			return color.RGBA{}
		}
		return res
	case OpNot, OpBitwiseNot:
		if len(bi.Inputs) > 0 && bi.Inputs[0] != nil {
			return bitwiseNotColor(bi.Inputs[0].At(x, y))
		}
		return color.RGBA{}
	}
	return color.RGBA{}
}

func isZeroColor(c color.Color) bool {
	if c == nil {
		return true
	}
	r, g, b, a := c.RGBA()
	return r == 0 && g == 0 && b == 0 && a == 0
}

func minColor(c1, c2 color.Color) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return color.RGBA64{
		R: uint16(min(r1, r2)),
		G: uint16(min(g1, g2)),
		B: uint16(min(b1, b2)),
		A: uint16(min(a1, a2)),
	}
}

func maxColor(c1, c2 color.Color) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return color.RGBA64{
		R: uint16(max(r1, r2)),
		G: uint16(max(g1, g2)),
		B: uint16(max(b1, b2)),
		A: uint16(max(a1, a2)),
	}
}

func absDiffColor(c1, c2 color.Color) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return color.RGBA64{
		R: uint16(absDiff(r1, r2)),
		G: uint16(absDiff(g1, g2)),
		B: uint16(absDiff(b1, b2)),
		A: uint16(absDiff(a1, a2)),
	}
}

func invertColor(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA64{
		R: uint16(0xFFFF - r),
		G: uint16(0xFFFF - g),
		B: uint16(0xFFFF - b),
		A: uint16(a),
	}
}

// Bitwise helpers
func bitwiseAndColor(c1, c2 color.Color) color.Color {
	n1 := color.NRGBA64Model.Convert(c1).(color.NRGBA64)
	n2 := color.NRGBA64Model.Convert(c2).(color.NRGBA64)

	return color.RGBA64Model.Convert(color.NRGBA64{
		R: n1.R & n2.R,
		G: n1.G & n2.G,
		B: n1.B & n2.B,
		A: n1.A & n2.A,
	})
}

func bitwiseOrColor(c1, c2 color.Color) color.Color {
	n1 := color.NRGBA64Model.Convert(c1).(color.NRGBA64)
	n2 := color.NRGBA64Model.Convert(c2).(color.NRGBA64)

	return color.RGBA64Model.Convert(color.NRGBA64{
		R: n1.R | n2.R,
		G: n1.G | n2.G,
		B: n1.B | n2.B,
		A: n1.A | n2.A,
	})
}

func bitwiseXorColor(c1, c2 color.Color) color.Color {
	n1 := color.NRGBA64Model.Convert(c1).(color.NRGBA64)
	n2 := color.NRGBA64Model.Convert(c2).(color.NRGBA64)

	return color.RGBA64Model.Convert(color.NRGBA64{
		R: n1.R ^ n2.R,
		G: n1.G ^ n2.G,
		B: n1.B ^ n2.B,
		A: n1.A ^ n2.A,
	})
}

func bitwiseNotColor(c color.Color) color.Color {
	n := color.NRGBA64Model.Convert(c).(color.NRGBA64)
	return color.RGBA64Model.Convert(color.NRGBA64{
		R: ^n.R,
		G: ^n.G,
		B: ^n.B,
		A: n.A,
	})
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func absDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
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

// SetBooleanMode sets the mode for the boolean operation.
type hasBooleanMode interface {
	SetBooleanMode(BooleanMode)
}

func (bi *BooleanImage) SetBooleanMode(m BooleanMode) {
	bi.Mode = m
}

func SetBooleanMode(m BooleanMode) func(any) {
	return func(i any) {
		if h, ok := i.(hasBooleanMode); ok {
			h.SetBooleanMode(m)
		}
	}
}

// SetThreshold sets the threshold for the boolean operation (used in ModeThreshold).
type hasThreshold interface {
	SetThreshold(float64)
}

func (bi *BooleanImage) SetThreshold(t float64) {
	bi.Threshold = t
}

func SetThreshold(t float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasThreshold); ok {
			h.SetThreshold(t)
		}
	}
}

// Specific Types

// And represents a boolean AND operation.
type And struct {
	BooleanImage
}

// Or represents a boolean OR operation.
type Or struct {
	BooleanImage
}

// Xor represents a boolean XOR operation.
type Xor struct {
	BooleanImage
}

// Not represents a boolean NOT operation.
type Not struct {
	BooleanImage
}

// Constructors

// NewAnd creates a new And pattern.
func NewAnd(inputs []image.Image, ops ...func(any)) image.Image {
	p := &And{
		BooleanImage: BooleanImage{
			Null:      Null{bounds: image.Rect(0, 0, 255, 255)},
			Op:        OpAnd,
			Inputs:    inputs,
			Predicate: DefaultPredicate,
		},
	}
	// Defaults are nil (zero) to allow Color Logic

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewOr creates a new Or pattern.
func NewOr(inputs []image.Image, ops ...func(any)) image.Image {
	p := &Or{
		BooleanImage: BooleanImage{
			Null:      Null{bounds: image.Rect(0, 0, 255, 255)},
			Op:        OpOr,
			Inputs:    inputs,
			Predicate: DefaultPredicate,
		},
	}
	// Defaults are nil

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewXor creates a new Xor pattern. It enforces exactly 2 inputs.
func NewXor(inputs []image.Image, ops ...func(any)) image.Image {
	if len(inputs) != 2 {
		panic(fmt.Sprintf("Xor requires exactly 2 inputs, got %d", len(inputs)))
	}
	p := &Xor{
		BooleanImage: BooleanImage{
			Null:      Null{bounds: image.Rect(0, 0, 255, 255)},
			Op:        OpXor,
			Inputs:    inputs,
			Predicate: DefaultPredicate,
		},
	}
	// Defaults are nil

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewNot creates a new Not pattern. It enforces exactly 1 input.
func NewNot(input image.Image, ops ...func(any)) image.Image {
	p := &Not{
		BooleanImage: BooleanImage{
			Null:      Null{bounds: image.Rect(0, 0, 255, 255)},
			Op:        OpNot,
			Inputs:    []image.Image{input},
			Predicate: DefaultPredicate,
		},
	}
	// Defaults are nil

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewBitwiseAnd creates a new Bitwise And pattern.
func NewBitwiseAnd(inputs []image.Image, ops ...func(any)) image.Image {
	p := &And{
		BooleanImage: BooleanImage{
			Null:   Null{bounds: image.Rect(0, 0, 255, 255)},
			Op:     OpBitwiseAnd,
			Inputs: inputs,
		},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NewBitwiseOr creates a new Bitwise Or pattern.
func NewBitwiseOr(inputs []image.Image, ops ...func(any)) image.Image {
	p := &Or{
		BooleanImage: BooleanImage{
			Null:   Null{bounds: image.Rect(0, 0, 255, 255)},
			Op:     OpBitwiseOr,
			Inputs: inputs,
		},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NewBitwiseXor creates a new Bitwise Xor pattern.
func NewBitwiseXor(inputs []image.Image, ops ...func(any)) image.Image {
	if len(inputs) != 2 {
		panic(fmt.Sprintf("BitwiseXor requires exactly 2 inputs, got %d", len(inputs)))
	}
	p := &Xor{
		BooleanImage: BooleanImage{
			Null:   Null{bounds: image.Rect(0, 0, 255, 255)},
			Op:     OpBitwiseXor,
			Inputs: inputs,
		},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NewBitwiseNot creates a new Bitwise Not pattern.
func NewBitwiseNot(input image.Image, ops ...func(any)) image.Image {
	p := &Not{
		BooleanImage: BooleanImage{
			Null:   Null{bounds: image.Rect(0, 0, 255, 255)},
			Op:     OpBitwiseNot,
			Inputs: []image.Image{input},
		},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// Demo variants

func NewDemoAnd(ops ...func(any)) image.Image {
	// Demo needs some inputs. We can create some default lines.
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	return NewAnd([]image.Image{h, v}, ops...)
}

func NewDemoOr(ops ...func(any)) image.Image {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	return NewOr([]image.Image{h, v}, ops...)
}

func NewDemoXor(ops ...func(any)) image.Image {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	return NewXor([]image.Image{h, v}, ops...)
}

func NewDemoNot(ops ...func(any)) image.Image {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	return NewNot(h, ops...)
}
