package pattern

import (
	"image"
	"image/png"
	"os"
)

var NoiseOutputFilename = "noise.png"
var NoiseZoomLevels = []int{} // Zooming is unnecessary for noise

const NoiseOrder = 20

// Noise Pattern
// Generates random noise using various algorithms (Crypto, Hash).
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
		"Hash": func(b image.Rectangle) image.Image {
			return NewNoise(SetBounds(b), SetNoiseAlgorithm(&HashNoise{Seed: 12345}))
		},
		"Hash2": func(b image.Rectangle) image.Image {
			return NewNoise(SetBounds(b), SetNoiseAlgorithm(&HashNoise{Seed: 67890}))
		},
	}, []string{"Crypto", "Hash", "Hash2"}
}

func init() {
	RegisterGenerator("Noise", GenerateNoise)
	RegisterReferences("Noise", GenerateNoiseReferences)
}
