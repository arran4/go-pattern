package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var SierpinskiTriangleOutputFilename = "sierpinski_triangle.png"
var SierpinskiTriangleZoomLevels = []int{}
const SierpinskiTriangleOrder = 40

// Sierpinski Triangle
// Generates a Sierpinski Triangle fractal (right-angled variant using Pascal's Triangle modulo 2).
func ExampleNewSierpinskiTriangle() {
	b := image.Rect(0, 0, 150, 150)
	i := NewSierpinskiTriangle(SetBounds(b), SetFillColor(color.Black), SetSpaceColor(color.White))
	f, err := os.Create(SierpinskiTriangleOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateSierpinskiTriangle(b image.Rectangle) image.Image {
	return NewSierpinskiTriangle(SetBounds(b), SetFillColor(color.Black), SetSpaceColor(color.White))
}

var SierpinskiCarpetOutputFilename = "sierpinski_carpet.png"
var SierpinskiCarpetZoomLevels = []int{}
const SierpinskiCarpetOrder = 41

// Sierpinski Carpet
// Generates a Sierpinski Carpet fractal.
func ExampleNewSierpinskiCarpet() {
	b := image.Rect(0, 0, 150, 150)
	i := NewSierpinskiCarpet(SetBounds(b), SetFillColor(color.Black), SetSpaceColor(color.White))
	f, err := os.Create(SierpinskiCarpetOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateSierpinskiCarpet(b image.Rectangle) image.Image {
	return NewSierpinskiCarpet(SetBounds(b), SetFillColor(color.Black), SetSpaceColor(color.White))
}

func init() {
	RegisterGenerator("SierpinskiTriangle", GenerateSierpinskiTriangle)
	RegisterGenerator("SierpinskiCarpet", GenerateSierpinskiCarpet)
}
