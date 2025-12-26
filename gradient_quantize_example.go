package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var GradientQuantizationOutputFilename = "gradient_quantize.png"
var GradientQuantizationZoomLevels = []int{}
const GradientQuantizationOrder = 36

func ExampleNewGradientQuantization() {
	grad := NewLinearGradient(
		SetStartColor(color.Black),
		SetEndColor(color.White),
	)
	p := NewQuantize(grad, 4)
	f, err := os.Create(GradientQuantizationOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, p); err != nil {
		panic(err)
	}
}

func GenerateGradientQuantization(b image.Rectangle) image.Image {
	grad := NewLinearGradient(
		SetStartColor(color.Black),
		SetEndColor(color.White),
		SetBounds(b),
	)
	return NewQuantize(grad, 4, SetBounds(b))
}

func GenerateGradientQuantizationReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("GradientQuantization", GenerateGradientQuantization)
	RegisterReferences("GradientQuantization", GenerateGradientQuantizationReferences)
}
