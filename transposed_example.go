package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var TransposedOutputFilename = "transposed.png"
var TransposedZoomLevels = []int{2, 4}

const TransposedOrder = 31

// Transposed Pattern
// Transposes the coordinates of an input pattern.
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

func demoTransposedInput(b image.Rectangle) image.Image {
	// The original demo used NewSimpleZoom(NewChecker..., 20)
	// Explicitly set SpaceSize to 1 to match legacy 1x1 checker behavior expected by this demo.
	return NewSimpleZoom(NewChecker(color.Black, color.White, SetBounds(b), SetSpaceSize(1)), 20, SetBounds(b))
}

func GenerateTransposed(b image.Rectangle) image.Image {
	return NewTransposed(demoTransposedInput(b), 10, 10, SetBounds(b))
}

func GenerateTransposedReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Input": demoTransposedInput,
	}, []string{"Input"}
}

func init() {
	RegisterGenerator("Transposed", GenerateTransposed)
	RegisterReferences("Transposed", GenerateTransposedReferences)
}
