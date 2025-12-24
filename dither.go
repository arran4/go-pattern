package pattern

import (
	"image"
	"image/color"
)

// Ensure BayerDither implements the image.Image interface.
var _ image.Image = (*BayerDither)(nil)

// BayerDither applies ordered dithering using a Bayer matrix.
type BayerDither struct {
	Null
	Input   image.Image
	Matrix  []uint8
	Size    int // 2, 4, 8
	Palette []color.Color
}

// Standard Bayer Matrices
var (
	Bayer2x2 = []uint8{
		0, 2,
		3, 1,
	}
	// 4x4 derived from 2x2
	Bayer4x4 = []uint8{
		0, 8, 2, 10,
		12, 4, 14, 6,
		3, 11, 1, 9,
		15, 7, 13, 5,
	}
)

func (p *BayerDither) At(x, y int) color.Color {
	if p.Input == nil {
		return color.RGBA{}
	}
	c := p.Input.At(x, y)

	// Get matrix value
	mx := x % p.Size
	if mx < 0 { mx += p.Size }
	my := y % p.Size
	if my < 0 { my += p.Size }

	threshold := p.Matrix[my*p.Size + mx]

	// Normalize threshold to 0-255?
	// Matrix values are 0..(size*size)-1.
	// We want to compare with pixel intensity.
	// Normalized threshold = (value + 0.5) / (size*size)

	n := p.Size * p.Size
	normT := float64(threshold) / float64(n)

	// Dither each channel? User said "apply to grayscale or per-channel".
	// Let's do per-channel.

	r, g, b, a := c.RGBA()

	// Convert 16-bit color to float 0-1
	fr := float64(r) / 65535.0
	fg := float64(g) / 65535.0
	fb := float64(b) / 65535.0

	// Apply dither
	// If val < threshold, it becomes darker?
	// Usually: val + (threshold - 0.5) ??
	// Or simply: if val > threshold ? 1 : 0 (for binary).
	// For multi-level palette?

	// If no palette is provided, we assume binary (black/white) or just thresholding?
	// The user mentions "Ordered dithering matrix".
	// Standard ordered dither on truecolor images is usually:
	// c_out = closest_palette_color(c_in + spread * (threshold - 0.5))

	// If we just want 1-bit per channel:
	dr := 0.0; if fr > normT { dr = 1.0 }
	dg := 0.0; if fg > normT { dg = 1.0 }
	db := 0.0; if fb > normT { db = 1.0 }

	return color.RGBA{
		R: uint8(dr * 255),
		G: uint8(dg * 255),
		B: uint8(db * 255),
		A: uint8(a >> 8),
	}
}

// NewBayerDither creates a new BayerDither pattern.
func NewBayerDither(input image.Image, size int, ops ...func(any)) image.Image {
	var mat []uint8
	if size == 4 {
		mat = Bayer4x4
	} else {
		size = 2
		mat = Bayer2x2
	}

	p := &BayerDither{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Input: input,
		Matrix: mat,
		Size: size,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}
