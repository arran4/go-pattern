package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
)

var BufferOutputFilename = "buffer.png"
var BufferZoomLevels = []int{}

const BufferOrder = 100

// Buffer Pattern
// A pattern that buffers a source image.
func ExampleNewBuffer() {
	// 1. Create a source pattern
	source := NewSolid(color.RGBA{255, 0, 0, 255})
	// 2. Create a buffer
	b := NewBuffer(source, SetExpiry(10*time.Second))
	// 3. Refresh the buffer explicitly to populate the cache
	b.Refresh()

	// Output:

	// Create the file for the example
	f, err := os.Create(BufferOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, b); err != nil {
		panic(err)
	}
}

// GenerateBuffer creates a buffer pattern for the gallery.
func GenerateBuffer(b image.Rectangle) image.Image {
	source := NewSolid(color.RGBA{100, 200, 100, 255})
	buf := NewBuffer(source, SetBounds(b))
	buf.Refresh()
	return buf
}

// NewSolid creates a uniform color pattern.
// Helper for the example.
func NewSolid(c color.Color) image.Image {
	return NewRect(SetFillColor(c))
}

func init() {
	RegisterGenerator("Buffer", GenerateBuffer)
}
