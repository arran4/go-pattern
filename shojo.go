package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Shojo implements the image.Image interface.
var _ image.Image = (*Shojo)(nil)

// Shojo is a pattern that generates scattered starbursts with glow halos ("Shōjo Sparkles").
type Shojo struct {
	Null
	FillColor  // Sparkle color
	SpaceColor // Background color
	Seed       int64
}

// hash returns a pseudo-random float64 in [0, 1) based on coordinates and seed.
func hash(x, y int, seed int64) float64 {
	// Simple hashing based on integer coordinates
	var h = uint64(seed) ^ (uint64(x)*73856093) ^ (uint64(y)*19349663)
	h ^= h >> 15
	h *= 0x735a2d97a6428a11
	h ^= h >> 15
	return float64(h&0xffffff) / 16777215.0
}

// At returns the color of the pixel at (x, y).
// It implements a grid-based scattered point algorithm to render starbursts.
func (p *Shojo) At(x, y int) color.Color {
	const (
		gridSize = 60
		// Influence radius of a star.
		maxRadius = 50.0
	)

	bgR, bgG, bgB, bgA := p.SpaceColor.SpaceColor.RGBA()

	// Base color accumulator
	accR := float64(bgR) / 65535.0
	accG := float64(bgG) / 65535.0
	accB := float64(bgB) / 65535.0
	accA := float64(bgA) / 65535.0

	sparkleR, sparkleG, sparkleB, sparkleA := p.FillColor.FillColor.RGBA()
	sR := float64(sparkleR) / 65535.0
	sG := float64(sparkleG) / 65535.0
	sB := float64(sparkleB) / 65535.0
	sA := float64(sparkleA) / 65535.0

	// Determine grid cell
	gx := int(math.Floor(float64(x) / float64(gridSize)))
	gy := int(math.Floor(float64(y) / float64(gridSize)))

	// Check 3x3 neighbor cells
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			cx := gx + dx
			cy := gy + dy

			// Deterministic random values using hash function
			r1 := hash(cx, cy, p.Seed)
			r2 := hash(cx, cy, p.Seed+1)
			r3 := hash(cx, cy, p.Seed+2)
			r4 := hash(cx, cy, p.Seed+3)
			r5 := hash(cx, cy, p.Seed+4)

			// Random number of sparkles per cell (0 or 1 for simplicity and spacing)
			if r1 > 0.4 { // 40% chance of a star (inverted check from previous > 0.6)
				continue
			}

			// Position within the cell
			starX := float64(cx*gridSize) + r2*float64(gridSize)
			starY := float64(cy*gridSize) + r3*float64(gridSize)

			// Properties
			size := 10.0 + r4*20.0
			rotation := r5 * math.Pi

			// Calculate intensity at pixel (x,y)
			distX := float64(x) - starX
			distY := float64(y) - starY
			dist := math.Sqrt(distX*distX + distY*distY)

			if dist >= maxRadius {
				continue
			}

			// Angle
			angle := math.Atan2(distY, distX) + rotation

			// Fade out at edges to avoid hard cut
			fade := 1.0 - (dist / maxRadius)
			if fade < 0 {
				fade = 0
			}
			// Apply a smoother fade curve (smoothstep)
			fade = fade * fade * (3 - 2*fade)

			// Starburst shape function
			// 1. Core glow (1/dist)
			d := math.Max(dist, 0.5)

			// Soft glow
			glow := 0.2 * size / d

			// Rays: 4 main rays
			rayWidth := 20.0
			rayIntensity := math.Pow(math.Abs(math.Sin(angle*2)), rayWidth)

			// Attenuate rays with distance more sharply than glow
			rays := rayIntensity * (size / d)

			intensity := (glow + rays) * fade

			// Additive blending
			accR += sR * intensity * sA
			accG += sG * intensity * sA
			accB += sB * intensity * sA
			// Accumulate alpha based on intensity
			accA = math.Max(accA, sA * math.Min(1.0, intensity))
		}
	}

	// Clamp values
	accR = math.Min(accR, 1.0)
	accG = math.Min(accG, 1.0)
	accB = math.Min(accB, 1.0)
	accA = math.Min(accA, 1.0)

	return color.RGBA64{
		R: uint16(accR * 65535),
		G: uint16(accG * 65535),
		B: uint16(accB * 65535),
		A: uint16(accA * 65535),
	}
}

func (p *Shojo) Bounds() image.Rectangle {
	return p.bounds
}

func (p *Shojo) ColorModel() color.Model {
	return color.RGBA64Model
}

// NewShojo creates a new Shojo Sparkles pattern.
func NewShojo(ops ...func(any)) image.Image {
	p := &Shojo{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Seed: 1,
	}
	p.FillColor.FillColor = color.White // Default sparkle color
	p.SpaceColor.SpaceColor = color.RGBA{10, 0, 30, 255} // Default dark purple background

	for _, op := range ops {
		op(p)
	}
	return p
}
