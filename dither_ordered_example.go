package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var OrderedDitherOutputFilename = "dither_ordered.png"
var OrderedDitherZoomLevels = []int{}

const OrderedDitherOrder = 102 // Arbitrary order

func ExampleNewOrderedDither() {
	i := NewDemoOrderedDither()
	f, err := os.Create(OrderedDitherOutputFilename)
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

func GenerateOrderedDither(b image.Rectangle) image.Image {
	return NewDemoOrderedDither(SetBounds(b))
}

func NewDemoOrderedDither(ops ...func(any)) image.Image {
	g := NewLinearGradient(SetStartColor(color.Black), SetEndColor(color.White))
	for _, op := range ops {
		op(g)
	}
	return NewBayer4x4Dither(g, nil)
}

func GenerateOrderedDitherReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	gopher := NewGopher()
	bw := color.Palette{color.Black, color.White}

	return map[string]func(image.Rectangle) image.Image{
		"Bayer2x2": func(b image.Rectangle) image.Image {
			return NewBayer2x2Dither(gopher, bw)
		},
		"Bayer4x4": func(b image.Rectangle) image.Image {
			return NewBayer4x4Dither(gopher, bw)
		},
		"Bayer8x8": func(b image.Rectangle) image.Image {
			return NewBayer8x8Dither(gopher, bw)
		},
		"Halftone4x4": func(b image.Rectangle) image.Image {
			return NewHalftoneDither(gopher, 4, bw)
		},
		"Halftone8x8": func(b image.Rectangle) image.Image {
			return NewHalftoneDither(gopher, 8, bw)
		},
		"Random": func(b image.Rectangle) image.Image {
			return NewRandomDither(gopher, bw, 12345)
		},
	}, []string{
		"Bayer2x2", "Bayer4x4", "Bayer8x8", "Halftone4x4", "Halftone8x8", "Random",
	}
}

func init() {
	RegisterGenerator("OrderedDither", GenerateOrderedDither)
	RegisterReferences("OrderedDither", GenerateOrderedDitherReferences)
}
