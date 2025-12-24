package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var BayerDitherOutputFilename = "bayer_dither.png"
var BayerDitherZoomLevels = []int{}
const BayerDitherOrder = 34

func ExampleNewBayerDither() {
	grad := NewLinearGradient(
		SetStartColor(color.Black),
		SetEndColor(color.White),
	)
	p := NewBayerDither(grad, 4)
	f, err := os.Create(BayerDitherOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, p)
}

func GenerateBayerDither(b image.Rectangle) image.Image {
	grad := NewLinearGradient(
		SetStartColor(color.Black),
		SetEndColor(color.White),
		SetBounds(b),
	)
	return NewBayerDither(grad, 4, SetBounds(b))
}

func GenerateBayerDitherReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("BayerDither", GenerateBayerDither)
	RegisterReferences("BayerDither", GenerateBayerDitherReferences)
}
