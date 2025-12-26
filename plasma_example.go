package pattern

import (
	"image"
	"image/png"
	"os"
)

var PlasmaOutputFilename = "plasma.png"
var PlasmaZoomLevels = []int{}
const PlasmaOrder = 33

func ExampleNewPlasma() {
	p := NewPlasma()
	f, err := os.Create(PlasmaOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, p); err != nil {
		panic(err)
	}
}

func GeneratePlasma(b image.Rectangle) image.Image {
	return NewPlasma(SetBounds(b))
}

func GeneratePlasmaReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("Plasma", GeneratePlasma)
	RegisterReferences("Plasma", GeneratePlasmaReferences)
}
