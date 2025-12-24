package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var BitwiseAndOutputFilename = "bitwise_and.png"
var BitwiseAndZoomLevels = []int{}
const BitwiseAndOrder = 37

func ExampleNewBitwiseAnd() {
	h := NewHorizontalLine(SetLineSize(50), SetSpaceSize(50), SetLineColor(color.RGBA{255, 0, 0, 255}))
	v := NewVerticalLine(SetLineSize(50), SetSpaceSize(50), SetLineColor(color.RGBA{0, 255, 0, 255}))
	p := NewBitwiseAnd([]image.Image{h, v})
	f, err := os.Create(BitwiseAndOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, p)
}

func GenerateBitwiseAnd(b image.Rectangle) image.Image {
	h := NewHorizontalLine(SetLineSize(50), SetSpaceSize(50), SetLineColor(color.RGBA{255, 0, 0, 255}), SetBounds(b))
	v := NewVerticalLine(SetLineSize(50), SetSpaceSize(50), SetLineColor(color.RGBA{0, 255, 0, 255}), SetBounds(b))
	return NewBitwiseAnd([]image.Image{h, v}, SetBounds(b))
}

func GenerateBitwiseAndReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("BitwiseAnd", GenerateBitwiseAnd)
	RegisterReferences("BitwiseAnd", GenerateBitwiseAndReferences)
}
