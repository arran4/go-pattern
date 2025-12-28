package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure WorleyNoise implements the image.Image interface.
var _ image.Image = (*WorleyNoise)(nil)

type DistanceMetric int

const (
	MetricEuclidean DistanceMetric = iota
	MetricManhattan
	MetricChebyshev
)

type WorleyOutput int

const (
	OutputF1 WorleyOutput = iota
	OutputF2
	OutputF2MinusF1
	OutputCellID
)

// WorleyNoise generates cellular noise (Worley noise).
type WorleyNoise struct {
	Null
	Seed
	Frequency
	Jitter float64
	Metric DistanceMetric
	Output WorleyOutput
}

func (w *WorleyNoise) At(x, y int) color.Color {
	freq := w.Frequency.Frequency
	if freq == 0 {
		freq = 0.05 // Default frequency
	}
	nx, ny := float64(x)*freq, float64(y)*freq
	ix, iy := math.Floor(nx), math.Floor(ny)
	fx, fy := nx-ix, ny-iy

	minDist := math.MaxFloat64
	secondMinDist := math.MaxFloat64
	closestHash := uint64(0)

	// Check 3x3 neighbor grids
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			neighborX := int(ix) + dx
			neighborY := int(iy) + dy

			// Hash to find point in neighbor cell
			h := w.hash(neighborX, neighborY)

			// Extract point position from hash
			// Use different bits for X and Y to decorrelate
			rX := float64(h&0xFFFF) / 65535.0
			rY := float64((h>>16)&0xFFFF) / 65535.0

			pointX := float64(dx) + rX*w.Jitter
			pointY := float64(dy) + rY*w.Jitter

			// Calculate distance
			var dist float64
			diffX := math.Abs(pointX - fx)
			diffY := math.Abs(pointY - fy)

			switch w.Metric {
			case MetricManhattan:
				dist = diffX + diffY
			case MetricChebyshev:
				dist = math.Max(diffX, diffY)
			case MetricEuclidean:
				fallthrough
			default:
				dist = math.Sqrt(diffX*diffX + diffY*diffY)
			}

			if dist < minDist {
				secondMinDist = minDist
				minDist = dist
				closestHash = h
			} else if dist < secondMinDist {
				secondMinDist = dist
			}
		}
	}

	var val float64
	switch w.Output {
	case OutputF1:
		val = minDist
	case OutputF2:
		val = secondMinDist
	case OutputF2MinusF1:
		val = secondMinDist - minDist
	case OutputCellID:
		// Map hash to grayscale color
		c := uint8(closestHash & 0xFF)
		return color.Gray{Y: c}
	}

	// Clamp and map to grayscale
	if val < 0 {
		val = 0
	}
	if val > 1 {
		val = 1
	}
	return color.Gray{Y: uint8(val * 255)}
}

// hash is a stateless hash function based on coordinates and seed.
func (w *WorleyNoise) hash(x, y int) uint64 {
	return StableHash(x, y, uint64(w.Seed.Seed))
}

// NewWorleyNoise creates a new WorleyNoise pattern.
func NewWorleyNoise(ops ...func(any)) image.Image {
	w := &WorleyNoise{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Jitter: 1.0, // Default full jitter
	}
	// Defaults
	w.Frequency.Frequency = 0.05

	for _, op := range ops {
		op(w)
	}
	return w
}

// Configuration options specific to WorleyNoise

// SetWorleyMetric sets the distance metric.
func SetWorleyMetric(m DistanceMetric) func(any) {
	return func(i any) {
		if w, ok := i.(*WorleyNoise); ok {
			w.Metric = m
		}
	}
}

// SetWorleyOutput sets the output type.
func SetWorleyOutput(o WorleyOutput) func(any) {
	return func(i any) {
		if w, ok := i.(*WorleyNoise); ok {
			w.Output = o
		}
	}
}

// SetWorleyJitter sets the jitter amount (0.0 to 1.0).
func SetWorleyJitter(j float64) func(any) {
	return func(i any) {
		if w, ok := i.(*WorleyNoise); ok {
			w.Jitter = j
		}
	}
}
