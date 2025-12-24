package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure EdgeDetect implements the image.Image interface.
var _ image.Image = (*EdgeDetect)(nil)

// EdgeDetect is a pattern that applies edge detection (Sobel operator) to an underlying image.
type EdgeDetect struct {
	img         image.Image
	sensitivity float64 // Multiplier for magnitude. Default 1.0 (clamped).
	// Or maybe normalization factor?
	// Let's call it "Gain".
}

func (e *EdgeDetect) ColorModel() color.Model {
	return color.GrayModel
}

func (e *EdgeDetect) Bounds() image.Rectangle {
	return e.img.Bounds()
}

// lum calculates the luminance of a color in the 0.0-1.0 range.
func lum(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	// Standard luminance conversion: 0.299R + 0.587G + 0.114B
	// Since RGBA returns premultiplied alpha 0-65535, we divide by 65535.0.
	y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return y / 65535.0
}

func (e *EdgeDetect) At(x, y int) color.Color {
	// Sobel operator kernels
	// Gx: -1 0 1
	//     -2 0 2
	//     -1 0 1
	//
	// Gy: -1 -2 -1
	//      0  0  0
	//      1  2  1

	// We need pixel values from x-1 to x+1, y-1 to y+1

	// Grid of 3x3
	var grid [3][3]float64
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			grid[j+1][i+1] = lum(e.img.At(x+i, y+j))
		}
	}

	gx := -1.0*grid[0][0] + 1.0*grid[0][2] +
		-2.0*grid[1][0] + 2.0*grid[1][2] +
		-1.0*grid[2][0] + 1.0*grid[2][2]

	gy := -1.0*grid[0][0] + -2.0*grid[0][1] + -1.0*grid[0][2] +
		1.0*grid[2][0] + 2.0*grid[2][1] + 1.0*grid[2][2]

	mag := math.Sqrt(gx*gx + gy*gy)

	// Normalization.
	// Max magnitude is approx sqrt(32) ≈ 5.66.
	// To make a full contrast edge (0 to 1) result in 1.0, we should divide by ~4 or ~5.66.
	// If we leave it as is, it's very sensitive (mag 1.0 reached with small gradient).
	// Let's normalize by 4.0 by default to give a reasonable "0-1" output for "0-1" input edges.

	val := mag / 4.0

	if val > 1.0 {
		val = 1.0
	}

	return color.Gray{Y: uint8(val * 255)}
}

// NewEdgeDetect creates a new EdgeDetect pattern from an existing image.
func NewEdgeDetect(img image.Image, ops ...func(any)) image.Image {
	e := &EdgeDetect{
		img: img,
	}
	for _, op := range ops {
		op(e)
	}
	return e
}

// NewDemoEdgeDetect produces a demo variant for readme.md
func NewDemoEdgeDetect(ops ...func(any)) image.Image {
	// A checker pattern has sharp edges.
	chk := NewChecker(color.Black, color.White, ops...)
	// Zoom it so we have larger blocks, edges are clearer
	zm := NewSimpleZoom(chk, 20, ops...)
	return NewEdgeDetect(zm, ops...)
}
