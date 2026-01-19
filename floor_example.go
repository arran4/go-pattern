package pattern

import (
	"image"
	"image/color"
)

var (
	FloorOutputFilename      = "floor.png"
	Floor_woodOutputFilename = "floor_wood.png"
)

func ExampleNewFloor() image.Image {
	// Tiled floor using Tile pattern?
	// We have `NewTile` which tiles an image.
	// We have `NewBrick` or `NewGrid`?
	// `brick.go` makes bricks.
	// `checker.go` makes checks.

	// Let's use `NewBrick` for a tile floor.
	// Large square tiles.

	// Create a marble texture for tiles
	marble := NewNoise(SetNoiseAlgorithm(&PerlinNoise{Frequency: 0.1}))
	marbleColor := NewColorMap(marble,
		ColorStop{0.0, color.RGBA{220, 220, 220, 255}},
		ColorStop{1.0, color.White},
	)

	// Create a slightly different marble for variation
	marble2 := NewRotate(marbleColor, 90)

	mortarColor := NewRect(SetFillColor(color.RGBA{50, 50, 50, 255}))

	return NewBrick(
		SetBrickSize(60, 60),
		SetMortarSize(3),
		SetBrickOffset(0),
		SetBrickImages(marbleColor, marble2),
		SetMortarImage(mortarColor),
	)
}

func GenerateFloor(rect image.Rectangle) image.Image {
	return ExampleNewFloor()
}

func init() {
	GlobalGenerators["Floor"] = GenerateFloor
}
