package pattern

import (
	"image"
)

var (
	AsphaltOutputFilename = "asphalt.png"
)

func ExampleNewAsphalt() image.Image {
	return NewAsphalt()
}

func GenerateAsphalt(rect image.Rectangle) image.Image {
	return ExampleNewAsphalt()
}

func init() {
	GlobalGenerators["Asphalt"] = GenerateAsphalt
}
