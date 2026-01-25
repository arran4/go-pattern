package pattern

import (
	"image"
	"image/color"
	"math"
	"sort"
	"sync"
)

// Yliluoma1Dither implements Yliluoma's ordered dithering algorithm 1.
// It mixes two colors from the palette to approximate the input color.
type Yliluoma1Dither struct {
	Null
	Input       image.Image
	Palette     []color.Color
	Matrix      []float64 // 0..1 values
	Size        int       // e.g. 8 for 8x8
	cache       map[int]mixingPlan
	paletteRGBA []cachedRGBA
	mu          sync.RWMutex
}

type cachedRGBA struct {
	r, g, b int
}

// NewYliluoma1Dither creates a new Yliluoma1Dither pattern.
func NewYliluoma1Dither(input image.Image, palette color.Palette, size int, ops ...func(any)) image.Image {
	if palette == nil {
		palette = color.Palette{color.Black, color.White}
	}
	if size <= 0 {
		size = 8
	}

	// Use standard Bayer matrix
	var mat []float64
	if size == 8 {
		mat = Bayer8x8
	} else if size == 4 {
		mat = Bayer4x4
	} else if size == 2 {
		mat = Bayer2x2
	} else {
		// Generate custom if needed, fallback to 8x8 for now or implement generator
		mat = Bayer8x8
		size = 8
	}

	p := &Yliluoma1Dither{
		Null: Null{
			bounds: image.Rect(0, 0, 100, 100),
		},
		Input:       input,
		Palette:     palette,
		Matrix:      mat,
		Size:        size,
		cache:       make(map[int]mixingPlan),
		paletteRGBA: make([]cachedRGBA, len(palette)),
	}
	for i, c := range palette {
		r, g, b, _ := c.RGBA()
		p.paletteRGBA[i] = cachedRGBA{int(r >> 8), int(g >> 8), int(b >> 8)}
	}
	if input != nil {
		p.bounds = input.Bounds()
	}

	for _, op := range ops {
		op(p)
	}
	return p
}

func (p *Yliluoma1Dither) At(x, y int) color.Color {
	if p.Input == nil {
		return color.Black
	}
	c := p.Input.At(x, y)
	r, g, b, _ := c.RGBA() // 0-65535
	// Normalize to 0-255 for internal calc matches article logic better
	ri, gi, bi := int(r>>8), int(g>>8), int(b>>8)

	// Matrix value 0..1
	mx := x % p.Size
	if mx < 0 {
		mx += p.Size
	}
	my := y % p.Size
	if my < 0 {
		my += p.Size
	}
	mapValue := p.Matrix[my*p.Size+mx]

	// Find best mixing plan
	key := (ri << 16) | (gi << 8) | bi

	p.mu.RLock()
	plan, ok := p.cache[key]
	p.mu.RUnlock()

	if !ok {
		plan = p.deviseBestMixingPlan(ri, gi, bi)
		p.mu.Lock()
		p.cache[key] = plan
		p.mu.Unlock()
	}

	if mapValue < plan.ratio {
		return p.Palette[plan.index2]
	}
	return p.Palette[plan.index1]
}

type mixingPlan struct {
	index1, index2 int
	ratio          float64 // 0..1. 0 = always index1, 1 = always index2
}

func (p *Yliluoma1Dither) deviseBestMixingPlan(r, g, b int) mixingPlan {
	if p.paletteRGBA == nil || len(p.paletteRGBA) != len(p.Palette) {
		p.paletteRGBA = make([]cachedRGBA, len(p.Palette))
		for i, c := range p.Palette {
			r, g, b, _ := c.RGBA()
			p.paletteRGBA[i] = cachedRGBA{int(r >> 8), int(g >> 8), int(b >> 8)}
		}
	}

	bestPlan := mixingPlan{0, 0, 0.5}
	leastPenalty := 1e99

	targetLuma := calculateLuma(r, g, b)

	// Iterate all unique pairs
	for i := 0; i < len(p.Palette); i++ {
		for j := i; j < len(p.Palette); j++ {
			c1 := p.paletteRGBA[i]
			c2 := p.paletteRGBA[j]
			ir1, ig1, ib1 := c1.r, c1.g, c1.b
			ir2, ig2, ib2 := c2.r, c2.g, c2.b

			// Analytic ratio calculation
			// solve(r1 + ratio*(r2-r1) = r) => ratio = (r - r1) / (r2 - r1)
			// Weighted average of ratios per channel

			// If colors are same, ratio doesn't matter, pick 0.5
			if i == j {
				penalty := evaluateMixingError(targetLuma, r, g, b, ir1, ig1, ib1, ir1, ig1, ib1, ir2, ig2, ib2, 0)
				if penalty < leastPenalty {
					leastPenalty = penalty
					bestPlan = mixingPlan{i, j, 0} // ratio irrelevant
				}
				continue
			}

			// Calculate optimal ratio
			// ratio = ( (r-r1)*(r2-r1) + (g-g1)*(g2-g1) + ... ) / ( (r2-r1)^2 + ... )
			// This is projecting vector (Target-C1) onto (C2-C1)

			dr := float64(ir2 - ir1)
			dg := float64(ig2 - ig1)
			db := float64(ib2 - ib1)

			numer := float64(r-ir1)*dr + float64(g-ig1)*dg + float64(b-ib1)*db
			denom := dr*dr + dg*dg + db*db

			ratio := 0.5
			if denom != 0 {
				ratio = numer / denom
			}

			if ratio < 0 {
				ratio = 0
			} else if ratio > 1 {
				ratio = 1
			}

			// Calculate mixed color
			mr := float64(ir1) + ratio*dr
			mg := float64(ig1) + ratio*dg
			mb := float64(ib1) + ratio*db

			penalty := evaluateMixingError(targetLuma, r, g, b, int(mr), int(mg), int(mb), ir1, ig1, ib1, ir2, ig2, ib2, ratio)
			if penalty < leastPenalty {
				leastPenalty = penalty
				bestPlan = mixingPlan{i, j, ratio}
			}
		}
	}
	return bestPlan
}

func evaluateMixingError(targetLuma float64, r, g, b, r0, g0, b0, r1, g1, b1, r2, g2, b2 int, ratio float64) float64 {
	// Using the improved comparison from article
	// ColorCompare(r,g,b, r0,g0,b0) + ColorCompare(r1,g1,b1, r2,g2,b2) * 0.1 * (fabs(ratio-0.5)+0.5);

	luma2 := (float64(r0)*299 + float64(g0)*587 + float64(b0)*114) / (255.0 * 1000)
	lumadiff := targetLuma - luma2

	diffR := float64(r - r0) / 255.0
	diffG := float64(g - g0) / 255.0
	diffB := float64(b - b0) / 255.0

	baseErr := (diffR*diffR*0.299+diffG*diffG*0.587+diffB*diffB*0.114)*0.75 + lumadiff*lumadiff

	mixErr := colorCompareLuma(r1, g1, b1, r2, g2, b2)

	factor := math.Abs(ratio - 0.5) + 0.5
	return baseErr + mixErr*0.1*factor
}

func calculateLuma(r, g, b int) float64 {
	return (float64(r)*299 + float64(g)*587 + float64(b)*114) / (255.0 * 1000)
}

func colorCompareLuma(r1, g1, b1, r2, g2, b2 int) float64 {
	luma1 := (float64(r1)*299 + float64(g1)*587 + float64(b1)*114) / (255.0 * 1000)
	luma2 := (float64(r2)*299 + float64(g2)*587 + float64(b2)*114) / (255.0 * 1000)
	lumadiff := luma1 - luma2

	diffR := float64(r1 - r2) / 255.0
	diffG := float64(g1 - g2) / 255.0
	diffB := float64(b1 - b2) / 255.0

	return (diffR*diffR*0.299 + diffG*diffG*0.587 + diffB*diffB*0.114)*0.75 + lumadiff*lumadiff
}

// Yliluoma2Dither implements Yliluoma's ordered dithering algorithm 2.
// It builds a candidate list of colors that average to the input color.
type Yliluoma2Dither struct {
	Null
	Input   image.Image
	Palette []color.Color
	Matrix  []float64
	Size    int
	cache   map[int][]int
	mu      sync.RWMutex
}

// NewYliluoma2Dither creates a new Yliluoma2Dither pattern.
func NewYliluoma2Dither(input image.Image, palette color.Palette, size int, ops ...func(any)) image.Image {
	if palette == nil {
		palette = color.Palette{color.Black, color.White}
	}
	if size <= 0 {
		size = 8
	}
	// ... reuse matrix selection logic ...
	var mat []float64
	if size == 8 {
		mat = Bayer8x8
	} else if size == 4 {
		mat = Bayer4x4
	} else if size == 2 {
		mat = Bayer2x2
	} else {
		mat = Bayer8x8
		size = 8
	}

	p := &Yliluoma2Dither{
		Null: Null{
			bounds: image.Rect(0, 0, 100, 100),
		},
		Input:   input,
		Palette: palette,
		Matrix:  mat,
		Size:    size,
		cache:   make(map[int][]int),
	}
	if input != nil {
		p.bounds = input.Bounds()
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

func (p *Yliluoma2Dither) At(x, y int) color.Color {
	if p.Input == nil {
		return color.Black
	}
	c := p.Input.At(x, y)
	r, g, b, _ := c.RGBA()
	ri, gi, bi := int(r>>8), int(g>>8), int(b>>8)

	mx := x % p.Size
	if mx < 0 { mx += p.Size }
	my := y % p.Size
	if my < 0 { my += p.Size }

	// Map value 0..1
	mapVal := p.Matrix[my*p.Size+mx]

	// Build candidate list
	key := (ri << 16) | (gi << 8) | bi

	p.mu.RLock()
	candidates, ok := p.cache[key]
	p.mu.RUnlock()

	if !ok {
		candidates = p.deviseMixingPlan(ri, gi, bi)
		p.mu.Lock()
		p.cache[key] = candidates
		p.mu.Unlock()
	}

	// candidates is sorted by luma
	// index = mapVal * size
	n := len(candidates)
	idx := int(mapVal * float64(n))
	if idx >= n { idx = n - 1 }
	if idx < 0 { idx = 0 }

	return p.Palette[candidates[idx]]
}

func (p *Yliluoma2Dither) deviseMixingPlan(r, g, b int) []int {
	targetCount := p.Size * p.Size // e.g. 64
	// Limit candidate list size to matrix size (or smaller for speed)
	// Article uses list size = X*Y

	candidates := make([]int, 0, targetCount)

	// Current sum
	sumR, sumG, sumB := 0, 0, 0

	for len(candidates) < targetCount {
		// Find best color to add to minimize error of average

		currentCount := len(candidates)

		bestIdx := 0
		bestCount := 1
		leastPenalty := -1.0

		maxTestCount := 1
		if currentCount > 0 {
			maxTestCount = currentCount
		}

		for i, palColor := range p.Palette {
			pr, pg, pb, _ := palColor.RGBA()
			pri, pgi, pbi := int(pr>>8), int(pg>>8), int(pb>>8)

			// Try adding 1, 2, 4... copies
			for count := 1; count <= maxTestCount; count *= 2 {
				// Test average
				tr := (sumR + pri*count)
				tg := (sumG + pgi*count)
				tb := (sumB + pbi*count)
				total := currentCount + count

				avgR := tr / total
				avgG := tg / total
				avgB := tb / total

				penalty := colorCompareLuma(r, g, b, avgR, avgG, avgB)

				if leastPenalty < 0 || penalty < leastPenalty {
					leastPenalty = penalty
					bestIdx = i
					bestCount = count
				}
			}
		}

		// Add best
		for k := 0; k < bestCount && len(candidates) < targetCount; k++ {
			candidates = append(candidates, bestIdx)
			pr, pg, pb, _ := p.Palette[bestIdx].RGBA()
			sumR += int(pr>>8)
			sumG += int(pg>>8)
			sumB += int(pb>>8)
		}
	}

	// Sort by luma
	type candidate struct {
		id   int
		luma float64
	}
	sorted := make([]candidate, len(candidates))
	for i, id := range candidates {
		sorted[i] = candidate{id, p.getLuma(id)}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].luma < sorted[j].luma
	})
	for i, s := range sorted {
		candidates[i] = s.id
	}

	return candidates
}

func (p *Yliluoma2Dither) getLuma(idx int) float64 {
	c := p.Palette[idx]
	r, g, b, _ := c.RGBA()
	return float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114
}

// KnollDither implements Thomas Knoll's pattern dithering (Photoshop).
type KnollDither struct {
	Null
	Input   image.Image
	Palette []color.Color
	Matrix  []int // Integer matrix 0..63 for 8x8
	Size    int
	cache   map[int][]int
	mu      sync.RWMutex
}

// NewKnollDither creates a new KnollDither pattern.
func NewKnollDither(img image.Image, palette color.Palette, size int, ops ...func(any)) image.Image {
	if palette == nil {
		palette = color.Palette{color.Black, color.White}
	}
	if size <= 0 {
		size = 8
	}

	// Need integer matrix 0..(N*N-1)
	mat := GenerateBayerInt(size)

	p := &KnollDither{
		Null: Null{
			bounds: image.Rect(0, 0, 100, 100),
		},
		Input:   img,
		Palette: palette,
		Matrix:  mat,
		Size:    size,
		cache:   make(map[int][]int),
	}
	if img != nil {
		p.bounds = img.Bounds()
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

func (p *KnollDither) At(x, y int) color.Color {
	if p.Input == nil {
		return color.Black
	}
	c := p.Input.At(x, y)
	r, g, b, _ := c.RGBA()
	ri, gi, bi := int(r>>8), int(g>>8), int(b>>8)

	mx := x % p.Size
	if mx < 0 { mx += p.Size }
	my := y % p.Size
	if my < 0 { my += p.Size }

	key := (ri << 16) | (gi << 8) | bi

	p.mu.RLock()
	candidates, ok := p.cache[key]
	p.mu.RUnlock()

	if !ok {
		candidates = p.devisePlan(ri, gi, bi)
		p.mu.Lock()
		p.cache[key] = candidates
		p.mu.Unlock()
	}

	idx := p.Matrix[my*p.Size + mx]

	// Protection against out of bounds if matrix has values >= len(candidates)
	if idx >= len(candidates) {
		idx = idx % len(candidates)
	}

	return p.Palette[candidates[idx]]
}

func (p *KnollDither) devisePlan(r, g, b int) []int {
	n := p.Size * p.Size
	candidates := make([]int, n)

	// Error accumulator
	er, eg, eb := 0, 0, 0

	// Threshold 0.0 .. 1.0?? Article says "Threshold = 0.5".
	// But algorithm says: Attempt = Input + Error * Threshold
	// "Error multiplier" X in code is 0.09?
	// The article pseudo code: Attempt = Input + Error * Threshold
	// Code: t = src + e * X.

	threshold := 0.5 // Default reasonable value

	for i := 0; i < n; i++ {
		// Attempt
		ar := r + int(float64(er) * threshold)
		ag := g + int(float64(eg) * threshold)
		ab := b + int(float64(eb) * threshold)

		// Clamp
		if ar < 0 { ar = 0 }; if ar > 255 { ar = 255 }
		if ag < 0 { ag = 0 }; if ag > 255 { ag = 255 }
		if ab < 0 { ab = 0 }; if ab > 255 { ab = 255 }

		// Find closest
		bestIdx := 0
		minDist := 1e99

		// "Start with guess c%16" - small optimization, ignore for now

		for pi, pc := range p.Palette {
			pr, pg, pb, _ := pc.RGBA()
			pri, pgi, pbi := int(pr>>8), int(pg>>8), int(pb>>8)

			dist := colorCompareLuma(ar, ag, ab, pri, pgi, pbi)
			if dist < minDist {
				minDist = dist
				bestIdx = pi
			}
		}

		candidates[i] = bestIdx

		// Update error
		pr, pg, pb, _ := p.Palette[bestIdx].RGBA()
		pri, pgi, pbi := int(pr>>8), int(pg>>8), int(pb>>8)

		er += r - pri
		eg += g - pgi
		eb += b - pbi
	}

	// Sort by luma
	type candidate struct {
		id   int
		luma float64
	}
	sorted := make([]candidate, len(candidates))
	for i, id := range candidates {
		sorted[i] = candidate{id, p.getLuma(id)}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].luma < sorted[j].luma
	})
	for i, s := range sorted {
		candidates[i] = s.id
	}

	return candidates
}

func (p *KnollDither) getLuma(idx int) float64 {
	c := p.Palette[idx]
	r, g, b, _ := c.RGBA()
	return float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114
}
