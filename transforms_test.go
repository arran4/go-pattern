package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestStandardLibraryTransforms(t *testing.T) {
	// create a dummy source
	src := NewNull()

	// 1. Rotate
	_ = NewRotate(src, 90)

	// 2. Scale
	_ = NewScale(src, ScaleX(0.5))

	// 3. Tile
	_ = NewTile(src, image.Rect(0, 0, 100, 100))

	// 4. Mirror
	_ = NewMirror(src, true, false)

	// 5. Warp
	_ = NewWarp(src)

	// 6. Clamp
	_ = NewClamp(src, image.Rect(0, 0, 100, 100))

	// 7. Remap
	mapImg := NewNull()
	_ = NewRemap(src, mapImg)

	// All transforms are present and accessible.
}

// copy of helper from remap_example.go for testing
type genericTest struct {
	f func(x, y int) color.Color
}

func (g *genericTest) ColorModel() color.Model { return color.RGBAModel }
func (g *genericTest) Bounds() image.Rectangle { return image.Rect(0, 0, 1000, 1000) }
func (g *genericTest) At(x, y int) color.Color { return g.f(x, y) }
func newGenericTest(f func(x, y int) color.Color) image.Image { return &genericTest{f: f} }

func TestClamp(t *testing.T) {
	// 10x10 white square
	src := newGenericTest(func(x, y int) color.Color {
		if x >= 0 && x < 10 && y >= 0 && y < 10 {
			return color.White
		}
		return color.Transparent
	})
	// Set actual bounds
	src = NewCrop(src, image.Rect(0, 0, 10, 10))

	// Clamp to 20x20
	clamped := NewClamp(src, image.Rect(0, 0, 20, 20))

	// Check inside
	if clamped.At(5, 5) != color.White {
		t.Error("Inside pixel should be white")
	}

	// Check extended edge (15, 5) -> should sample at x=9
	if clamped.At(15, 5) != color.White {
		t.Error("Extended pixel should be white")
	}

	// Check outside if it was transparent originally?
	// Our source function returns transparent outside 0..10.
	// But Clamp clamps coordinates to 0..9.
	// So it should always be white if the source is fully white in 0..10.
}

func TestRemap(t *testing.T) {
	// Source: Left half Black, Right half White
	src := newGenericTest(func(x, y int) color.Color {
		if x < 50 {
			return color.Black
		}
		return color.White
	})
	// 100x100
	src = NewCrop(src, image.Rect(0, 0, 100, 100))

	// Map: Flip Horizontal
	// u = 1 - x/w
	uv := newGenericTest(func(x, y int) color.Color {
		u := 1.0 - float64(x)/100.0
		v := float64(y)/100.0
		if u < 0 { u = 0 }

		return color.RGBA64{
			R: uint16(u * 65535),
			G: uint16(v * 65535),
			B: 0,
			A: 65535,
		}
	})
	uv = NewCrop(uv, image.Rect(0, 0, 100, 100))

	remapped := NewRemap(src, uv)

	// At x=10 (Left), u ~ 0.9 (Right). Source Right is White.
	// So result should be White.
	c := remapped.At(10, 50)
	r, g, b, _ := c.RGBA()
	if r < 30000 || g < 30000 || b < 30000 {
		t.Error("Expected White at left side (flipped), got", c)
	}

	// At x=90 (Right), u ~ 0.1 (Left). Source Left is Black.
	// So result should be Black.
	c = remapped.At(90, 50)
	r, g, b, _ = c.RGBA()
	if r > 30000 {
		t.Error("Expected Black at right side (flipped), got", c)
	}
}
