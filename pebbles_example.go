package pattern

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var PebblesOutputFilename = "pebbles.png"

const PebblesBaseLabel = "Pebbles"

// Pebbles Example (Chipped Stone / Gravel)
// Demonstrates using the Scatter pattern to create overlapping, irregular stones.
func ExampleNewPebbles() {
	// Re-implement Pebbles using Scatter for true overlapping geometry.
	pebbles := NewScatter(
		SetScatterFrequency(0.04), // Size control
		SetScatterDensity(1.0),    // Packed tight
		SetScatterMaxOverlap(1),
		SetScatterGenerator(func(u, v float64, hash uint64) (color.Color, float64) {
			// Randomize size slightly
			rSize := float64(hash&0xFF)/255.0
			radius := 12.0 + rSize*6.0 // 12 to 18 pixels radius

			// Perturb the shape using simple noise (simulated by sin/cos of hash+angle)
			// to make it "chipped" or irregular.
			angle := math.Atan2(v, u)
			dist := math.Sqrt(u*u + v*v)

			// Simple radial noise
			noise := math.Sin(angle*5 + float64(hash%10)) * 0.1
			noise += math.Cos(angle*13 + float64(hash%7)) * 0.05

			effectiveRadius := radius * (1.0 + noise)

			if dist > effectiveRadius {
				return color.Transparent, 0
			}

			// Stone Color: Grey/Brown variations
			grey := 100 + int(hash%100)
			col := color.RGBA{uint8(grey), uint8(grey - 5), uint8(grey - 10), 255}

			// Shading (diffuse)
			// Normal estimation for a flattened spheroid
			nx := u / effectiveRadius
			ny := v / effectiveRadius
			nz := math.Sqrt(math.Max(0, 1.0 - nx*nx - ny*ny))

			// Light dir
			lx, ly, lz := -0.5, -0.5, 0.7
			lLen := math.Sqrt(lx*lx + ly*ly + lz*lz)
			lx, ly, lz = lx/lLen, ly/lLen, lz/lLen

			diffuse := math.Max(0, nx*lx + ny*ly + nz*lz)

			// Apply shading
			r := float64(col.R) * (0.1 + 0.9*diffuse)
			g := float64(col.G) * (0.1 + 0.9*diffuse)
			b := float64(col.B) * (0.1 + 0.9*diffuse)

			// Soft edge anti-aliasing
			alpha := 1.0
			edgeDist := effectiveRadius - dist
			if edgeDist < 1.0 {
				alpha = edgeDist
			}

			// Use hash for random Z-ordering
			z := float64(hash) / 18446744073709551615.0

			return color.RGBA{
				R: uint8(math.Min(255, r)),
				G: uint8(math.Min(255, g)),
				B: uint8(math.Min(255, b)),
				A: uint8(alpha * 255),
			}, z
		}),
	)

	f, err := os.Create(PebblesOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, pebbles); err != nil {
		panic(err)
	}
}

func GeneratePebbles(b image.Rectangle) image.Image {
	return NewScatter(
		SetBounds(b),
		SetScatterFrequency(0.04),
		SetScatterDensity(1.0),
		SetScatterMaxOverlap(1),
		SetScatterGenerator(func(u, v float64, hash uint64) (color.Color, float64) {
			rSize := float64(hash&0xFF)/255.0
			radius := 12.0 + rSize*6.0
			angle := math.Atan2(v, u)
			dist := math.Sqrt(u*u + v*v)
			noise := math.Sin(angle*5 + float64(hash%10)) * 0.1
			noise += math.Cos(angle*13 + float64(hash%7)) * 0.05
			effectiveRadius := radius * (1.0 + noise)
			if dist > effectiveRadius {
				return color.Transparent, 0
			}
			grey := 100 + int(hash%100)
			col := color.RGBA{uint8(grey), uint8(grey - 5), uint8(grey - 10), 255}
			nx := u / effectiveRadius
			ny := v / effectiveRadius
			nz := math.Sqrt(math.Max(0, 1.0 - nx*nx - ny*ny))
			lx, ly, lz := -0.5, -0.5, 0.7
			lLen := math.Sqrt(lx*lx + ly*ly + lz*lz)
			lx, ly, lz = lx/lLen, ly/lLen, lz/lLen
			diffuse := math.Max(0, nx*lx + ny*ly + nz*lz)
			r := float64(col.R) * (0.1 + 0.9*diffuse)
			g := float64(col.G) * (0.1 + 0.9*diffuse)
			b := float64(col.B) * (0.1 + 0.9*diffuse)
			alpha := 1.0
			edgeDist := effectiveRadius - dist
			if edgeDist < 1.0 {
				alpha = edgeDist
			}

			// Use hash for random Z-ordering
			z := float64(hash) / 18446744073709551615.0

			return color.RGBA{
				R: uint8(math.Min(255, r)),
				G: uint8(math.Min(255, g)),
				B: uint8(math.Min(255, b)),
				A: uint8(alpha * 255),
			}, z
		}),
	)
}

func init() {
	RegisterGenerator(PebblesBaseLabel, GeneratePebbles)
}
