package pattern

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var FibonacciOutputFilename = "fibonacci.png"
var FibonacciZoomLevels = []int{}

const FibonacciOrder = 25 // Adjust as needed to fit in the list

const FibonacciBaseLabel = "Fibonacci"

func ExampleNewFibonacci() {
	// Create a simple Fibonacci spiral
	c := NewFibonacci(SetLineColor(color.Black), SetSpaceColor(color.White))
	fmt.Printf("Fibonacci bounds: %v\n", c.Bounds())
	// Output:
	// Fibonacci bounds: (0,0)-(255,255)

	f, err := os.Create(FibonacciOutputFilename)
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

func GenerateFibonacci(b image.Rectangle) image.Image {
	v1 := NewFibonacci(
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
		SetBounds(b),
	)
	v2 := NewFibonacci(
		SetLineSize(5),
		SetLineColor(color.RGBA{0, 0, 255, 255}),
		SetSpaceColor(color.White),
		SetBounds(b),
	)
	v3 := NewFibonacci(
		SetLineSize(2),
		SetLineColor(color.White),
		SetSpaceColor(color.Black),
		SetBounds(b),
	)
	v4 := NewFibonacci(
		SetLineSize(10),
		SetLineColor(color.RGBA{255, 0, 0, 255}),
		SetSpaceColor(color.RGBA{255, 255, 0, 100}), // Transparent Yellow
		SetBounds(b),
	)

	return stitchImagesForDemo(v1, v2, v3, v4)
}

func GenerateFibonacciReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"BlueSpiral": func(b image.Rectangle) image.Image {
			return NewFibonacci(
				SetLineSize(3),
				SetLineColor(color.RGBA{0, 0, 255, 255}),
				SetSpaceColor(color.White),
				SetBounds(b),
			)
		},
		"ThinBlack": func(b image.Rectangle) image.Image {
			return NewFibonacci(
				SetLineSize(1),
				SetLineColor(color.Black),
				SetSpaceColor(color.White),
				SetBounds(b),
			)
		},
		"TransparentBackground": func(b image.Rectangle) image.Image {
			return NewFibonacci(
				SetLineSize(2),
				SetLineColor(color.Black),
				// No SpaceColor set, defaults to transparent
				SetBounds(b),
			)
		},
	}, []string{"ThinBlack", "BlueSpiral", "TransparentBackground"}
}

func init() {
	RegisterGenerator("Fibonacci", GenerateFibonacci)
	RegisterReferences("Fibonacci", GenerateFibonacciReferences)
}
