package pattern

import (
	"image"
	"image/color"
)

// BooleanModes demonstrates the different operating modes of BooleanImage.

var BooleanModesOutputFilename = "boolean_modes.png"
var BooleanModesZoomLevels = []int{}

const BooleanModesOrder = 38

// ExampleNewBooleanModes is a placeholder for documentation.
func ExampleNewBooleanModes() {
	// See GenerateBooleanModes for actual implementation.
}

func GenerateBooleanModes(b image.Rectangle) image.Image {
	// Setup Inputs
	// 1. Gopher (Color image)
	gopher := demoGopher(b)
	// 2. Stripes (Black and White)
	stripes := NewHorizontalLine(
		SetLineSize(20),
		SetSpaceSize(20),
		SetLineColor(color.RGBA{255, 0, 0, 128}),  // Red semi-transparent
		SetSpaceColor(color.RGBA{0, 0, 255, 128}), // Blue semi-transparent
		SetBounds(b),
	)

	// Mode 1: ComponentWise (Default for these inputs if colors not set, but let's be explicit)
	// Red&Gopher, Blue&Gopher
	m1 := NewAnd([]image.Image{gopher, stripes}, SetBooleanMode(ModeComponentWise), SetBounds(b))

	// Mode 2: Bitwise
	m2 := NewAnd([]image.Image{gopher, stripes}, SetBooleanMode(ModeBitwise), SetBounds(b))

	// Mode 3: Threshold (Strict Boolean Logic)
	// Using PredicateAverageGrayAbove(128)
	// If Gopher > 50% gray AND Stripes > 50% gray -> Green, else Black
	m3 := NewAnd(
		[]image.Image{gopher, stripes},
		SetBooleanMode(ModeThreshold),
		SetThreshold(0.5),
		SetPredicate(PredicateAverageGrayAbove(128)),
		SetTrueColor(color.RGBA{0, 255, 0, 255}),
		SetFalseColor(color.Black),
		SetBounds(b),
	)

	// Mode 4: Fuzzy (Interpolated)
	// Min(GopherGray, StripesGray) mapped to White->Green gradient
	m4 := NewAnd(
		[]image.Image{gopher, stripes},
		SetBooleanMode(ModeFuzzy),
		SetPredicate(PredicateFuzzyAlpha()), // Use Alpha channel
		SetTrueColor(color.RGBA{0, 255, 0, 255}),
		SetFalseColor(color.White),
		SetBounds(b),
	)

	return stitchImagesForDemo(m1, m2, m3, m4)
}

func GenerateBooleanModesReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": demoGopher,
		"Stripes": func(b image.Rectangle) image.Image {
			return NewHorizontalLine(
				SetLineSize(20),
				SetSpaceSize(20),
				SetLineColor(color.RGBA{255, 0, 0, 128}),
				SetSpaceColor(color.RGBA{0, 0, 255, 128}),
				SetBounds(b),
			)
		},
	}, []string{"Gopher", "Stripes"}
}

func init() {
	RegisterGenerator("BooleanModes", GenerateBooleanModes)
	RegisterReferences("BooleanModes", GenerateBooleanModesReferences)
}
