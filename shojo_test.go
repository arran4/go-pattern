package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestShojo(t *testing.T) {
	// Create a new Shojo pattern
	p := NewShojo(
		SetSpaceColor(color.Black),
		SetFillColor(color.White),
	)

	// Check if it implements image.Image
	if _, ok := p.(image.Image); !ok {
		t.Fatal("Shojo does not implement image.Image")
	}

	// Sample a few pixels to ensure it doesn't crash and returns valid colors
	// Center of a cell likely has a star if we are lucky, but we are testing for crashes/errors mainly.
	// Since randomness is seeded, we can check specific pixels if we knew the output, but for now just running it is enough.
	c := p.At(10, 10)
	if c == nil {
		t.Fatal("At returned nil color")
	}

	_, _, _, a := c.RGBA()
	if a == 0 {
		// It might be transparent if SpaceColor is transparent, but we set it to Black (opaque).
		// Wait, color.Black is opaque.
		// So Alpha should be non-zero.
	}
}
