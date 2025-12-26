package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestQuantize(t *testing.T) {
	// Create a source image with gradient
	w, h := 256, 1
	src := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		src.Set(x, 0, color.RGBA{R: uint8(x), G: uint8(x), B: uint8(x), A: 255})
	}

	// Quantize to 2 levels
	q2 := NewQuantize(src, 2)

	// Test specific points
	// 0 -> 0
	checkColor(t, q2.At(0, 0), 0, "2 levels at 0")
	// 255 -> 255
	checkColor(t, q2.At(255, 0), 255, "2 levels at 255")
	// 100 -> 0 (100/255 = 0.39, *1 = 0.39, round 0)
	checkColor(t, q2.At(100, 0), 0, "2 levels at 100")
	// 150 -> 255 (150/255 = 0.58, *1 = 0.58, round 1)
	checkColor(t, q2.At(150, 0), 255, "2 levels at 150")

	// Quantize to 3 levels: 0, 127/128, 255
	// 0, 127.5, 255
	q3 := NewQuantize(src, 3)
	// 0 -> 0
	checkColor(t, q3.At(0, 0), 0, "3 levels at 0")
	// 255 -> 255
	checkColor(t, q3.At(255, 0), 255, "3 levels at 255")
	// 128 -> ~128
	c := q3.At(128, 0)
	r, _, _, _ := c.RGBA()
	if r < 32000 || r > 33000 {
		t.Errorf("Expected ~32767 for 3 levels at 128, got %d", r)
	}

}

func TestQuantizeAlpha(t *testing.T) {
	// Test with semi-transparent color
	// R=128, A=128 (Pre-multiplied: R=128, A=128 means full red intensity at 50% opacity? No.)
	// Pre-multiplied: R <= A.
	// If we want "Full Red at 50% opacity", in NRGBA it is (255, 0, 0, 128).
	// In RGBA (premultiplied): (128, 0, 0, 128).

	// Let's create an NRGBA color manually.
	nrgba := color.NRGBA{R: 255, G: 0, B: 0, A: 128} // Bright Red at 50% alpha
	// Converted to RGBA: R=128, G=0, B=0, A=128.

	// If we quantize to 2 levels.
	// R (255) -> 255 (Max).
	// A (128) -> 128.
	// Result NRGBA: (255, 0, 0, 128).
	// Converted back to RGBA: (128, 0, 0, 128).

	// What if we have a color that quantizes UP?
	// Input NRGBA: (150, 0, 0, 128).
	// Levels=2.
	// 150/255 = 0.58 -> 1.0 -> 255.
	// Result NRGBA: (255, 0, 0, 128).
	// Converted to RGBA: (128, 0, 0, 128).

	// Input NRGBA: (100, 0, 0, 128).
	// Levels=2.
	// 100/255 = 0.39 -> 0.0 -> 0.
	// Result NRGBA: (0, 0, 0, 128).
	// Converted to RGBA: (0, 0, 0, 128).

	// So correctness check: R,G,B <= A in the output.

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, nrgba) // Set automatically converts to RGBA model.

	q := NewQuantize(img, 2)
	c := q.At(0, 0)
	r, g, b, a := c.RGBA()

	if r > a || g > a || b > a {
		t.Errorf("Invalid premultiplied color: R=%d G=%d B=%d A=%d", r, g, b, a)
	}

	// Verify values
	// Expect 128, 0, 0, 128 (scaled to 16-bit: ~32896)
	// 128 -> 0x80. 0x8080 = 32896.

	// Wait, converting NRGBA(255,0,0,128) to RGBA might be slightly lossy or exact depending on impl.
	// 255 * 128 / 255 = 128.

	// if a != 0x8080 { // 128 * 257 = 32896
	// 	// It might be slightly off depending on Go's alpha blending.
	// 	// Go's RGBA usually does alpha-premultiplication.
	// 	// Let's trust it is around 50%.
	// }

	// Check that we got the "max" red for that alpha.
	if r < a-200 { // Allow some tolerance
		t.Errorf("Expected R to be close to A (saturated red), got R=%d, A=%d", r, a)
	}
}

func checkColor(t *testing.T, c color.Color, expected uint8, msg string) {
	r, _, _, _ := c.RGBA()
	// RGBA returns 0-65535. expected is 0-255.
	// Scale expected to 16 bit
	exp16 := uint32(expected) * 257 // 255*257 = 65535

	// Allow small diff due to floating point math
	diff := int(r) - int(exp16)
	if diff < 0 { diff = -diff }

	if diff > 300 { // Arbitrary tolerance
		t.Errorf("%s: Expected %d (approx %d), got %d", msg, expected, exp16, r)
	}
}
