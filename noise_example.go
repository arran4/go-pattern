package pattern

import (
	"image"
	"image/png"
	"os"
)

var NoiseOutputFilename = "noise.png"
var NoiseZoomLevels = []int{2}

const NoiseOrder = 20

// Noise Pattern
// Generates random noise using various algorithms (Crypto, Uniform, Pi).
func ExampleNewNoise() {
	// Create a noise pattern with default (Crypto) algorithm
	i := NewNoise()
	f, err := os.Create(NoiseOutputFilename)
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

func GenerateNoise(b image.Rectangle) image.Image {
	return NewNoise(SetBounds(b))
}

func GenerateNoiseReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Crypto": func(b image.Rectangle) image.Image {
			return NewNoise(SetBounds(b), SetNoiseAlgorithm(&CryptoNoise{}))
		},
		"Uniform": func(b image.Rectangle) image.Image {
			return NewNoise(SetBounds(b), SetNoiseAlgorithm(&UniformNoise{Seed: 12345}))
		},
		"Pi": func(b image.Rectangle) image.Image {
			return NewNoise(SetBounds(b), SetNoiseAlgorithm(&PiNoise{Stride: 30}))
		},
	}, []string{"Crypto", "Uniform", "Pi"}
}

func init() {
	RegisterGenerator("Noise", GenerateNoise)
	RegisterReferences("Noise", GenerateNoiseReferences)
}
