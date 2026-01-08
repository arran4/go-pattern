package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	GlyphRingOutputFilename = "glyph_ring.png"
	GlyphRingZoomLevels     = []int{}
	GlyphRingOrder          = 101
)

const GlyphRingBaseLabel = "GlyphRing"

func init() {
	RegisterGenerator(GlyphRingBaseLabel, GenerateGlyphRing)
	RegisterReferences(GlyphRingBaseLabel, GenerateGlyphRingReferences)
}

// GenerateGlyphRing renders the circular glyph band pattern.
func GenerateGlyphRing(rect image.Rectangle) image.Image {
	return NewGlyphRing(
		SetBounds(rect),
		SetRadius(rect.Dx()/2-24),
		SetDensity(0.045),
		SetGlowColor(color.NRGBA{R: 140, G: 200, B: 255, A: 255}),
	)
}

// GenerateGlyphRingReferences returns no additional reference patterns.
func GenerateGlyphRingReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{}, []string{}
}

// ExampleNewGlyphRing produces a demo PNG for documentation.
func ExampleNewGlyphRing() {
	i := NewGlyphRing()
	f, err := os.Create(GlyphRingOutputFilename)
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
