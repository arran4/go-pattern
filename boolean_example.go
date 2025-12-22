package pattern

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

var AndOutputFilename = "boolean_and.png"
var AndZoomLevels = []int{2, 4}
const AndOrder = 20

func ExampleNewAnd() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	i := NewAnd([]image.Image{h, v})

	f, err := os.Create(AndOutputFilename)
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

func GenerateAnd(b image.Rectangle) image.Image {
	hWhite := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black), SetBounds(b))
	vWhite := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black), SetBounds(b))

	pred := PredicateAverageGrayAbove(128)

	i := NewAnd([]image.Image{hWhite, vWhite}, SetPredicate(pred), SetBounds(b))

	return stitchImages(hWhite, vWhite, i)
}

var OrOutputFilename = "boolean_or.png"
var OrZoomLevels = []int{2, 4}
const OrOrder = 21

func ExampleNewOr() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	i := NewOr([]image.Image{h, v})

	f, err := os.Create(OrOutputFilename)
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

func GenerateOr(b image.Rectangle) image.Image {
	hWhite := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black), SetBounds(b))
	vWhite := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black), SetBounds(b))

	pred := PredicateAverageGrayAbove(128)

	i := NewOr([]image.Image{hWhite, vWhite}, SetPredicate(pred), SetBounds(b))
	return stitchImages(hWhite, vWhite, i)
}

var XorOutputFilename = "boolean_xor.png"
var XorZoomLevels = []int{2, 4}
const XorOrder = 22

func ExampleNewXor() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	i := NewXor([]image.Image{h, v})

	f, err := os.Create(XorOutputFilename)
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

func GenerateXor(b image.Rectangle) image.Image {
	hWhite := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black), SetBounds(b))
	vWhite := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black), SetBounds(b))

	pred := PredicateAverageGrayAbove(128)

	i := NewXor([]image.Image{hWhite, vWhite}, SetPredicate(pred), SetBounds(b))
	return stitchImages(hWhite, vWhite, i)
}

var NotOutputFilename = "boolean_not.png"
var NotZoomLevels = []int{2, 4}
const NotOrder = 23

func ExampleNewNot() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.Black))
	i := NewNot(h)

	f, err := os.Create(NotOutputFilename)
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

func GenerateNot(b image.Rectangle) image.Image {
	hWhite := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White), SetSpaceColor(color.Black), SetBounds(b))

	pred := PredicateAverageGrayAbove(128)

	i := NewNot(hWhite, SetPredicate(pred), SetBounds(b))
	return stitchImages(hWhite, i)
}

// Helper to stitch images horizontally
func stitchImages(images ...image.Image) image.Image {
	if len(images) == 0 {
		return nil
	}

	width := 0
	height := 0
	for _, img := range images {
		b := img.Bounds()
		width += b.Dx()
		if b.Dy() > height {
			height = b.Dy()
		}
	}
	// Add padding
	padding := 10
	width += padding * (len(images) - 1)

	out := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with gray background to show extent
	draw.Draw(out, out.Bounds(), &image.Uniform{color.RGBA{200, 200, 200, 255}}, image.Point{}, draw.Src)

	x := 0
	for _, img := range images {
		b := img.Bounds()
		r := image.Rect(x, 0, x+b.Dx(), b.Dy())
		draw.Draw(out, r, img, b.Min, draw.Over)
		x += b.Dx() + padding
	}

	return out
}

func init() {
	RegisterGenerator("And", GenerateAnd)
	RegisterGenerator("Or", GenerateOr)
	RegisterGenerator("Xor", GenerateXor)
	RegisterGenerator("Not", GenerateNot)
}
