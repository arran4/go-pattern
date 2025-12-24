package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var EdgeDetectOutputFilename = "edgedetect.png"
var EdgeDetectZoomLevels = []int{} // No zoom levels by default for now, or maybe 1?

const EdgeDetectOrder = 100 // Arbitrary order

// EdgeDetect Pattern
// Applies Sobel edge detection to an input image.
func ExampleNewEdgeDetect() {
	i := NewDemoEdgeDetect()
	f, err := os.Create(EdgeDetectOutputFilename)
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

func GenerateEdgeDetect(b image.Rectangle) image.Image {
	return NewDemoEdgeDetect(SetBounds(b))
}

func GenerateEdgeDetectReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	// We want to show the original image vs the edge detected one

	sourceGen := func(b image.Rectangle) image.Image {
		chk := NewChecker(color.Black, color.White, SetBounds(b))
		return NewSimpleZoom(chk, 20, SetBounds(b))
	}

	gopherGen := func(b image.Rectangle) image.Image {
		// NewGopher returns a fixed size image. We might want to scale it or just return it?
		// The pattern framework usually expects bounds.
		// If we use NewGopher, it ignores 'b' unless we wrap it.
		// Let's just return it, but maybe centered or tiled if 'b' is large?
		// For demo, the bootstrap tool calls with 150x150. Gopher is likely larger/smaller.
		// Let's use it as is.
		return NewGopher()
	}

	gopherEdgesGen := func(b image.Rectangle) image.Image {
		return NewEdgeDetect(NewGopher())
	}

	return map[string]func(image.Rectangle) image.Image{
		"Source":       sourceGen,
		"Gopher":       gopherGen,
		"Gopher Edges": gopherEdgesGen,
	}, []string{"Source", "Gopher", "Gopher Edges"}
}

func init() {
	RegisterGenerator("EdgeDetect", GenerateEdgeDetect)
	RegisterReferences("EdgeDetect", GenerateEdgeDetectReferences)
}
