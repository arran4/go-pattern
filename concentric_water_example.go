package pattern

import (
	"image"
	"image/color"
)

var ConcentricWaterOutputFilename = "concentric_water.png"
var ConcentricWaterZoomLevels = []int{}

const ConcentricWaterOrder = 370

// ExampleNewConcentricWater demonstrates concentric distance-field ripples with
// sine-driven heights that tint and bend the normals of the surface.
func ExampleNewConcentricWater() image.Image {
	return NewConcentricWater(
		ConcentricWaterRingSpacing(14.0),
		ConcentricWaterAmplitude(1.1),
		ConcentricWaterAmplitudeFalloff(0.018),
		ConcentricWaterBaseTint(color.RGBA{24, 104, 168, 255}),
		ConcentricWaterNormalStrength(4.0),
	)
}

func GenerateConcentricWater(b image.Rectangle) image.Image {
	cx := b.Dx() / 2
	cy := b.Dy() / 2
	return NewConcentricWater(
		ConcentricWaterRingSpacing(14.0),
		ConcentricWaterAmplitude(1.1),
		ConcentricWaterAmplitudeFalloff(0.018),
		ConcentricWaterBaseTint(color.RGBA{24, 104, 168, 255}),
		ConcentricWaterNormalStrength(4.0),
		SetBounds(b),
		SetCenter(cx, cy),
	)
}

func GenerateConcentricWaterReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("ConcentricWater", GenerateConcentricWater)
	RegisterReferences("ConcentricWater", GenerateConcentricWaterReferences)
}
