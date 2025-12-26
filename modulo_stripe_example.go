package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var ModuloStripeOutputFilename = "modulo_stripe.png"
var ModuloStripeZoomLevels = []int{}
const ModuloStripeOrder = 31

func ExampleNewModuloStripe() {
	p := NewModuloStripe([]color.Color{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
	})
	f, err := os.Create(ModuloStripeOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, p); err != nil {
		panic(err)
	}
}

func GenerateModuloStripe(b image.Rectangle) image.Image {
	return NewModuloStripe([]color.Color{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
	}, SetBounds(b))
}

func GenerateModuloStripeReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("ModuloStripe", GenerateModuloStripe)
	RegisterReferences("ModuloStripe", GenerateModuloStripeReferences)
}
