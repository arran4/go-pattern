package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var MirrorOutputFilename = "mirror.png"
var MirrorZoomLevels = []int{}

const MirrorOrder = 32
const MirrorBaseLabel = "Grid"

// Mirror Pattern
// Mirrors the input pattern horizontally or vertically.
func ExampleNewMirror() {
	i := NewMirror(NewDemoMirrorInput(image.Rect(0, 0, 40, 40)), true, false)
	f, err := os.Create(MirrorOutputFilename)
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

func NewDemoMirrorInput(b image.Rectangle) image.Image {
	// Create an asymmetric image to demonstrate mirroring.
	// Use white background to ensure visibility.
	return NewText("Go", TextSize(30), TextColorColor(color.Black), TextBackgroundColorColor(color.White))
}

func GenerateMirror(b image.Rectangle) image.Image {
	input := NewDemoMirrorInput(b)
	// We want to show Original, Mirror H, Mirror V, Mirror HV
	// Create a 2x2 grid.
	return NewGrid(
		Row(input, NewMirror(input, true, false)),
		Row(NewMirror(input, false, true), NewMirror(input, true, true)),
	)
}

func GenerateMirrorReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Input": NewDemoMirrorInput,
	}, []string{"Input"}
}

func init() {
	RegisterGenerator("Mirror", GenerateMirror)
	RegisterReferences("Mirror", GenerateMirrorReferences)
}
