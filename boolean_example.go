package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// Helpers for consistent demo inputs
// Black lines on White background ("Ink on Paper")
func demoHorizontal(b image.Rectangle) image.Image {
	return NewHorizontalLine(
		SetLineSize(20),
		SetSpaceSize(20),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
		SetBounds(b),
	)
}

func demoVertical(b image.Rectangle) image.Image {
	return NewVerticalLine(
		SetLineSize(20),
		SetSpaceSize(20),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
		SetBounds(b),
	)
}

// PredicateInk returns 1.0 for Black (ink), 0.0 for White (paper).
// Used to perform logic on the "ink" rather than luminance.
func PredicateInk(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	// Calculate luminance or average.
	// White is 0xFFFF. Black is 0.
	avg := float64(r+g+b) / 3.0
	// Normalize to 0-1 (White=1, Black=0)
	v := avg / 65535.0
	// Invert: Black=1, White=0
	return 1.0 - v
}

// AND Pattern

var AndOutputFilename = "boolean_and.png"
var AndZoomLevels = []int{}
const AndOrder = 20

func ExampleNewAnd() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetSpaceColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetSpaceColor(color.White))

	// Use PredicateInk so Logic operates on Black lines.
	// Black=True, White=False.
	// AND(Black, Black) = Black.
	// Result should be Black (Ink). So we need SetTrueColor(Black).
	i := NewAnd([]image.Image{h, v}, SetPredicate(PredicateInk), SetTrueColor(color.Black), SetFalseColor(color.White))

	f, err := os.Create(AndOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateAnd(b image.Rectangle) image.Image {
	return NewAnd(
		[]image.Image{demoHorizontal(b), demoVertical(b)},
		SetPredicate(PredicateInk),
		SetTrueColor(color.Black),
		SetFalseColor(color.White),
		SetBounds(b),
	)
}

func GenerateAndReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Horizontal": demoHorizontal,
		"Vertical":   demoVertical,
	}, []string{"Horizontal", "Vertical"}
}


// OR Pattern

var OrOutputFilename = "boolean_or.png"
var OrZoomLevels = []int{}
const OrOrder = 21

func ExampleNewOr() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetSpaceColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetSpaceColor(color.White))

	i := NewOr([]image.Image{h, v}, SetPredicate(PredicateInk), SetTrueColor(color.Black), SetFalseColor(color.White))

	f, err := os.Create(OrOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateOr(b image.Rectangle) image.Image {
	return NewOr(
		[]image.Image{demoHorizontal(b), demoVertical(b)},
		SetPredicate(PredicateInk),
		SetTrueColor(color.Black),
		SetFalseColor(color.White),
		SetBounds(b),
	)
}

func GenerateOrReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Horizontal": demoHorizontal,
		"Vertical":   demoVertical,
	}, []string{"Horizontal", "Vertical"}
}


// XOR Pattern

var XorOutputFilename = "boolean_xor.png"
var XorZoomLevels = []int{}
const XorOrder = 22

func ExampleNewXor() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetSpaceColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetSpaceColor(color.White))

	i := NewXor([]image.Image{h, v}, SetPredicate(PredicateInk), SetTrueColor(color.Black), SetFalseColor(color.White))

	f, err := os.Create(XorOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateXor(b image.Rectangle) image.Image {
	return NewXor(
		[]image.Image{demoHorizontal(b), demoVertical(b)},
		SetPredicate(PredicateInk),
		SetTrueColor(color.Black),
		SetFalseColor(color.White),
		SetBounds(b),
	)
}

func GenerateXorReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Horizontal": demoHorizontal,
		"Vertical":   demoVertical,
	}, []string{"Horizontal", "Vertical"}
}


// NOT Pattern

var NotOutputFilename = "boolean_not.png"
var NotZoomLevels = []int{}
const NotOrder = 23

func ExampleNewNot() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetSpaceColor(color.White))

	i := NewNot(h, SetPredicate(PredicateInk), SetTrueColor(color.Black), SetFalseColor(color.White))

	f, err := os.Create(NotOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateNot(b image.Rectangle) image.Image {
	return NewNot(
		demoHorizontal(b),
		SetPredicate(PredicateInk),
		SetTrueColor(color.Black),
		SetFalseColor(color.White),
		SetBounds(b),
	)
}

func GenerateNotReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Horizontal": demoHorizontal,
	}, []string{"Horizontal"}
}

func init() {
	RegisterGenerator("And", GenerateAnd)
	RegisterReferences("And", GenerateAndReferences)

	RegisterGenerator("Or", GenerateOr)
	RegisterReferences("Or", GenerateOrReferences)

	RegisterGenerator("Xor", GenerateXor)
	RegisterReferences("Xor", GenerateXorReferences)

	RegisterGenerator("Not", GenerateNot)
	RegisterReferences("Not", GenerateNotReferences)
}
