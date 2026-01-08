package pattern

import "image"

// PCBTracesOutputFilename is the default output filename for the PCBTraces pattern.
var PCBTracesOutputFilename = "pcbtraces.png"

// PCBTracesZoomLevels is the default zoom levels for the PCBTraces pattern.
var PCBTracesZoomLevels = []int{}

// GeneratePCBTraces creates a PCB trace layout in the provided bounds.
func GeneratePCBTraces(bounds image.Rectangle) image.Image {
	return NewPCBTraces(SetBounds(bounds))
}

// ExampleNewPCBTraces returns a sample PCB trace layout with default options.
func ExampleNewPCBTraces() image.Image {
	return GeneratePCBTraces(image.Rect(0, 0, 192, 192))
}

func init() {
	RegisterGenerator("PCBTraces", GeneratePCBTraces)
}
