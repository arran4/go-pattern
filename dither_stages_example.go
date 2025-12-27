package pattern

import (
	"image"
	"image/color"
)

// DitherStagesOutputFilename defines the output filename for the stages example.
var DitherStagesOutputFilename = "dither_stages.png"
var DitherStagesZoomLevels = []int{}
const DitherStagesOrder = 103

// ExampleNewDitherStages demonstrates the progression of dithering techniques
// on a linear gradient, illustrating the "stages" or levels of detail each matrix provides.
func ExampleNewDitherStages() image.Image {
	// Default view: Bayer 8x8 on a gradient
	return NewBayer8x8Dither(NewLinearGradient(), nil)
}

// GenerateDitherStages is the generator function for the main example image.
func GenerateDitherStages(b image.Rectangle) image.Image {
	g := NewLinearGradient(SetStartColor(color.Black), SetEndColor(color.White), SetBounds(b))
	return NewBayer8x8Dither(g, nil)
}

// GenerateDitherStagesReferences generates a comprehensive set of dithering examples
// using a gradient source to visualize the patterns.
func GenerateDitherStagesReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	bw := color.Palette{color.Black, color.White}

	// Create a gradient source function
	grad := func(b image.Rectangle) image.Image {
		// Ensure the gradient fits the bounds
		return NewLinearGradient(SetStartColor(color.Black), SetEndColor(color.White), SetBounds(b))
	}

	// Custom 3x3 Matrix
	// 0 7 3
	// 6 5 2
	// 4 1 8
	mat3x3 := []float64{
		0.0/9.0, 7.0/9.0, 3.0/9.0,
		6.0/9.0, 5.0/9.0, 2.0/9.0,
		4.0/9.0, 1.0/9.0, 8.0/9.0,
	}

	return map[string]func(image.Rectangle) image.Image{
		"1_Threshold": func(b image.Rectangle) image.Image {
			// Thresholding is effectively Bayer 1x1 (scalar 0.5)
			// But lets simulate it with 1x1 matrix [0.5]
			return NewOrderedDither(grad(b), []float64{0.5}, 1, bw, 0)
		},
		"2_Random": func(b image.Rectangle) image.Image {
			return NewRandomDither(grad(b), bw, 1) // Seed 1
		},
		"3_Bayer2x2": func(b image.Rectangle) image.Image {
			// 5 levels
			return NewBayer2x2Dither(grad(b), bw)
		},
		"4_Bayer3x3": func(b image.Rectangle) image.Image {
			// 10 levels
			return NewOrderedDither(grad(b), mat3x3, 3, bw, 0)
		},
		"5_Bayer4x4": func(b image.Rectangle) image.Image {
			// 17 levels
			return NewBayer4x4Dither(grad(b), bw)
		},
		"6_Bayer8x8": func(b image.Rectangle) image.Image {
			// 65 levels
			return NewBayer8x8Dither(grad(b), bw)
		},
		"7_Halftone4x4": func(b image.Rectangle) image.Image {
			return NewHalftoneDither(grad(b), 4, bw)
		},
		"8_Halftone6x6": func(b image.Rectangle) image.Image {
			return NewHalftoneDither(grad(b), 6, bw)
		},
		"9_BlueNoise": func(b image.Rectangle) image.Image {
			// Void and Cluster / Blue Noise
			return NewBlueNoiseDither(grad(b), bw)
		},
		"10_Yliluoma1_8x8": func(b image.Rectangle) image.Image {
			return NewYliluoma1Dither(grad(b), bw, 8)
		},
		"11_Knoll_8x8": func(b image.Rectangle) image.Image {
			return NewKnollDither(grad(b), bw, 8)
		},
	}, []string{
		"1_Threshold",
		"2_Random",
		"3_Bayer2x2",
		"4_Bayer3x3",
		"5_Bayer4x4",
		"6_Bayer8x8",
		"7_Halftone4x4",
		"8_Halftone6x6",
		"9_BlueNoise",
		"10_Yliluoma1_8x8",
		"11_Knoll_8x8",
	}
}

func init() {
	RegisterGenerator("DitherStages", GenerateDitherStages)
	RegisterReferences("DitherStages", GenerateDitherStagesReferences)
}
