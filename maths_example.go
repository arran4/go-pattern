package pattern

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"
)

var MathsMandelbrotOutputFilename = "maths_mandelbrot.png"
var MathsMandelbrotZoomLevels = []int{}
const MathsMandelbrotOrder = 20

// Mandelbrot Set
// Generates a Mandelbrot set visualization.
func ExampleNewMathsMandelbrot() {
	// See GenerateMathsMandelbrot for implementation details
}

func GenerateMathsMandelbrot(b image.Rectangle) image.Image {
	w, h := float64(b.Dx()), float64(b.Dy())
	f := func(x, y int) color.Color {
		// Map pixel coordinates to complex plane
		// Center: -0.5, 0
		// Scale: 2.0 / min(w, h)
		scale := 3.0 / math.Min(w, h)
		cx := float64(x) - w/2
		cy := float64(y) - h/2
		z := complex(cx*scale-0.7, cy*scale)

		// Check if point is in Mandelbrot set
		// z_n+1 = z_n^2 + c
		// c is the point in the complex plane
		c := z
		z = 0
		iters := 200
		for i := 0; i < iters; i++ {
			if cmplx.Abs(z) > 2 {
				// Escaped
				gray := uint8(255 - (255 * i / iters))
				return color.RGBA{gray, gray, gray, 255}
			}
			z = z*z + c
		}
		return color.Black
	}
	return NewMaths(f, SetBounds(b))
}

var MathsJuliaOutputFilename = "maths_julia.png"
var MathsJuliaZoomLevels = []int{}
const MathsJuliaOrder = 21

// Julia Set
// Generates a Julia set visualization.
func ExampleNewMathsJulia() {
	// See GenerateMathsJulia for implementation details
}

func GenerateMathsJulia(b image.Rectangle) image.Image {
	w, h := float64(b.Dx()), float64(b.Dy())
	c := complex(-0.7, 0.27015) // Constant for Julia set

	f := func(x, y int) color.Color {
		scale := 3.0 / math.Min(w, h)
		cx := float64(x) - w/2
		cy := float64(y) - h/2
		z := complex(cx*scale, cy*scale)

		iters := 200
		for i := 0; i < iters; i++ {
			if cmplx.Abs(z) > 2 {
				gray := uint8(255 - (255 * i / iters))
				return color.RGBA{gray, gray, gray, 255}
			}
			z = z*z + c
		}
		return color.Black
	}
	return NewMaths(f, SetBounds(b))
}

var MathsSineOutputFilename = "maths_sine.png"
var MathsSineZoomLevels = []int{}
const MathsSineOrder = 22

// Sine Waves
// Generates a sine wave pattern.
func ExampleNewMathsSine() {
	// See GenerateMathsSine for implementation details
}

func GenerateMathsSine(b image.Rectangle) image.Image {
	f := func(x, y int) color.Color {
		// y = A * sin(B * x + C) + D
		// Draw a sine wave

		val := math.Sin(float64(x) * 0.1)

		// Map -1..1 to screen coordinates
		// Center line at h/2
		// Amplitude h/4

		targetY := float64(b.Dy())/2 + val * float64(b.Dy())/4

		dist := math.Abs(float64(y) - targetY)

		if dist < 2.0 {
			return color.Black
		}
		return color.White
	}
	return NewMaths(f, SetBounds(b))
}

var MathsWavesOutputFilename = "maths_waves.png"
var MathsWavesZoomLevels = []int{}
const MathsWavesOrder = 23

// Interference Waves
// Generates an interference pattern from multiple sine waves.
func ExampleNewMathsWaves() {
	// See GenerateMathsWaves for implementation details
}

func GenerateMathsWaves(b image.Rectangle) image.Image {
	f := func(x, y int) color.Color {
		cx, cy := float64(b.Dx())/2, float64(b.Dy())/2

		// Distance from center
		dx := float64(x) - cx
		dy := float64(y) - cy
		dist := math.Sqrt(dx*dx + dy*dy)

		// Sin(dist)
		v := math.Sin(dist * 0.1)

		// Map -1..1 to 0..255
		gray := uint8((v + 1.0) / 2.0 * 255.0)
		return color.RGBA{gray, gray, gray, 255}
	}
	return NewMaths(f, SetBounds(b))
}


func init() {
	RegisterGenerator("MathsMandelbrot", GenerateMathsMandelbrot)
	RegisterGenerator("MathsJulia", GenerateMathsJulia)
	RegisterGenerator("MathsSine", GenerateMathsSine)
	RegisterGenerator("MathsWaves", GenerateMathsWaves)
}
