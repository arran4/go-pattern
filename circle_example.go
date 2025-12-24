package pattern

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var CircleOutputFilename = "circle.png"
var CircleZoomLevels = []int{}
const CircleOrder = 24

const CircleBaseLabel = "Circle"

func ExampleNewCircle() {
	// Create a simple circle
	c := NewCircle(SetLineColor(color.Black), SetSpaceColor(color.White))
	fmt.Printf("Circle bounds: %v\n", c.Bounds())
	// Output:
	// Circle bounds: (0,0)-(255,255)

	f, err := os.Create(CircleOutputFilename)
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

func GenerateCircle(b image.Rectangle) image.Image {
	v1 := NewCircle(
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
		SetBounds(b),
	)
	v2 := NewCircle(
		SetLineColor(color.RGBA{255, 0, 0, 255}),
		SetSpaceColor(color.RGBA{255, 255, 0, 255}),
		SetBounds(b),
	)
	return stitchImagesForDemo(v1, v2)
}

func GenerateCircleReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"RedCircle": func(b image.Rectangle) image.Image {
			return NewCircle(
				SetLineColor(color.RGBA{255, 0, 0, 255}),
				SetSpaceColor(color.RGBA{255, 255, 255, 255}),
				SetBounds(b),
			)
		},
		"TransparentBackground": func(b image.Rectangle) image.Image {
			return NewCircle(
				SetLineColor(color.Black),
				// No SpaceColor set, defaults to transparent
				SetBounds(b),
			)
		},
	}, []string{"RedCircle", "TransparentBackground"}
}

func init() {
	RegisterGenerator("Circle", GenerateCircle)
	RegisterReferences("Circle", GenerateCircleReferences)
}
