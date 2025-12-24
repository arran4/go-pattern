package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var PerlinNoiseOutputFilename = "perlin_noise.png"
var PerlinNoiseZoomLevels = []int{}
const PerlinNoiseOrder = 20

// Perlin Noise Pattern
// Generates smooth procedural noise using the Perlin algorithm.
func ExampleNewPerlinNoise() {
	i := NewPerlinNoise(
		SetNoiseColorLow(color.Black),
		SetNoiseColorHigh(color.White),
		SetFrequency(0.1),
	)
	f, err := os.Create(PerlinNoiseOutputFilename)
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

func GeneratePerlinNoise(b image.Rectangle) image.Image {
	return NewPerlinNoise(SetBounds(b), SetFrequency(0.1))
}

func init() {
	RegisterGenerator("PerlinNoise", GeneratePerlinNoise)
}
