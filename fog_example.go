package pattern

import (
	"image"
	"image/color"
)

var (
	FogOutputFilename = "fog.png"
	FogZoomLevels     = []int{}
	FogBaseLabel      = "Fog"
)

const FogOrder = 105

// ExampleNewFog renders soft Perlin/fBm fog with a radial falloff so the center stays clearer.
func ExampleNewFog() image.Image {
	return NewFog(
		SetDensity(0.85),
		SetFalloffCurve(1.8),
		SetFillColor(color.RGBA{185, 205, 230, 255}),
	)
}

func GenerateFog(b image.Rectangle) image.Image {
	return NewFog(
		SetBounds(b),
		SetDensity(0.9),
		SetFalloffCurve(2.0),
		SetFillColor(color.RGBA{190, 205, 230, 255}),
	)
}

func GenerateFogReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"LightMist": func(b image.Rectangle) image.Image {
			return NewFog(
				SetBounds(b),
				SetDensity(0.45),
				SetFalloffCurve(1.2),
				SetFillColor(color.RGBA{215, 225, 235, 255}),
			)
		},
		"ColdBlue": func(b image.Rectangle) image.Image {
			return NewFog(
				SetBounds(b),
				SetDensity(0.75),
				SetFalloffCurve(1.6),
				SetFillColor(color.RGBA{160, 190, 235, 255}),
			)
		},
		"EmberGlow": func(b image.Rectangle) image.Image {
			return NewFog(
				SetBounds(b),
				SetDensity(1.1),
				SetFalloffCurve(2.4),
				SetFillColor(color.RGBA{255, 185, 150, 255}),
			)
		},
	}, []string{"LightMist", "ColdBlue", "EmberGlow"}
}

func init() {
	RegisterGenerator(FogBaseLabel, GenerateFog)
	RegisterReferences(FogBaseLabel, GenerateFogReferences)
}
