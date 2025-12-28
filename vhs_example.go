package pattern

import (
	"image"
	"image/png"
	"math"
	"os"
)

var VHSOutputFilename = "vhs.png"
var VHSZoomLevels = []int{}
const VHSOrder = 40 // Adjust order as needed

// Retro VHS Effect
// Demonstrates the VHS scanline, color shift, and noise effect.
func ExampleNewVHS() {
	// Use the embedded Gopher image as the source
	src := NewGopher()

	// Apply VHS effect
	i := NewVHS(src,
		SetScanlineFrequency(math.Pi),
		SetScanlineIntensity(0.3),
		SetColorOffset(4),
		SetNoiseIntensity(0.15),
		SetSeed(42),
	)

	f, err := os.Create(VHSOutputFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateVHS(b image.Rectangle) image.Image {
	// Use the embedded Gopher image as the source
	src := NewGopher()
	// Or maybe a more colorful image to show off the channel shift?
	// NewGopher is good.

	// Apply VHS effect
	// Default frequency is math.Pi (every 2 pixels).

	return NewVHS(src,
		SetScanlineFrequency(math.Pi),
		SetScanlineIntensity(0.3),
		SetColorOffset(4),
		SetNoiseIntensity(0.15),
		SetSeed(42),
	)
}

func init() {
	RegisterGenerator("VHS", GenerateVHS)
}
