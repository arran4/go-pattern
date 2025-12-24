package pattern

import (
	"image"
	"image/png"
	"os"
)

var QuantizeOutputFilename = "quantize.png"
var QuantizeZoomLevels = []int{}

const QuantizeOrder = 30
const QuantizeBaseLabel = "Quantize"

// Quantize Pattern
// Example of quantizing the colors of an image (Posterization).
// This example reduces the Gopher image to 4 levels per channel.
func ExampleNewQuantize() {
	i := NewQuantize(NewGopher(), 4)
	f, err := os.Create(QuantizeOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateQuantize(b image.Rectangle) image.Image {
	return NewQuantize(NewGopher(), 4)
}

func GenerateQuantizeReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Quantize (2 levels)": func(b image.Rectangle) image.Image {
			return NewQuantize(NewGopher(), 2)
		},
		"Quantize (4 levels)": func(b image.Rectangle) image.Image {
			return NewQuantize(NewGopher(), 4)
		},
		"Quantize (8 levels)": func(b image.Rectangle) image.Image {
			return NewQuantize(NewGopher(), 8)
		},
	}, []string{"Quantize (2 levels)", "Quantize (4 levels)", "Quantize (8 levels)"}
}

func init() {
	RegisterGenerator("Quantize", GenerateQuantize)
	RegisterReferences("Quantize", GenerateQuantizeReferences)
}
