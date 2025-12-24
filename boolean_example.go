package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func demoGopher(b image.Rectangle) image.Image {
	return NewGopher()
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

// PredicateAnyAlpha returns 1.0 if there is any alpha (opaque), 0.0 if transparent.
// Use average with threshold.
func PredicateAnyAlpha(c color.Color) float64 {
	_, _, _, a := c.RGBA()
	if a > 0 {
		return 1.0
	}
	return 0.0
}

// AND Pattern

var AndOutputFilename = "boolean_and.png"
var AndZoomLevels = []int{}
const AndOrder = 20

func ExampleNewAnd() {
	// Gopher AND Horizontal Stripes
	g := NewGopher()
	// Line: Black (Alpha 1). Space: White (Alpha 1).
	// With Color Logic:
	// And(Gopher, Stripes) -> Min(Gopher, Stripes)
	// Stripes: Black lines, White spaces.
	// Where Space(White): Min(Gopher, White) -> Gopher.
	// Where Line(Black): Min(Gopher, Black) -> Black.
	// Result: Gopher visible through stripes.
	h := NewHorizontalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White))

	// Default uses component-wise min if no TrueColor/FalseColor set.
	i := NewAnd([]image.Image{g, h})

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
	h := NewHorizontalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))
	return NewAnd(
		[]image.Image{demoGopher(b), h},
		SetBounds(b),
	)
}

func GenerateAndReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
		"Stripes": func(b image.Rectangle) image.Image {
			return NewHorizontalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))
		},
	}, []string{"Gopher", "Stripes"}
}


// OR Pattern

var OrOutputFilename = "boolean_or.png"
var OrZoomLevels = []int{}
const OrOrder = 21

func ExampleNewOr() {
	g := NewGopher()
	v := NewVerticalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White))

	// OR(Gopher, Stripes) -> Max(Gopher, Stripes)
	// Stripes: Black lines, White spaces.
	// Where Space(White): Max(Gopher, White) -> White.
	// Where Line(Black): Max(Gopher, Black) -> Gopher.
	// Result: Gopher visible where Lines are Black. Masked White where Spaces are White.
	// This is effectively "Gopher masked by stripes (inverted)".
	i := NewOr([]image.Image{g, v})

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
	v := NewVerticalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))
	return NewOr(
		[]image.Image{demoGopher(b), v},
		SetBounds(b),
	)
}

func GenerateOrReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
		"Stripes": func(b image.Rectangle) image.Image {
			return NewVerticalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))
		},
	}, []string{"Gopher", "Stripes"}
}


// XOR Pattern

var XorOutputFilename = "boolean_xor.png"
var XorZoomLevels = []int{}
const XorOrder = 22

func ExampleNewXor() {
	g := NewGopher()
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))

	// XOR(Gopher, Stripes)
	// Use explicit colors to preserve the "Yellow" demo requested.
	i := NewXor([]image.Image{g, v}, SetTrueColor(color.RGBA{255, 255, 0, 255}), SetFalseColor(color.Transparent))

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
	vAlpha := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetBounds(b))
	return NewXor(
		[]image.Image{demoGopher(b), vAlpha},
		SetTrueColor(color.RGBA{255, 255, 0, 255}),
		SetFalseColor(color.Transparent),
		SetBounds(b),
	)
}

func GenerateXorReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
		"Stripes": func(b image.Rectangle) image.Image {
			return NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetBounds(b))
		},
	}, []string{"Gopher", "Stripes"}
}


// NOT Pattern

var NotOutputFilename = "boolean_not.png"
var NotZoomLevels = []int{}
const NotOrder = 23

func ExampleNewNot() {
	g := NewGopher()

	// Not Gopher.
	// Default component-wise: Invert colors.
	i := NewNot(g)

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
		demoGopher(b),
		SetBounds(b),
	)
}

func GenerateNotReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
	}, []string{"Gopher"}
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
