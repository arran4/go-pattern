package pattern

import (
	"image"
	"image/png"
	"os"
)

var XorGridOutputFilename = "xor_pattern.png"
var XorGridZoomLevels = []int{}
const XorGridOrder = 30

func ExampleNewXorGrid() {
	p := NewXorPattern()
	f, err := os.Create(XorGridOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, p)
}

func GenerateXorGrid(b image.Rectangle) image.Image {
	return NewXorPattern(SetBounds(b))
}

func GenerateXorGridReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("XorGrid", GenerateXorGrid)
	RegisterReferences("XorGrid", GenerateXorGridReferences)
}
