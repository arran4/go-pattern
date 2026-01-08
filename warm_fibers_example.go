package pattern

import (
	"image"
	"image/color"
)

// Warm Fibers Gradient

var WarmFibersGradientOutputFilename = "warm_fibers.png"
var WarmFibersGradientZoomLevels = []int{}

const WarmFibersGradientOrder = 105
const WarmFibersGradientBaseLabel = "Base"

func ExampleNewWarmFibersGradient() {
	NewWarmFibersGradient(
		SetStartColor(color.RGBA{245, 200, 160, 255}),
		SetEndColor(color.RGBA{130, 70, 40, 255}),
		SetFiberDensity(1.25),
		SetVignetteStrength(0.4),
	)
}

func GenerateWarmFibersGradient(b image.Rectangle) image.Image {
	return NewWarmFibersGradient(
		SetStartColor(color.RGBA{245, 200, 160, 255}),
		SetEndColor(color.RGBA{130, 70, 40, 255}),
		SetFiberDensity(1.25),
		SetVignetteStrength(0.4),
		SetBounds(b),
	)
}

func GenerateWarmFibersReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Fine Fibers": func(b image.Rectangle) image.Image {
			return NewWarmFibersGradient(
				SetStartColor(color.RGBA{240, 192, 150, 255}),
				SetEndColor(color.RGBA{143, 72, 45, 255}),
				SetFiberDensity(0.7),
				SetVignetteStrength(0.25),
				SetBounds(b),
			)
		},
		"Strong Vignette": func(b image.Rectangle) image.Image {
			return NewWarmFibersGradient(
				SetStartColor(color.RGBA{250, 205, 170, 255}),
				SetEndColor(color.RGBA{120, 60, 35, 255}),
				SetFiberDensity(1.1),
				SetVignetteStrength(0.65),
				SetBounds(b),
			)
		},
	}, []string{"Fine Fibers", "Strong Vignette"}
}

func init() {
	RegisterGenerator("WarmFibersGradient", GenerateWarmFibersGradient)
	RegisterReferences("WarmFibersGradient", GenerateWarmFibersReferences)
}
