package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var HorizontalLineOutputFilename = "horizontal_line.png"
var HorizontalLineZoomLevels = []int{2, 4}

const HorizontalLineOrder = 10

func ExampleNewHorizontalLine() {
	i := NewHorizontalLine(
		SetLineSize(5),
		SetSpaceSize(5),
		SetLineColor(color.RGBA{255, 0, 0, 255}),
	)
	f, err := os.Create(HorizontalLineOutputFilename)
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

func GenerateHorizontalLine(b image.Rectangle) image.Image {
	return NewDemoHorizontalLine(
		SetBounds(b),
		SetLineSize(10),
		SetSpaceSize(10),
		SetLineColor(color.Black), // More visible than default? Default is Black.
	)
}

var VerticalLineOutputFilename = "vertical_line.png"
var VerticalLineZoomLevels = []int{2, 4}

const VerticalLineOrder = 11

func ExampleNewVerticalLine() {
	i := NewVerticalLine(
		SetLineSize(5),
		SetSpaceSize(5),
		SetLineColor(color.RGBA{0, 0, 255, 255}),
	)
	f, err := os.Create(VerticalLineOutputFilename)
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

func GenerateVerticalLine(b image.Rectangle) image.Image {
	return NewDemoVerticalLine(
		SetBounds(b),
		SetLineSize(10),
		SetSpaceSize(10),
		SetLineColor(color.Black),
	)
}

func init() {
	RegisterGenerator("HorizontalLine", GenerateHorizontalLine)
	RegisterGenerator("VerticalLine", GenerateVerticalLine)
}
