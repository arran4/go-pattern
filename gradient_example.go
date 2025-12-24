package pattern

import (
	"image"
	"image/color"
)

// Linear Gradient (Horizontal) Pattern

var LinearGradientOutputFilename = "linear_gradient.png"
var LinearGradientZoomLevels = []int{}
const LinearGradientOrder = 30
const LinearGradientBaseLabel = "Horizontal"

func ExampleNewLinearGradient() {
	// Linear Gradient (Horizontal)
	NewLinearGradient(
		SetStartColor(color.RGBA{255, 0, 0, 255}),
		SetEndColor(color.RGBA{0, 0, 255, 255}),
	)
}

func GenerateLinearGradient(b image.Rectangle) image.Image {
	return NewLinearGradient(
		SetStartColor(color.RGBA{255, 0, 0, 255}),
		SetEndColor(color.RGBA{0, 0, 255, 255}),
		SetBounds(b),
	)
}

func GenerateLinearGradientReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Vertical": func(b image.Rectangle) image.Image {
			return NewLinearGradient(
				SetStartColor(color.RGBA{0, 255, 0, 255}),
				SetEndColor(color.RGBA{255, 255, 0, 255}),
				GradientVertical(),
				SetBounds(b),
			)
		},
	}, []string{"Vertical"}
}


// Radial Gradient Pattern

var RadialGradientOutputFilename = "radial_gradient.png"
var RadialGradientZoomLevels = []int{}
const RadialGradientOrder = 31

func ExampleNewRadialGradient() {
	// Radial Gradient
	NewRadialGradient(
		SetStartColor(color.RGBA{255, 0, 0, 255}),
		SetEndColor(color.RGBA{0, 0, 255, 255}),
	)
}

func GenerateRadialGradient(b image.Rectangle) image.Image {
	return NewRadialGradient(
		SetStartColor(color.RGBA{255, 0, 0, 255}),
		SetEndColor(color.RGBA{0, 0, 255, 255}),
		SetBounds(b),
	)
}


// Conic Gradient Pattern

var ConicGradientOutputFilename = "conic_gradient.png"
var ConicGradientZoomLevels = []int{}
const ConicGradientOrder = 32

func ExampleNewConicGradient() {
	// Conic Gradient
	NewConicGradient(
		SetStartColor(color.RGBA{255, 0, 255, 255}),
		SetEndColor(color.RGBA{0, 255, 255, 255}),
	)
}

func GenerateConicGradient(b image.Rectangle) image.Image {
	return NewConicGradient(
		SetStartColor(color.RGBA{255, 0, 255, 255}),
		SetEndColor(color.RGBA{0, 255, 255, 255}),
		SetBounds(b),
	)
}

func init() {
	RegisterGenerator("LinearGradient", GenerateLinearGradient)
	RegisterReferences("LinearGradient", GenerateLinearGradientReferences)

	RegisterGenerator("RadialGradient", GenerateRadialGradient)

	RegisterGenerator("ConicGradient", GenerateConicGradient)
}
