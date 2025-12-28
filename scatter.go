package pattern

import (
	"image"
	"image/color"
	"math"
	"sort"
)

// Ensure Scatter implements the image.Image interface.
var _ image.Image = (*Scatter)(nil)

// ScatterItem represents a single item placed by the Scatter pattern.
// It is calculated by the user-provided generator function.
type ScatterItem struct {
	Color color.Color
	Z     float64 // Depth for sorting (higher is on top)
}

// ScatterGenerator is a function that returns the color/alpha of an item at a given local coordinate.
// u, v: Local coordinates relative to the item's center (approx -0.5 to 0.5 range depending on cell).
// hash: Random hash for this specific item instance.
// returns: The color of the pixel. If Alpha is 0, it's transparent.
type ScatterGenerator func(u, v float64, hash uint64) (color.Color, float64)

// Scatter places generated items in a grid with random offsets.
// It supports overlapping items by sorting them by a Z-index derived from the hash.
type Scatter struct {
	Null
	SpaceColor         // Background color
	Frequency  float64 // Controls cell size (1/Frequency)
	Density    float64 // 0.0 to 1.0, chance of item appearing in a cell
	Generator  ScatterGenerator
	Seed       int64
	MaxOverlap int // Radius of neighbor cells to check (default 1 for 3x3)
}

func (s *Scatter) SetSeed(v int64) {
	s.Seed = v
}

func (s *Scatter) SetSeedUint64(v uint64) {
	s.Seed = int64(v)
}

func (s *Scatter) At(x, y int) color.Color {
	freq := s.Frequency
	if freq == 0 {
		freq = 0.05
	}
	cellSize := 1.0 / freq

	// Determine grid cell
	gx := int(math.Floor(float64(x) * freq))
	gy := int(math.Floor(float64(y) * freq))

	// Candidates for rendering at this pixel
	type candidate struct {
		c color.Color
		z float64
	}
	var candidates []candidate

	overlap := s.MaxOverlap
	if overlap == 0 {
		overlap = 1
	}

	// Check neighbor cells
	for dy := -overlap; dy <= overlap; dy++ {
		for dx := -overlap; dx <= overlap; dx++ {
			cx := gx + dx
			cy := gy + dy

			// Hash for this cell
			h := s.hash(cx, cy)

			// Deterministic random float [0, 1)
			r1 := float64(h&0xFFFF) / 65535.0

			// Density check
			if r1 > s.Density {
				continue
			}

			// Random position within the cell
			rX := float64((h>>16)&0xFFFF) / 65535.0
			rY := float64((h>>32)&0xFFFF) / 65535.0

			// Center of the item in pixel coordinates
			centerX := (float64(cx) + rX) * cellSize
			centerY := (float64(cy) + rY) * cellSize

			// Local coordinates (u, v) relative to center
			// Normalized so that 1.0 is roughly the size of a cell?
			// Let's pass pixel delta. Generator can decide scaling.
			u := float64(x) - centerX
			v := float64(y) - centerY

			// Call generator
			if s.Generator != nil {
				col, z := s.Generator(u, v, h)
				_, _, _, a := col.RGBA()
				if a > 0 {
					candidates = append(candidates, candidate{col, z})
				}
			}
		}
	}

	// Sort candidates by Z (low to high -> Painter's algorithm)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].z < candidates[j].z
	})

	// Composite background
	var fr, fg, fb, fa float64
	if s.SpaceColor.SpaceColor != nil {
		finalR, finalG, finalB, finalA := s.SpaceColor.SpaceColor.RGBA()
		fr = float64(finalR) / 65535.0
		fg = float64(finalG) / 65535.0
		fb = float64(finalB) / 65535.0
		fa = float64(finalA) / 65535.0
	} else {
		// Default to black if SpaceColor is nil?
		// Actually Null might not have initialized SpaceColor correctly if not embedded or set.
		// SpaceColor struct has a field SpaceColor which is an interface.
		fr, fg, fb, fa = 0, 0, 0, 1.0 // Opaque black default
	}

	// Helper for alpha blending
	blend := func(bgR, bgG, bgB, bgA, fgR, fgG, fgB, fgA float64) (float64, float64, float64, float64) {
		outA := fgA + bgA*(1.0-fgA)
		if outA == 0 {
			return 0, 0, 0, 0
		}
		outR := (fgR*fgA + bgR*bgA*(1.0-fgA)) / outA
		outG := (fgG*fgA + bgG*bgA*(1.0-fgA)) / outA
		outB := (fgB*fgA + bgB*bgA*(1.0-fgA)) / outA
		return outR, outG, outB, outA
	}

	for _, cand := range candidates {
		r, g, b, a := cand.c.RGBA()
		sr := float64(r) / 65535.0
		sg := float64(g) / 65535.0
		sb := float64(b) / 65535.0
		sa := float64(a) / 65535.0

		fr, fg, fb, fa = blend(fr, fg, fb, fa, sr, sg, sb, sa)
	}

	return color.RGBA64{
		R: uint16(fr * 65535),
		G: uint16(fg * 65535),
		B: uint16(fb * 65535),
		A: uint16(fa * 65535),
	}
}

// hash is a stateless hash function.
func (s *Scatter) hash(x, y int) uint64 {
	return StableHash(x, y, uint64(s.Seed))
}

// NewScatter creates a new Scatter pattern.
func NewScatter(ops ...func(any)) image.Image {
	p := &Scatter{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Frequency: 0.05,
		Density:   1.0,
		MaxOverlap: 1,
		Generator: func(u, v float64, hash uint64) (color.Color, float64) {
			return color.Transparent, 0
		},
	}
	// Default background black
	p.SpaceColor.SpaceColor = color.Black

	for _, op := range ops {
		op(p)
	}
	return p
}

// SetScatterFrequency sets the frequency (density of grid cells).
func SetScatterFrequency(f float64) func(any) {
	return func(i any) {
		if p, ok := i.(*Scatter); ok {
			p.Frequency = f
		}
	}
}

// SetScatterDensity sets the probability (0-1) of an item appearing in a cell.
func SetScatterDensity(d float64) func(any) {
	return func(i any) {
		if p, ok := i.(*Scatter); ok {
			p.Density = d
		}
	}
}

// SetScatterGenerator sets the item generator function.
func SetScatterGenerator(g ScatterGenerator) func(any) {
	return func(i any) {
		if p, ok := i.(*Scatter); ok {
			p.Generator = g
		}
	}
}

// SetScatterMaxOverlap sets the radius of neighbor cells to check.
func SetScatterMaxOverlap(overlap int) func(any) {
	return func(i any) {
		if p, ok := i.(*Scatter); ok {
			p.MaxOverlap = overlap
		}
	}
}
