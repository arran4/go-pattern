package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var RotateOutputFilename = "rotate.png"
var RotateZoomLevels = []int{}

const RotateOrder = 33

// Rotate Pattern
// Rotates the input pattern by 90, 180, or 270 degrees.
func ExampleNewRotate() {
	i := NewRotate(NewDemoRotateInput(image.Rect(0, 0, 40, 60)), 90)
	f, err := os.Create(RotateOutputFilename)
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

func NewDemoRotateInput(b image.Rectangle) image.Image {
	// Asymmetric input: Width != Height to see rotation effect on bounds.
	// Use white background.
	return NewText("Go", TextSize(80), TextColorColor(color.Black), TextBackgroundColorColor(color.White))
}

func GenerateRotate(b image.Rectangle) image.Image {
	input := NewDemoRotateInput(b)
	return NewGrid(
		Row(input, NewRotate(input, 90)),
		Row(NewRotate(input, 180), NewRotate(input, 270)),
	)
}

func GenerateRotateReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Input": NewDemoRotateInput,
	}, []string{"Input"}
}

func init() {
	RegisterGenerator("Rotate", GenerateRotate)
	RegisterReferences("Rotate", GenerateRotateReferences)
}
