package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var PolkaOutputFilename = "polka.png"

// Polka Pattern
// A pattern of dots (circles) arranged in a grid.
func ExampleNewPolka() {
	i := NewPolka(
		SetRadius(10),
		SetSpacing(40),
		SetFillColor(color.Black),
		SetSpaceColor(color.White),
	)
	f, err := os.Create(PolkaOutputFilename)
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

func GeneratePolka(b image.Rectangle) image.Image {
	return NewDemoPolka(SetBounds(b))
}

func init() {
	RegisterGenerator("Polka", GeneratePolka)
}
