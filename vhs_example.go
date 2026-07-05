package pattern

import (
	"image"
	"math"
)

var VHSOutputFilename = "vhs.png"
var VHSZoomLevels = []int{}

const VHSOrder = 40 // Adjust order as needed

// Retro VHS Effect
// Demonstrates the VHS scanline, color shift, and noise effect.
func ExampleNewVHS() {
	// See GenerateVHS for implementation details
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
