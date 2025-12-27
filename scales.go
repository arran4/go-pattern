package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Scales implements the image.Image interface.
var _ image.Image = (*Scales)(nil)

// Scales generates an overlapping scales pattern.
type Scales struct {
	Null
	Radius  int
	SpacingX int
	SpacingY int
}

func (s *Scales) At(x, y int) color.Color {
	r := s.Radius
	sx := s.SpacingX
	sy := s.SpacingY
	if r <= 0 { r = 20 }
	if sx <= 0 { sx = r }
	if sy <= 0 { sy = r }

	// Calculate grid cell indices
	// Rows are spaced by sy
	// Columns are spaced by sx
	// But odd rows are offset by sx/2

	// We need to check neighbors because scales overlap.
	// We want the "highest" z-index scale to win.
	// Let's assume Z increases with Y (lower rows overlap upper rows).
	// Within a row, maybe right overlaps left? Or no overlap in row?
	// Fish scales usually: Row Y overlaps Row Y-1.

	// Check a range of cells around the pixel.
	// Rough approximation: Convert pixel to approx cell coord, check +/- 1 or 2.

	// cy := int(math.Floor(float64(y) / float64(sy))) // Removed unused variable

	// We need to check sufficient rows above and below depending on radius vs spacing.
	// Max overlap distance is Radius.
	// So we check rows where abs(center_y - y) < Radius.

	minRow := int(math.Floor(float64(y - r) / float64(sy)))
	maxRow := int(math.Ceil(float64(y + r) / float64(sy)))

	bestZ := -100000
	var bestVal float64 = 0 // 0 to 1 (center)
	found := false

	for row := minRow; row <= maxRow; row++ {
		// Calculate X offset for this row
		offsetX := 0.0
		if row % 2 != 0 {
			offsetX = float64(sx) / 2.0
		}

		// Center Y of this row
		centerY := float64(row) * float64(sy)

		// Determine relevant columns
		// abs(centerX - x) < Radius
		minCol := int(math.Floor((float64(x) - float64(r) - offsetX) / float64(sx)))
		maxCol := int(math.Ceil((float64(x) + float64(r) - offsetX) / float64(sx)))

		for col := minCol; col <= maxCol; col++ {
			centerX := float64(col)*float64(sx) + offsetX

			// Distance check
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distSq := dx*dx + dy*dy
			radSq := float64(r * r)

			if distSq < radSq {
				// This scale covers the pixel.
				// Determine Z-index.
				// Higher row = higher Z.
				// If same row, maybe secondary Z on col?
				// Let's say higher Y wins.
				z := row

				if z > bestZ {
					bestZ = z
					// Calculate height profile
					// Simple spherical: sqrt(1 - d^2/r^2)
					// Conical: 1 - d/r
					// Let's use spherical for a nice curve
					// dist := math.Sqrt(distSq) // Removed unused variable
					val := math.Sqrt(1.0 - (distSq / radSq))
					// Or maybe just 1 - dist/r
					// val = 1.0 - (dist / float64(r))

					bestVal = val
					found = true
				}
			}
		}
	}

	if !found {
		return color.Black
	}

	// Map 0-1 to Grayscale
	c := uint8(bestVal * 255)
	return color.Gray{Y: c}
}

// NewScales creates a new Scales pattern.
func NewScales(ops ...func(any)) image.Image {
	p := &Scales{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Radius: 30,
		SpacingX: 30,
		SpacingY: 20, // Default overlap in Y
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// Configuration options

type ScaleRadius struct { Radius int }
func (s *ScaleRadius) SetScaleRadius(v int) { s.Radius = v }
type hasScaleRadius interface { SetScaleRadius(int) }
func SetScaleRadius(v int) func(any) {
	return func(i any) { if h, ok := i.(hasScaleRadius); ok { h.SetScaleRadius(v) } }
}

type ScaleXSpacing struct { SpacingX int }
func (s *ScaleXSpacing) SetScaleXSpacing(v int) { s.SpacingX = v }
type hasScaleXSpacing interface { SetScaleXSpacing(int) }
func SetScaleXSpacing(v int) func(any) {
	return func(i any) { if h, ok := i.(hasScaleXSpacing); ok { h.SetScaleXSpacing(v) } }
}

type ScaleYSpacing struct { SpacingY int }
func (s *ScaleYSpacing) SetScaleYSpacing(v int) { s.SpacingY = v }
type hasScaleYSpacing interface { SetScaleYSpacing(int) }
func SetScaleYSpacing(v int) func(any) {
	return func(i any) { if h, ok := i.(hasScaleYSpacing); ok { h.SetScaleYSpacing(v) } }
}

// Implement interfaces on Scales struct
func (s *Scales) SetScaleRadius(v int) { s.Radius = v }
func (s *Scales) SetScaleXSpacing(v int) { s.SpacingX = v }
func (s *Scales) SetScaleYSpacing(v int) { s.SpacingY = v }
