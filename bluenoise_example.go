package pattern

import (
	"image"
	"image/png"
	"os"
)

var BlueNoiseOutputFilename = "bluenoise.png"
var BlueNoiseZoomLevels = []int{}

const BlueNoiseOrder = 35

func ExampleNewBlueNoise() {
	p := NewBlueNoise()
	f, err := os.Create(BlueNoiseOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, p); err != nil {
		panic(err)
	}
}

func GenerateBlueNoise(b image.Rectangle) image.Image {
	return NewBlueNoise(SetBounds(b))
}

func GenerateBlueNoiseReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return nil, nil
}

func init() {
	RegisterGenerator("BlueNoise", GenerateBlueNoise)
	RegisterReferences("BlueNoise", GenerateBlueNoiseReferences)
}
