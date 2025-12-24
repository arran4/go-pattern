package pattern

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var SierpinskiOutputFilename = "sierpinski.png"
var SierpinskiZoomLevels = []int{}
const SierpinskiOrder = 25

const SierpinskiBaseLabel = "Sierpinski"

func ExampleNewSierpinski() {
	// Create a simple Sierpiński triangle
	c := NewSierpinski(SetFillColor(color.Black), SetSpaceColor(color.White))
	fmt.Printf("Sierpinski bounds: %v\n", c.Bounds())
	// Output:
	// Sierpinski bounds: (0,0)-(255,255)

	f, err := os.Create(SierpinskiOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, c); err != nil {
		panic(err)
	}
}

func GenerateSierpinski(b image.Rectangle) image.Image {
	v1 := NewSierpinski(
		SetFillColor(color.Black),
		SetSpaceColor(color.White),
		SetBounds(b),
	)
	v2 := NewSierpinski(
		SetFillColor(color.RGBA{255, 0, 0, 255}), // Red
		SetSpaceColor(color.RGBA{0, 0, 0, 255}), // Black Background
		SetBounds(b),
	)
	v3 := NewSierpinski(
		SetFillColor(color.White),
		SetSpaceColor(color.RGBA{0, 0, 255, 255}), // Blue Background
		SetBounds(b),
	)
	v4 := NewSierpinski(
		SetFillColor(color.RGBA{0, 255, 0, 255}), // Green
		SetSpaceColor(color.Transparent),
		SetBounds(b),
	)

	return stitchImagesForDemo(v1, v2, v3, v4)
}

func GenerateSierpinskiReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"RedBlack": func(b image.Rectangle) image.Image {
			return NewSierpinski(
				SetFillColor(color.RGBA{255, 0, 0, 255}),
				SetSpaceColor(color.Black),
				SetBounds(b),
			)
		},
		"GreenTransparent": func(b image.Rectangle) image.Image {
			return NewSierpinski(
				SetFillColor(color.RGBA{0, 255, 0, 255}),
				SetSpaceColor(color.Transparent),
				SetBounds(b),
			)
		},
	}, []string{"RedBlack", "GreenTransparent"}
}

func init() {
	RegisterGenerator("Sierpinski", GenerateSierpinski)
	RegisterReferences("Sierpinski", GenerateSierpinskiReferences)
}
