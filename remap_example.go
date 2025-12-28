package pattern

import (
	"image"
	"image/color"
)

var (
	RemapOutputFilename = "remap.png"
	RemapZoomLevels     = []int{}
	RemapOrder          = 21
	RemapBaseLabel      = "Remap"
)

func init() {
	RegisterGenerator("Remap", func(bounds image.Rectangle) image.Image {
		return ExampleNewRemap(SetBounds(bounds))
	})
	RegisterReferences("Remap", GenerateRemapReferences)
}

func ExampleNewRemap(ops ...func(any)) image.Image {
	src := NewGopher()
	// Create a UV map that flips the image horizontally
	uv := newGeneric(func(x, y int) color.Color {
		w, h := 150.0, 150.0
		u := 1.0 - (float64(x) / w) // Flip X
		v := float64(y) / h
		return color.RGBA64{
			R: uint16(u * 65535),
			G: uint16(v * 65535),
			B: 0,
			A: 65535,
		}
	})
	uv = NewCrop(uv, image.Rect(0, 0, 150, 150))
	// ops from bootstrap might override bounds.
	// But NewRemap signature changed to NewRemap(src, ops...). Map must be passed via option.

	newOps := append([]func(any){RemapMap(uv)}, ops...)
	return NewRemap(src, newOps...)
}

func GenerateRemapReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Swirl": func(bounds image.Rectangle) image.Image {
			src := NewChecker(color.Black, color.White)
			src = NewCrop(src, image.Rect(0, 0, 150, 150))

			uv := newGeneric(func(x, y int) color.Color {
				w, h := 150.0, 150.0
				u := float64(x) / w
				v := float64(y) / h

				cx, cy := 0.5, 0.5
				dx, dy := u-cx, v-cy
				dist := dx*dx + dy*dy

				factor := 1.0
				u += (v - 0.5) * factor * (1.0 - dist*2.0)

				if u < 0 { u = 0 }
				if u > 1 { u = 1 }
				if v < 0 { v = 0 }
				if v > 1 { v = 1 }

				return color.RGBA64{
					R: uint16(u * 65535),
					G: uint16(v * 65535),
					B: 0,
					A: 65535,
				}
			})
			uv = NewCrop(uv, bounds)
			return NewRemap(src, RemapMap(uv))
		},
	}, []string{"Swirl"}
}

// Helper to create a pattern from a function
type generic struct {
	f func(x, y int) color.Color
}

func (g *generic) ColorModel() color.Model { return color.RGBAModel }
func (g *generic) Bounds() image.Rectangle { return image.Rect(0, 0, 1000, 1000) }
func (g *generic) At(x, y int) color.Color { return g.f(x, y) }
func newGeneric(f func(x, y int) color.Color) image.Image { return &generic{f: f} }
