package pattern

import (
	"image"
	"image/color"
	"image/draw"
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

func stitchImagesForDemo(images ...image.Image) image.Image {
	if len(images) == 0 {
		return nil
	}

	width := 0
	height := 0
	for _, img := range images {
		b := img.Bounds()
		width += b.Dx()
		if b.Dy() > height {
			height = b.Dy()
		}
	}
	padding := 10
	width += padding * (len(images) - 1)

	out := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(out, out.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	x := 0
	for _, img := range images {
		b := img.Bounds()
		r := image.Rect(x, 0, x+b.Dx(), b.Dy())
		draw.Draw(out, r, img, b.Min, draw.Over)
		x += b.Dx() + padding
	}

	return out
}

// BooleanAnd Pattern

var BooleanAndOutputFilename = "boolean_and.png"
var BooleanAndZoomLevels = []int{}
const BooleanAndOrder = 20

func ExampleNewBooleanAnd() {
	// Gopher AND Horizontal Stripes
	g := NewGopher()
	// Line: Black (Alpha 1). Space: White (Alpha 1).
	h := NewHorizontalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White))

	// Default uses component-wise min if no TrueColor/FalseColor set.
	i := NewAnd([]image.Image{g, h})

	f, err := os.Create(BooleanAndOutputFilename)
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

func GenerateBooleanAnd(b image.Rectangle) image.Image {
	h := NewHorizontalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))

	// Variant 1: Component-wise Min (Default)
	v1 := NewAnd([]image.Image{demoGopher(b), h}, SetBounds(b))

	// Variant 2: Silhouette (Cyan mask)
	v2 := NewAnd(
		[]image.Image{demoGopher(b), h},
		SetTrueColor(color.RGBA{0, 255, 255, 255}), // Cyan
		SetFalseColor(color.Transparent),
		SetBounds(b),
	)

	return stitchImagesForDemo(v1, v2)
}

func GenerateBooleanAndReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
		"Stripes": func(b image.Rectangle) image.Image {
			return NewHorizontalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))
		},
	}, []string{"Gopher", "Stripes"}
}


// BooleanOr Pattern

var BooleanOrOutputFilename = "boolean_or.png"
var BooleanOrZoomLevels = []int{}
const BooleanOrOrder = 21

func ExampleNewBooleanOr() {
	g := NewGopher()
	v := NewVerticalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White))

	// OR(Gopher, Stripes) -> Max(Gopher, Stripes)
	i := NewOr([]image.Image{g, v})

	f, err := os.Create(BooleanOrOutputFilename)
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

func GenerateBooleanOr(b image.Rectangle) image.Image {
	v := NewVerticalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))

	// Variant 1: Component-wise Max (Default)
	v1 := NewOr([]image.Image{demoGopher(b), v}, SetBounds(b))

	// Variant 2: Silhouette (Magenta mask)
	v2 := NewOr(
		[]image.Image{demoGopher(b), v},
		SetTrueColor(color.RGBA{255, 0, 255, 255}), // Magenta
		SetFalseColor(color.Transparent),
		SetBounds(b),
	)

	return stitchImagesForDemo(v1, v2)
}

func GenerateBooleanOrReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
		"Stripes": func(b image.Rectangle) image.Image {
			return NewVerticalLine(SetLineSize(10), SetSpaceSize(10), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))
		},
	}, []string{"Gopher", "Stripes"}
}


// BooleanXor Pattern

var BooleanXorOutputFilename = "boolean_xor.png"
var BooleanXorZoomLevels = []int{}
const BooleanXorOrder = 22

func ExampleNewBooleanXor() {
	g := NewGopher()
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))

	// XOR(Gopher, Stripes)
	i := NewXor([]image.Image{g, v}, SetTrueColor(color.RGBA{255, 255, 0, 255}), SetFalseColor(color.Transparent))

	f, err := os.Create(BooleanXorOutputFilename)
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

func GenerateBooleanXor(b image.Rectangle) image.Image {
	vAlpha := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetBounds(b))

	// Variant 1: Component-wise AbsDiff (Default)
	// Stripes need to be white background for component-wise logic to match visual expectations?
	// VerticalLine default uses White space.
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetSpaceColor(color.White), SetBounds(b))
	v1 := NewXor([]image.Image{demoGopher(b), v}, SetBounds(b))

	// Variant 2: Silhouette (Yellow mask)
	// Uses the Alpha version of lines
	v2 := NewXor(
		[]image.Image{demoGopher(b), vAlpha},
		SetTrueColor(color.RGBA{255, 255, 0, 255}),
		SetFalseColor(color.Transparent),
		SetBounds(b),
	)

	return stitchImagesForDemo(v1, v2)
}

func GenerateBooleanXorReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
		"Stripes": func(b image.Rectangle) image.Image {
			return NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black), SetBounds(b))
		},
	}, []string{"Gopher", "Stripes"}
}


// BooleanNot Pattern

var BooleanNotOutputFilename = "boolean_not.png"
var BooleanNotZoomLevels = []int{}
const BooleanNotOrder = 23

func ExampleNewBooleanNot() {
	g := NewGopher()

	// Not Gopher.
	// Default component-wise: Invert colors.
	i := NewNot(g)

	f, err := os.Create(BooleanNotOutputFilename)
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

func GenerateBooleanNot(b image.Rectangle) image.Image {
	// Variant 1: Inverted colors (Default)
	v1 := NewNot(demoGopher(b), SetBounds(b))

	// Variant 2: Silhouette (Green mask)
	v2 := NewNot(
		demoGopher(b),
		SetTrueColor(color.RGBA{0, 255, 0, 255}),
		SetFalseColor(color.Transparent),
		SetBounds(b),
	)

	return stitchImagesForDemo(v1, v2)
}

func GenerateBooleanNotReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
	}, []string{"Gopher"}
}

func init() {
	RegisterGenerator("BooleanAnd", GenerateBooleanAnd)
	RegisterReferences("BooleanAnd", GenerateBooleanAndReferences)

	RegisterGenerator("BooleanOr", GenerateBooleanOr)
	RegisterReferences("BooleanOr", GenerateBooleanOrReferences)

	RegisterGenerator("BooleanXor", GenerateBooleanXor)
	RegisterReferences("BooleanXor", GenerateBooleanXorReferences)

	RegisterGenerator("BooleanNot", GenerateBooleanNot)
	RegisterReferences("BooleanNot", GenerateBooleanNotReferences)
}
