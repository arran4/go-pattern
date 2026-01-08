package pattern

import (
	"image"
)

var (
	ChippedBrickOutputFilename = "brick_weathered.png"
	ChippedBrickZoomLevels     = []int{}
	ChippedBrickBaseLabel      = "Chipped Brick"
)

func init() {
	RegisterGenerator("ChippedBrick", GenerateChippedBrick)
}

// GenerateChippedBrick builds a chipped brick wall example with hue variation and recessed mortar.
func GenerateChippedBrick(bounds image.Rectangle) image.Image {
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

// ExampleNewChippedBrick provides a sample for documentation use.
func ExampleNewChippedBrick() image.Image {
	return GenerateChippedBrick(image.Rect(0, 0, 300, 300))
}
