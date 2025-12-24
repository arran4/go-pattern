package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var MirrorOutputFilename = "mirror.png"
var MirrorZoomLevels = []int{2, 4}

const MirrorOrder = 32

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
	return NewText("Go", TextSize(40), TextColorColor(color.Black))
}

func GenerateMirror(b image.Rectangle) image.Image {
	input := NewDemoMirrorInput(b)
	// We want to show Original, Mirror H, Mirror V, Mirror HV
	// Create a 2x2 grid.
	// Since Grid expects components, we can build it.

	// Helper to add bounds to an image if needed, though NewGrid handles layout.
	// But Mirror needs to know the bounds of the input to mirror correctly?
	// Mirror uses img.Bounds().
	// Text image has fixed bounds (0,0, w, h).

	return NewGrid(
		2, 2,
		NewDemoMirrorInput(b), // Original
		NewMirror(input, true, false), // Mirror H
		NewMirror(input, false, true), // Mirror V
		NewMirror(input, true, true), // Mirror HV
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
