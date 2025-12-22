package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// Helpers for consistent demo inputs
func demoHorizontal(b image.Rectangle) image.Image {
	return NewHorizontalLine(
		SetLineSize(20),
		SetSpaceSize(20),
		SetLineColor(color.White),
		SetSpaceColor(color.Black),
		SetBounds(b),
	)
}

func demoVertical(b image.Rectangle) image.Image {
	return NewVerticalLine(
		SetLineSize(20),
		SetSpaceSize(20),
		SetLineColor(color.White),
		SetSpaceColor(color.Black),
		SetBounds(b),
	)
}

// AND Pattern

var AndOutputFilename = "boolean_and.png"
var AndZoomLevels = []int{2, 4}
const AndOrder = 20

func ExampleNewAnd() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black))

	// Use a predicate that considers White=True, Black=False.
	// Default is FuzzyAlpha which sees both as 1.0 (Opaque).
	// So we use AverageGrayAbove(128).
	pred := PredicateAverageGrayAbove(128)

	i := NewAnd([]image.Image{h, v}, SetPredicate(pred))

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
		SetPredicate(PredicateAverageGrayAbove(128)),
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
var OrZoomLevels = []int{2, 4}
const OrOrder = 21

func ExampleNewOr() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black))

	pred := PredicateAverageGrayAbove(128)

	i := NewOr([]image.Image{h, v}, SetPredicate(pred))

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
		SetPredicate(PredicateAverageGrayAbove(128)),
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
var XorZoomLevels = []int{2, 4}
const XorOrder = 22

func ExampleNewXor() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black))

	pred := PredicateAverageGrayAbove(128)

	i := NewXor([]image.Image{h, v}, SetPredicate(pred))

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
		SetPredicate(PredicateAverageGrayAbove(128)),
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
var NotZoomLevels = []int{2, 4}
const NotOrder = 23

func ExampleNewNot() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black))

	pred := PredicateAverageGrayAbove(128)

	i := NewNot(h, SetPredicate(pred))

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
		SetPredicate(PredicateAverageGrayAbove(128)),
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
