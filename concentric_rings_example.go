package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var ConcentricRingsOutputFilename = "concentric_rings.png"
var ConcentricRingsZoomLevels = []int{}
const ConcentricRingsOrder = 32

func ExampleNewConcentricRings() {
	p := NewConcentricRings([]color.Color{
		color.Black,
		color.White,
		color.RGBA{255, 0, 0, 255},
	})
	f, err := os.Create(ConcentricRingsOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, p); err != nil {
		panic(err)
	}
}

func GenerateConcentricRings(b image.Rectangle) image.Image {
	return NewConcentricRings([]color.Color{
		color.Black,
		color.White,
		color.RGBA{255, 0, 0, 255},
	}, SetBounds(b))
}

func GenerateConcentricRingsReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("ConcentricRings", GenerateConcentricRings)
	RegisterReferences("ConcentricRings", GenerateConcentricRingsReferences)
}
