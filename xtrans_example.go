package pattern

import (
	"image"
	"image/png"
	"os"
)

var XTransOutputFilename = "xtrans.png"
var XTransZoomLevels = []int{}

const XTransOrder = 34
const XTransBaseLabel = "XTrans"

// ExampleNewXTrans generates an example of applying an X-Trans CFA pattern.
func ExampleNewXTrans() {
	// Create a linear gradient input
	grad := NewLinearGradient()

	// Apply X-Trans pattern
	p := NewXTrans(grad)

	f, err := os.Create(XTransOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()

	if err = png.Encode(f, p); err != nil {
		panic(err)
	}
}

func GenerateXTrans(b image.Rectangle) image.Image {
	grad := NewLinearGradient(SetBounds(b))
	return NewXTrans(grad, SetBounds(b))
}

func GenerateXTransReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"X-Trans Filtered Gradient": func(b image.Rectangle) image.Image {
			grad := NewLinearGradient(SetBounds(b))
			return NewXTrans(grad, SetBounds(b))
		},
		"X-Trans CFA": func(b image.Rectangle) image.Image {
			return NewXTrans(nil, SetBounds(b))
		},
	}, []string{"X-Trans Filtered Gradient", "X-Trans CFA"}
}

func init() {
	RegisterGenerator("XTrans", GenerateXTrans)
	RegisterReferences("XTrans", GenerateXTransReferences)
}
