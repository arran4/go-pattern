package pattern

import (
	"image"
)

var (
	Brick_weatheredOutputFilename = "brick_weathered.png"
	Brick_weatheredZoomLevels     = []int{}
	Brick_weatheredBaseLabel      = "Brick Weathered"
)

func init() {
	RegisterGenerator(Brick_weatheredBaseLabel, GenerateBrickWeathered)
}

// GenerateBrickWeathered builds a chipped brick wall example with hue variation and recessed mortar.
func GenerateBrickWeathered(bounds image.Rectangle) image.Image {
	return NewChippedBrick(
		SetBounds(bounds),
		SetBrickSize(48, 22),
		SetMortarSize(4),
		SetChipIntensity(0.45),
		SetMortarDepth(0.8),
		SetHueJitter(0.18),
		SetSeed(2024),
	)
}

// ExampleNewBrick_weathered provides a sample for documentation use.
func ExampleNewBrick_weathered() image.Image {
	return GenerateBrickWeathered(image.Rect(0, 0, 300, 300))
}
