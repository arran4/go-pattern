package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestMaskedBlend(t *testing.T) {
	// Background: Red
	bg := image.NewUniform(color.RGBA{255, 0, 0, 255})
	// Foreground: Blue
	fg := image.NewUniform(color.RGBA{0, 0, 255, 255})
	// Mask: Linear Gradient (Black to White) horizontal
	// 0..100 -> 0..255
	mask := NewGeneric(func(x, _ int) color.Color {
		val := x * 255 / 100
		if val < 0 { val = 0 }
		if val > 255 { val = 255 }
		v := uint8(val)
		return color.RGBA{v, v, v, 255}
	})

	mb := NewMaskedBlend(bg, fg, mask)

	// Check x=0 (Mask=0 -> Red)
	c0 := mb.At(0, 0)
	r0, _, b0, _ := c0.RGBA()
	// Red component should be high (near 0xFFFF), Blue low (0)
	// c0.RGBA returns premultiplied alpha, but alpha is 255 (0xFFFF).
	if r0 < 0xF000 {
		t.Errorf("At(0,0) expected Red (R>=0xF000), got R=%x", r0)
	}
	if b0 > 0x1000 {
		t.Errorf("At(0,0) expected Red (B<=0x1000), got B=%x", b0)
	}

	// Check x=100 (Mask=1 -> Blue)
	c100 := mb.At(100, 0)
	r100, _, b100, _ := c100.RGBA()
	if r100 > 0x1000 {
		t.Errorf("At(100,0) expected Blue (R<=0x1000), got R=%x", r100)
	}
	if b100 < 0xF000 {
		t.Errorf("At(100,0) expected Blue (B>=0xF000), got B=%x", b100)
	}

	// Check x=50 (Mask=0.5 -> Purple)
	c50 := mb.At(50, 0)
	r50, _, b50, _ := c50.RGBA()
	// Should be around 127/255 * 0xFFFF = 0x7F7F
	// Let's allow a range 0x7000 - 0x9000
	if r50 < 0x7000 || r50 > 0x9000 {
		t.Errorf("At(50,0) expected ~Purple (R~0x8000), got R=%x", r50)
	}
	if b50 < 0x7000 || b50 > 0x9000 {
		t.Errorf("At(50,0) expected ~Purple (B~0x8000), got B=%x", b50)
	}
}
