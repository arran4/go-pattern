package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure NormalMap implements the image.Image interface.
var _ image.Image = (*NormalMap)(nil)

// NormalMap generates a tangent-space normal map from a source height map.
// The red channel represents the X (horizontal) vector.
// The green channel represents the Y (vertical) vector.
// The blue channel represents the Z (depth) vector.
type NormalMap struct {
	Source   image.Image
	Strength float64
}

// At returns the normal map color at the specified coordinates.
func (nm *NormalMap) At(x, y int) color.Color {
	// Sobel operator kernels
	// Gx: -1 0 1
	//     -2 0 2
	//     -1 0 1
	//
	// Gy: -1 -2 -1
	//      0  0  0
	//      1  2  1

	// Helper to get height (0.0-1.0)
	getHeight := func(x, y int) float64 {
		c := nm.Source.At(x, y)
		r, g, b, _ := c.RGBA()
		// Convert to grayscale/luminance
		// RGBA returns 0-65535.
		lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
		return lum / 65535.0
	}

	// Grid of 3x3
	var grid [3][3]float64
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			grid[j+1][i+1] = getHeight(x+i, y+j)
		}
	}

	gx := -1.0*grid[0][0] + 1.0*grid[0][2] +
		-2.0*grid[1][0] + 2.0*grid[1][2] +
		-1.0*grid[2][0] + 1.0*grid[2][2]

	gy := -1.0*grid[0][0] + -2.0*grid[0][1] + -1.0*grid[0][2] +
		1.0*grid[2][0] + 2.0*grid[2][1] + 1.0*grid[2][2]

	// The vector is (-gx, -gy, 1.0).
	// Strength scales the slope. Higher strength -> steeper slope.
	// We apply strength to gx and gy.

	dx := -gx * nm.Strength
	dy := -gy * nm.Strength
	dz := 1.0

	// Normalize
	len := math.Sqrt(dx*dx + dy*dy + dz*dz)
	nx := dx / len
	ny := dy / len
	nz := dz / len

	// Map [-1, 1] to [0, 255]
	r := uint8((nx+1.0)*0.5*255.0 + 0.5)
	g := uint8((ny+1.0)*0.5*255.0 + 0.5)
	b := uint8((nz+1.0)*0.5*255.0 + 0.5)

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// Bounds delegates to the source image.
func (nm *NormalMap) Bounds() image.Rectangle {
	if nm.Source == nil {
		return image.Rect(0, 0, 255, 255)
	}
	return nm.Source.Bounds()
}

// ColorModel returns RGB.
func (nm *NormalMap) ColorModel() color.Model {
	return color.RGBAModel
}

// NewNormalMap creates a new NormalMap from a source image.
// Default Strength is 1.0.
func NewNormalMap(source image.Image, ops ...func(any)) image.Image {
	nm := &NormalMap{
		Source:   source,
		Strength: 1.0,
	}
	for _, op := range ops {
		op(nm)
	}
	return nm
}

// NormalMapStrength sets the strength of the normal map.
func NormalMapStrength(strength float64) func(any) {
	return func(i any) {
		if nm, ok := i.(*NormalMap); ok {
			nm.Strength = strength
		}
	}
}
