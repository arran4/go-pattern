package pattern

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

// OrderedDither applies ordered dithering using a threshold matrix.
type OrderedDither struct {
	Null
	img     image.Image
	matrix  []float64
	dim     int
	palette color.Palette
	spread  float64
}

// NewOrderedDither creates a new OrderedDither pattern.
// matrix should be a square matrix flattened. dim is the width/height.
// values in matrix should be normalized 0..1.
// spread controls the intensity of dithering. If 0, it auto-calculates based on palette size.
func NewOrderedDither(img image.Image, matrix []float64, dim int, palette color.Palette, spread float64, ops ...func(any)) image.Image {
	if palette == nil {
		palette = color.Palette{color.Black, color.White}
	}
	if spread == 0 {
		spread = 255.0 / float64(len(palette))
		if len(palette) <= 2 {
			spread = 255.0
		}
	}

	b := image.Rect(0, 0, 100, 100)
	if img != nil {
		b = img.Bounds()
	}
	od := &OrderedDither{
		img:     img,
		matrix:  matrix,
		dim:     dim,
		palette: palette,
		spread:  spread,
		Null: Null{
			bounds: b,
		},
	}
	for _, op := range ops {
		op(od)
	}
	return od
}

func (d *OrderedDither) At(x, y int) color.Color {
	if d.img == nil {
		return color.Black
	}
	c := d.img.At(x, y)
	r, g, b, a := c.RGBA()

	mx := x % d.dim
	my := y % d.dim
	if mx < 0 {
		mx += d.dim
	}
	if my < 0 {
		my += d.dim
	}

	mVal := d.matrix[my*d.dim + mx]

	// shift = (mVal - 0.5) * spread
	shift := (mVal - 0.5) * d.spread

	rf := float64(r) / 257.0
	gf := float64(g) / 257.0
	bf := float64(b) / 257.0

	rf += shift
	gf += shift
	bf += shift

	rf = float64(clamp(rf))
	gf = float64(clamp(gf))
	bf = float64(clamp(bf))

	target := color.RGBA{
		R: uint8(rf),
		G: uint8(gf),
		B: uint8(bf),
		A: uint8(a >> 8),
	}

	return d.palette.Convert(target)
}

// Predefined Bayer matrices

func generateBayer(n int) []float64 {
	if n == 1 {
		return []float64{0}
	}

	return normalizeMatrix(GenerateBayerInt(n))
}

func GenerateBayerInt(n int) []int {
	if n == 1 {
		return []int{0}
	}
	prev := GenerateBayerInt(n / 2)
	curr := make([]int, n*n)
	half := n / 2

	for y := 0; y < half; y++ {
		for x := 0; x < half; x++ {
			val := prev[y*half + x]
			curr[y*n+x] = 4 * val
			curr[y*n+(x+half)] = 4*val + 2
			curr[(y+half)*n+x] = 4*val + 3
			curr[(y+half)*n+(x+half)] = 4*val + 1
		}
	}
	return curr
}

func normalizeMatrix(ints []int) []float64 {
	size := len(ints)
	res := make([]float64, size)
	maxVal := float64(size)
	for i, v := range ints {
		res[i] = float64(v) / maxVal
	}
	return res
}

// Bayer2x2 Matrix
var Bayer2x2 = generateBayer(2)

// Bayer4x4 Matrix
var Bayer4x4 = generateBayer(4)

// Bayer8x8 Matrix
var Bayer8x8 = generateBayer(8)

// NewBayer2x2Dither is a convenience function.
func NewBayer2x2Dither(img image.Image, palette color.Palette, ops ...func(any)) image.Image {
	return NewOrderedDither(img, Bayer2x2, 2, palette, 0, ops...)
}

// NewBayer4x4Dither is a convenience function.
func NewBayer4x4Dither(img image.Image, palette color.Palette, ops ...func(any)) image.Image {
	return NewOrderedDither(img, Bayer4x4, 4, palette, 0, ops...)
}

// NewBayer8x8Dither is a convenience function.
func NewBayer8x8Dither(img image.Image, palette color.Palette, ops ...func(any)) image.Image {
	return NewOrderedDither(img, Bayer8x8, 8, palette, 0, ops...)
}

// --- Halftone ---

// generateHalftoneMatrix generates a clustered dot (halftone) matrix.
// n is the size (e.g., 8 for 8x8).
// It creates a "hill" of thresholds peaking at the center.
func generateHalftoneMatrix(n int) []float64 {
	cx := float64(n-1) / 2.0
	cy := float64(n-1) / 2.0
	maxDist := math.Sqrt(cx*cx + cy*cy)
	if maxDist == 0 {
		return []float64{0.5}
	}

	res := make([]float64, n*n)
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			// Normalized distance 0..1
			nd := dist / maxDist
			if nd > 1 {
				nd = 1
			}

			// We want center to be high threshold (black dot in center).
			// Distance 0 -> 1.0
			// Distance Max -> 0.0
			res[y*n+x] = 1.0 - nd
		}
	}
	return res
}

// NewHalftoneDither creates a halftone dither effect using a clustered dot matrix.
// size determines the grid size of the dots (e.g. 8x8).
func NewHalftoneDither(img image.Image, size int, palette color.Palette, ops ...func(any)) image.Image {
	if size < 2 {
		size = 4
	}
	matrix := generateHalftoneMatrix(size)
	return NewOrderedDither(img, matrix, size, palette, 0, ops...)
}

// --- Random Dither ---

// RandomDither applies random noise dithering.
type RandomDither struct {
	Null
	img     image.Image
	palette color.Palette
	seed    int64
	spread  float64
}

// NewRandomDither creates a new RandomDither pattern.
func NewRandomDither(img image.Image, palette color.Palette, seed int64, ops ...func(any)) image.Image {
	if palette == nil {
		palette = color.Palette{color.Black, color.White}
	}

	// Default spread
	spread := 255.0 / float64(len(palette))
	if len(palette) <= 2 {
		spread = 255.0
	}

	b := image.Rect(0, 0, 100, 100)
	if img != nil {
		b = img.Bounds()
	}
	rd := &RandomDither{
		img:     img,
		palette: palette,
		seed:    seed,
		spread:  spread,
		Null: Null{
			bounds: b,
		},
	}
	for _, op := range ops {
		op(rd)
	}
	return rd
}

func (d *RandomDither) At(x, y int) color.Color {
	if d.img == nil {
		return color.Black
	}
	c := d.img.At(x, y)
	r, g, b, a := c.RGBA()

	// Hash x, y, seed to get random value 0..1
	// Simple hash
	h := int64(x)*48611 + int64(y)*50329 + d.seed
	h = (h ^ (h >> 16)) * 0x85ebca6b
	h = (h ^ (h >> 13)) * 0xc2b2ae35
	h = (h ^ (h >> 16))

	rnd := float64(h&0xFFFF) / 65535.0

	// Apply noise
	shift := (rnd - 0.5) * d.spread

	rf := float64(r)/257.0 + shift
	gf := float64(g)/257.0 + shift
	bf := float64(b)/257.0 + shift

	target := color.RGBA{
		R: uint8(clamp(rf)),
		G: uint8(clamp(gf)),
		B: uint8(clamp(bf)),
		A: uint8(a >> 8),
	}

	return d.palette.Convert(target)
}

// --- Blue Noise ---

// generateBlueNoise generates a small (16x16) blue noise matrix using a simplified best-candidate algorithm.
func generateBlueNoise(size int) []float64 {
	if size > 32 {
		size = 32
	}
	width := size
	height := size

	ranks := make([]int, width*height)
	for i := range ranks {
		ranks[i] = -1
	}

	// Pick first point
	r := rand.New(rand.NewSource(1))
	first := r.Intn(width * height)
	ranks[first] = 0

	points := []struct{ x, y int }{{first % width, first / width}}

	for k := 1; k < width*height; k++ {
		bestDist := -1.0
		bestIdx := -1

		// Check all candidates (pixels with rank -1)
		for i := 0; i < width*height; i++ {
			if ranks[i] != -1 {
				continue
			}
			px, py := i%width, i/width

			minDist := 1000000.0
			for _, pt := range points {
				// Toroidal distance
				dx := float64(abs(px - pt.x))
				dy := float64(abs(py - pt.y))
				if dx > float64(width)/2 {
					dx = float64(width) - dx
				}
				if dy > float64(height)/2 {
					dy = float64(height) - dy
				}
				d := dx*dx + dy*dy
				if d < minDist {
					minDist = d
				}
			}

			if minDist > bestDist {
				bestDist = minDist
				bestIdx = i
			}
		}

		ranks[bestIdx] = k
		points = append(points, struct{ x, y int }{bestIdx % width, bestIdx / width})
	}

	return normalizeMatrix(ranks)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// NewBlueNoiseDither creates a blue noise dither pattern.
// Uses a generated 16x16 Blue Noise mask.
func NewBlueNoiseDither(img image.Image, palette color.Palette, ops ...func(any)) image.Image {
	size := 16
	matrix := generateBlueNoise(size)
	return NewOrderedDither(img, matrix, size, palette, 0, ops...)
}

// --- Multi-Scale Ordered Dither ---

// MultiScaleOrderedDither blends between two matrices based on local variance.
type MultiScaleOrderedDither struct {
	Null
	img         image.Image
	palette     color.Palette
	matrixSmall []float64 // e.g. 2x2
	dimSmall    int
	matrixLarge []float64 // e.g. 8x8
	dimLarge    int
	spread      float64
}

func NewMultiScaleOrderedDither(img image.Image, palette color.Palette, ops ...func(any)) image.Image {
	if palette == nil {
		palette = color.Palette{color.Black, color.White}
	}
	spread := 255.0 / float64(len(palette))
	if len(palette) <= 2 {
		spread = 255.0
	}

	ms := &MultiScaleOrderedDither{
		img:         img,
		palette:     palette,
		matrixSmall: Bayer2x2,
		dimSmall:    2,
		matrixLarge: Bayer8x8,
		dimLarge:    8,
		spread:      spread,
	}
	b := image.Rect(0, 0, 100, 100)
	if img != nil {
		b = img.Bounds()
	}
	ms.bounds = b

	for _, op := range ops {
		op(ms)
	}
	return ms
}

func (d *MultiScaleOrderedDither) At(x, y int) color.Color {
	if d.img == nil {
		return color.Black
	}

	// Calculate local variance (simplified: difference between pixel and neighbors)
	c := d.img.At(x, y)
	r, g, b, a := c.RGBA()
	lum := (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) / 65535.0

	// Simple edge detection/variance estimate
	cR := d.img.At(x+1, y)
	cB := d.img.At(x, y+1)
	rR, gR, bR, _ := cR.RGBA()
	rB, gB, bB, _ := cB.RGBA()
	lumR := (float64(rR)*0.299 + float64(gR)*0.587 + float64(bR)*0.114) / 65535.0
	lumB := (float64(rB)*0.299 + float64(gB)*0.587 + float64(bB)*0.114) / 65535.0

	diff := math.Abs(lum-lumR) + math.Abs(lum-lumB)

	// Blend factor 0..1. 0 = Large, 1 = Small.
	t := diff * 4.0 // Sensitivity
	if t > 1 {
		t = 1
	}

	// Sample both matrices
	// Small
	mSmall := d.matrixSmall[(y%d.dimSmall)*d.dimSmall+(x%d.dimSmall)]
	// Large
	mLarge := d.matrixLarge[(y%d.dimLarge)*d.dimLarge+(x%d.dimLarge)]

	// Blend thresholds
	mVal := mLarge*(1-t) + mSmall*t

	// Apply dither
	shift := (mVal - 0.5) * d.spread

	rf := float64(r)/257.0 + shift
	gf := float64(g)/257.0 + shift
	bf := float64(b)/257.0 + shift

	target := color.RGBA{
		R: uint8(clamp(rf)),
		G: uint8(clamp(gf)),
		B: uint8(clamp(bf)),
		A: uint8(a >> 8),
	}

	return d.palette.Convert(target)
}
