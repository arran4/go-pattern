package pattern

import (
	"image"
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
	return NewSimpleZoom(NewGopher(), 2)
}

func GenerateMirror(b image.Rectangle) image.Image {
	return NewMirror(NewDemoMirrorInput(b), true, false, SetBounds(b))
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
