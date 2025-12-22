package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var AndOutputFilename = "boolean_and.png"
var AndZoomLevels = []int{2, 4}
const AndOrder = 20

func ExampleNewAnd() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
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

func BootstrapAnd(b image.Rectangle) image.Image {
	return NewDemoAnd(SetBounds(b))
}

var OrOutputFilename = "boolean_or.png"
var OrZoomLevels = []int{2, 4}
const OrOrder = 21

func ExampleNewOr() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
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

func BootstrapOr(b image.Rectangle) image.Image {
	return NewDemoOr(SetBounds(b))
}

var XorOutputFilename = "boolean_xor.png"
var XorZoomLevels = []int{2, 4}
const XorOrder = 22

func ExampleNewXor() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
	v := NewVerticalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
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

func BootstrapXor(b image.Rectangle) image.Image {
	return NewDemoXor(SetBounds(b))
}

var NotOutputFilename = "boolean_not.png"
var NotZoomLevels = []int{2, 4}
const NotOrder = 23

func ExampleNewNot() {
	h := NewHorizontalLine(SetLineSize(20), SetSpaceSize(20), SetLineColor(color.White))
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

func BootstrapNot(b image.Rectangle) image.Image {
	return NewDemoNot(SetBounds(b))
}

func init() {
	RegisterGenerator("And", BootstrapAnd)
	RegisterGenerator("Or", BootstrapOr)
	RegisterGenerator("Xor", BootstrapXor)
	RegisterGenerator("Not", BootstrapNot)
}
