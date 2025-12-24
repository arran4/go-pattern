package pattern

import (
	"image"
	"image/color"
)

const GradientBaseLabel = "Gradient"

func init() {
	// Register the generators
	// This makes them available for the bootstrap tool to generate readme.md
}

func GenerateLinearGradient(ops ...func(any)) image.Image {
	return NewLinearGradient(ops...)
}

func GenerateRadialGradient(ops ...func(any)) image.Image {
	return NewRadialGradient(ops...)
}

func GenerateConicGradient(ops ...func(any)) image.Image {
	return NewConicGradient(ops...)
}

func ExampleNewLinearGradient() {
	// Linear Gradient (Horizontal)
	NewLinearGradient(
		SetStartColor(color.RGBA{255, 0, 0, 255}),
		SetEndColor(color.RGBA{0, 0, 255, 255}),
	)
}

func ExampleNewLinearGradient_vertical() {
	// Linear Gradient (Vertical)
	NewLinearGradient(
		SetStartColor(color.RGBA{0, 255, 0, 255}),
		SetEndColor(color.RGBA{255, 255, 0, 255}),
		GradientVertical(),
	)
}

func ExampleNewRadialGradient() {
	// Radial Gradient
	NewRadialGradient(
		SetStartColor(color.RGBA{255, 0, 0, 255}),
		SetEndColor(color.RGBA{0, 0, 255, 255}),
	)
}

func ExampleNewConicGradient() {
	// Conic Gradient
	NewConicGradient(
		SetStartColor(color.RGBA{255, 0, 255, 255}),
		SetEndColor(color.RGBA{0, 255, 255, 255}),
	)
}
