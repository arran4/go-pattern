package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var ScreenToneOutputFilename = "screentone.png"

// ScreenTone Pattern
// A halftone dot matrix pattern with adjustable density (Spacing) and angle.
func ExampleNewScreenTone() {
	i := NewScreenTone(
		SetRadius(3),
		SetSpacing(10),
		SetAngle(45),
		SetFillColor(color.Black),
		SetSpaceColor(color.White),
	)
	f, err := os.Create(ScreenToneOutputFilename)
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

func GenerateScreenTone(b image.Rectangle) image.Image {
	v1 := NewScreenTone(
		SetBounds(b),
		SetRadius(3),
		SetSpacing(10),
		SetAngle(45),
		SetFillColor(color.Black),
		SetSpaceColor(color.White),
	)

	v2 := NewScreenTone(
		SetBounds(b),
		SetRadius(3),
		SetSpacing(10),
		SetAngle(15),
		SetFillColor(color.RGBA{255, 0, 0, 255}),
		SetSpaceColor(color.White),
	)

	v3 := NewScreenTone(
		SetBounds(b),
		SetRadius(5),
		SetSpacing(12),
		SetAngle(75),
		SetFillColor(color.RGBA{0, 0, 255, 255}),
		SetSpaceColor(color.White),
	)

	return stitchImagesForDemo(v1, v2, v3)
}

func init() {
	RegisterGenerator("ScreenTone", GenerateScreenTone)
}
