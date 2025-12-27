package pattern

import (
	"image"
	"image/color"
)

var (
	GlobeOutputFilename = "globe.png"
	Globe_wireframeOutputFilename = "globe_wireframe.png"
)

// ExampleNewGlobe generates a globe pattern.
func ExampleNewGlobe() image.Image {
	return NewGlobe(
		SetLatitudeLines(5),
		SetLongitudeLines(12),
		SetLineSize(2),
		SetLineColor(color.RGBA{200, 200, 200, 255}),
		SetFillColor(color.RGBA{0, 0, 50, 255}),
		SetSpaceColor(color.Transparent),
		SetAngle(15),
	)
}

// ExampleNewGlobe_wireframe generates a wireframe globe.
func ExampleNewGlobe_wireframe() image.Image {
	return NewGlobe(
		SetLatitudeLines(10),
		SetLongitudeLines(20),
		SetLineSize(1),
		SetLineColor(color.Black),
		SetSpaceColor(color.White),
		SetAngle(-23.5), // Earth tilt
	)
}
