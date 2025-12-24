package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var BayerDitherOutputFilename = "bayer_dither.png"
var BayerDitherZoomLevels = []int{}
const BayerDitherOrder = 34

func ExampleNewBayerDither() {
	grad := NewLinearGradient(
		SetStartColor(color.Black),
		SetEndColor(color.White),
	)
	p := NewBayerDither(grad, 4)
	f, err := os.Create(BayerDitherOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, p)
}

func GenerateBayerDither(b image.Rectangle) image.Image {
	grad := NewLinearGradient(
		SetStartColor(color.Black),
		SetEndColor(color.White),
		SetBounds(b),
	)
	return NewBayerDither(grad, 4, SetBounds(b))
}

func GenerateBayerDitherReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("BayerDither", GenerateBayerDither)
	RegisterReferences("BayerDither", GenerateBayerDitherReferences)
}

var ErrorDiffusionOutputFilename = "dither_errordiffusion.png"
var ErrorDiffusionZoomLevels = []int{}

const ErrorDiffusionOrder = 101 // Arbitrary order

func ExampleNewErrorDiffusion() {
	// Standard example
	i := NewDemoErrorDiffusion()
	f, err := os.Create(ErrorDiffusionOutputFilename)
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

func GenerateErrorDiffusion(b image.Rectangle) image.Image {
	return NewDemoErrorDiffusion(SetBounds(b))
}

func NewDemoErrorDiffusion(ops ...func(any)) image.Image {
	// Gradient to show smooth dither
	g := NewLinearGradient(SetStartColor(color.Black), SetEndColor(color.White))
	// Apply options (including bounds) to gradient
	for _, op := range ops {
		op(g)
	}
	// Dither to BW
	return NewErrorDiffusion(g, FloydSteinberg, nil)
}

func GenerateErrorDiffusionReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	gopher := NewGopher()
	bw := color.Palette{color.Black, color.White}
	web := color.Palette{
		color.Black, color.White, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255},
		color.RGBA{255, 255, 0, 255}, color.RGBA{0, 255, 255, 255}, color.RGBA{255, 0, 255, 255},
	}

	return map[string]func(image.Rectangle) image.Image{
			"FloydSteinberg": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, FloydSteinberg, bw)
			},
			"JarvisJudiceNinke": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, JarvisJudiceNinke, bw)
			},
			"Stucki": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, Stucki, bw)
			},
			"Atkinson": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, Atkinson, bw)
			},
			"Burkes": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, Burkes, bw)
			},
			"SierraLite": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, SierraLite, bw)
			},
			"Sierra2": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, Sierra2, bw)
			},
			"Sierra3": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, Sierra3, bw)
			},
			"StevensonArce": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, StevensonArce, bw)
			},
			"FloydSteinberg_WebSafe": func(b image.Rectangle) image.Image {
				return NewErrorDiffusion(gopher, FloydSteinberg, web)
			},
		}, []string{
			"FloydSteinberg", "JarvisJudiceNinke", "Stucki", "Atkinson", "Burkes",
			"SierraLite", "Sierra2", "Sierra3", "StevensonArce", "FloydSteinberg_WebSafe",
		}
}

func init() {
	RegisterGenerator("ErrorDiffusion", GenerateErrorDiffusion)
	RegisterReferences("ErrorDiffusion", GenerateErrorDiffusionReferences)
}
