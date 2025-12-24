package pattern

import (
	"image"
	"image/color"
)

var SierpinskiTriangleOutputFilename = "sierpinski_triangle.png"
var SierpinskiTriangleZoomLevels = []int{}
const SierpinskiTriangleOrder = 40

// Sierpinski Triangle
// Generates a Sierpinski Triangle fractal (right-angled variant using Pascal's Triangle modulo 2).
func ExampleNewSierpinskiTriangle() {
	// See GenerateSierpinskiTriangle for implementation details
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
	// See GenerateSierpinskiCarpet for implementation details
}

func GenerateSierpinskiCarpet(b image.Rectangle) image.Image {
	return NewSierpinskiCarpet(SetBounds(b), SetFillColor(color.Black), SetSpaceColor(color.White))
}

func init() {
	RegisterGenerator("SierpinskiTriangle", GenerateSierpinskiTriangle)
	RegisterGenerator("SierpinskiCarpet", GenerateSierpinskiCarpet)
}
