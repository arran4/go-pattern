package pattern

import (
	"image"
	"image/png"
	"os"
)

var GopherOutputFilename = "gopher.png"
var GopherZoomLevels = []int{}

const GopherOrder = 20
const GopherBaseLabel = "Gopher"

// Gopher Pattern
// A static image of the Go Gopher.
func ExampleNewGopher() {
	i := NewGopher()
	f, err := os.Create(GopherOutputFilename)
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

func BootstrapGopher(b image.Rectangle) image.Image {
	return NewGopher()
}

func BootstrapGopherReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Gopher": func(b image.Rectangle) image.Image {
			return NewGopher()
		},
	}, []string{"Gopher"}
}

func init() {
	RegisterGenerator("Gopher", BootstrapGopher)
	RegisterReferences("Gopher", BootstrapGopherReferences)
}
