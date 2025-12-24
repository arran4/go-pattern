package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var HorizontalLineOutputFilename = "horizontal_line.png"
var HorizontalLineZoomLevels = []int{} // Disabled as per feedback

const HorizontalLineOrder = 10

func ExampleNewHorizontalLine() {
	i := NewHorizontalLine(
		SetLineSize(5),
		SetSpaceSize(5),
		SetLineColor(color.RGBA{255, 0, 0, 255}),
		SetSpaceColor(color.White),
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
	// Show variations stitched together instead of zooms
	v1 := NewDemoHorizontalLine(
		SetBounds(b),
		SetLineSize(10),
		SetSpaceSize(10),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
	)
	v2 := NewDemoHorizontalLine(
		SetBounds(b),
		SetLineSize(2),
		SetSpaceSize(18),
		SetLineColor(color.RGBA{255, 0, 0, 255}), // Red thin lines
		SetSpaceColor(color.White),
	)
	v3 := NewDemoHorizontalLine(
		SetBounds(b),
		SetLineSize(18),
		SetSpaceSize(2),
		SetLineColor(color.RGBA{0, 0, 255, 255}), // Blue thick lines
		SetSpaceColor(color.White),
	)

	return stitchImagesForDemo(v1, v2, v3)
}

var VerticalLineOutputFilename = "vertical_line.png"
var VerticalLineZoomLevels = []int{}

const VerticalLineOrder = 11

func ExampleNewVerticalLine() {
	i := NewVerticalLine(
		SetLineSize(5),
		SetSpaceSize(5),
		SetLineColor(color.RGBA{0, 0, 255, 255}),
		SetSpaceColor(color.White),
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
	v1 := NewDemoVerticalLine(
		SetBounds(b),
		SetLineSize(10),
		SetSpaceSize(10),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
	)
	v2 := NewDemoVerticalLine(
		SetBounds(b),
		SetLineSize(5),
		SetSpaceSize(15),
		SetLineColor(color.RGBA{0, 128, 0, 255}), // Green
		SetSpaceColor(color.White),
	)
	return stitchImagesForDemo(v1, v2)
}

// stitchImagesForDemo is defined in boolean_example.go (shared via package scope)

func init() {
	RegisterGenerator("HorizontalLine", GenerateHorizontalLine)
	RegisterGenerator("VerticalLine", GenerateVerticalLine)
}
