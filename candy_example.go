package pattern

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var CandyOutputFilename = "candy.png"

const CandyBaseLabel = "Candy"

// Candy Example (M&Ms / Smarties)
// Demonstrates using the Scatter pattern to draw overlapping, colored candy circles.
func ExampleNewCandy() {
	// 1. Define colors for our candy.
	colors := []color.RGBA{
		{255, 0, 0, 255},   // Red
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{255, 255, 0, 255}, // Yellow
		{255, 165, 0, 255}, // Orange
		{139, 69, 19, 255}, // Brown
	}

	// 2. Create the Scatter pattern.
	candy := NewScatter(
		SetScatterFrequency(0.04), // Controls size/spacing relative to pixels
		SetScatterDensity(0.9),    // High density
		SetScatterGenerator(func(u, v float64, hash uint64) (color.Color, float64) {
			// Radius of the candy
			radius := 14.0

			// Distance from center
			distSq := u*u + v*v
			if distSq > radius*radius {
				return color.Transparent, 0
			}
			dist := math.Sqrt(distSq)

			// Pick a random color based on hash
			colIdx := hash % uint64(len(colors))
			baseCol := colors[colIdx]

			// Simple shading: slightly darker at edges, highlight at top-left
			// Spherical shading approx
			// Normal vector (nx, ny, nz)
			// z = sqrt(1 - x^2 - y^2)
			nx := u / radius
			ny := v / radius
			nz := math.Sqrt(math.Max(0, 1.0-nx*nx-ny*ny))

			// Light source direction (top-left)
			lx, ly, lz := -0.5, -0.5, 0.7
			lLen := math.Sqrt(lx*lx + ly*ly + lz*lz)
			lx, ly, lz = lx/lLen, ly/lLen, lz/lLen

			// Diffuse
			dot := nx*lx + ny*ly + nz*lz
			diffuse := math.Max(0, dot)

			// Specular (Glossy plastic look)
			// Reflected light vector
			// R = 2(N.L)N - L
			rx := 2*dot*nx - lx
			ry := 2*dot*ny - ly
			rz := 2*dot*nz - lz
			// View vector (straight up)
			vx, vy, vz := 0.0, 0.0, 1.0
			specDot := rx*vx + ry*vy + rz*vz
			specular := math.Pow(math.Max(0, specDot), 20) // Shininess

			// Apply lighting
			r := float64(baseCol.R)*(0.2+0.8*diffuse) + 255*specular*0.6
			g := float64(baseCol.G)*(0.2+0.8*diffuse) + 255*specular*0.6
			b := float64(baseCol.B)*(0.2+0.8*diffuse) + 255*specular*0.6

			// Clamp
			r = math.Min(255, math.Max(0, r))
			g = math.Min(255, math.Max(0, g))
			b = math.Min(255, math.Max(0, b))

			// Anti-aliasing at edge
			alpha := 1.0
			if dist > radius-1.0 {
				alpha = radius - dist
			}

			// Use hash for random Z-ordering
			z := float64(hash) / 18446744073709551615.0

			return color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(alpha * 255),
			}, z
		}),
	)

	f, err := os.Create(CandyOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, candy); err != nil {
		panic(err)
	}
}

func GenerateCandy(b image.Rectangle) image.Image {
	colors := []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
		{255, 255, 0, 255},
		{255, 165, 0, 255},
		{139, 69, 19, 255},
	}

	return NewScatter(
		SetBounds(b),
		SetScatterFrequency(0.04),
		SetScatterDensity(0.9),
		SetScatterGenerator(func(u, v float64, hash uint64) (color.Color, float64) {
			radius := 14.0
			distSq := u*u + v*v
			if distSq > radius*radius {
				return color.Transparent, 0
			}
			dist := math.Sqrt(distSq)
			colIdx := hash % uint64(len(colors))
			baseCol := colors[colIdx]

			nx := u / radius
			ny := v / radius
			nz := math.Sqrt(math.Max(0, 1.0-nx*nx-ny*ny))
			lx, ly, lz := -0.5, -0.5, 0.7
			lLen := math.Sqrt(lx*lx + ly*ly + lz*lz)
			lx, ly, lz = lx/lLen, ly/lLen, lz/lLen
			dot := nx*lx + ny*ly + nz*lz
			diffuse := math.Max(0, dot)
			rx := 2*dot*nx - lx
			ry := 2*dot*ny - ly
			rz := 2*dot*nz - lz
			specDot := rx*0 + ry*0 + rz*1
			specular := math.Pow(math.Max(0, specDot), 20)
			r := float64(baseCol.R)*(0.2+0.8*diffuse) + 255*specular*0.6
			g := float64(baseCol.G)*(0.2+0.8*diffuse) + 255*specular*0.6
			b := float64(baseCol.B)*(0.2+0.8*diffuse) + 255*specular*0.6
			alpha := 1.0
			if dist > radius-1.0 {
				alpha = radius - dist
			}

			// Use hash for random Z-ordering
			z := float64(hash) / 18446744073709551615.0

			return color.RGBA{uint8(math.Min(255, math.Max(0, r))), uint8(math.Min(255, math.Max(0, g))), uint8(math.Min(255, math.Max(0, b))), uint8(alpha * 255)}, z
		}),
	)
}

func init() {
	RegisterGenerator(CandyBaseLabel, GenerateCandy)
}
