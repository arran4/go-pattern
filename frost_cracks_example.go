package pattern

import (
	"image"
	"image/png"
	"os"
)

var FrostCracksOutputFilename = "frost_cracks.png"
var FrostCracksZoomLevels = []int{}

const FrostCracksBaseLabel = "FrostCracks"
const FrostCracksOrder = 120

// ExampleNewFrostCracks saves a PNG demonstrating icy fracture lines carved from fBm.
func ExampleNewFrostCracks() {
	img := GenerateFrostCracks(image.Rect(0, 0, 300, 300))

	f, err := os.Create(FrostCracksOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func GenerateFrostCracks(b image.Rectangle) image.Image {
	return NewFrostCracks(
		SetBounds(b),
		SetDensity(0.64),
		SetGlowAmount(0.55),
		SetBlurAmount(2),
	)
}

func init() {
	RegisterGenerator(FrostCracksBaseLabel, GenerateFrostCracks)
	RegisterReferences(FrostCracksBaseLabel, func() (map[string]func(image.Rectangle) image.Image, []string) {
		return nil, nil
	})
}
