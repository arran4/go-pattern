package pattern

import (
	"image"
	"image/png"
	"os"
)

var TransposedOutputFilename = "transposed.png"
var TransposedZoomLevels = []int{}

const TransposedOrder = 3
const TransposedBaseLabel = "Transposed"

// Transposed Pattern
// Transposes the X and Y coordinates of an underlying image.
func ExampleNewTransposed() {
	i := NewTransposed(NewDemoNull(), 10, 10)
	f, err := os.Create(TransposedOutputFilename)
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

func BootstrapTransposed(b image.Rectangle) image.Image {
	return nil
}

func BootstrapTransposedReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Original": func(b image.Rectangle) image.Image {
			return NewSimpleZoom(NewDemoChecker(SetBounds(b)), 10, SetBounds(b))
		},
		"Transposed": func(b image.Rectangle) image.Image {
			base := NewSimpleZoom(NewDemoChecker(SetBounds(b)), 10, SetBounds(b))
			return NewTransposed(base, 5, 5, SetBounds(b))
		},
	}, []string{"Original", "Transposed"}
}

func init() {
	RegisterGenerator("Transposed", BootstrapTransposed)
	RegisterReferences("Transposed", BootstrapTransposedReferences)
}
