package pattern

import (
	"image"
	"image/png"
	"os"
)

var XorPatternOutputFilename = "xor_pattern.png"
var XorPatternZoomLevels = []int{}
const XorPatternOrder = 30

func ExampleNewXorPattern() {
	p := NewXorPattern()
	f, err := os.Create(XorPatternOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, p)
}

func GenerateXorPattern(b image.Rectangle) image.Image {
	return NewXorPattern(SetBounds(b))
}

func GenerateXorPatternReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("XorPattern", GenerateXorPattern)
	RegisterReferences("XorPattern", GenerateXorPatternReferences)
}
