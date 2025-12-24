package pattern

import (
	"image"
	"image/color"
)

// Ensure Heatmap implements the image.Image interface.
var _ image.Image = (*Heatmap)(nil)

// HeatmapFunc is the function signature for the heatmap generator.
// It accepts logical coordinates (x, y) and returns a scalar value z.
type HeatmapFunc func(x, y float64) float64

// Heatmap generates a color gradient based on a 2D scalar function.
type Heatmap struct {
	Null
	StartColor
	EndColor
	Func       HeatmapFunc
	MinX, MaxX float64
	MinY, MaxY float64
	MinZ, MaxZ float64
}

// At returns the color at (x, y).
func (h *Heatmap) At(x, y int) color.Color {
	b := h.Bounds()
	if b.Empty() {
		return color.RGBA{}
	}

	// Map pixel coordinates to logical coordinates
	// We map the center of the pixel or the top-left?
	// Usually mapping the coordinate directly is fine.
	// x spans from b.Min.X to b.Max.X (exclusive)
	// We want to map [0, width) to [MinX, MaxX)

	width := float64(b.Dx())
	height := float64(b.Dy())

	if width == 0 || height == 0 {
		return color.RGBA{}
	}

	// Normalized coordinates 0..1
	nx := float64(x-b.Min.X) / width
	ny := float64(y-b.Min.Y) / height

	// Logical coordinates
	u := h.MinX + nx*(h.MaxX-h.MinX)
	v := h.MinY + ny*(h.MaxY-h.MinY)

	val := h.Func(u, v)

	// Normalize val to t for interpolation
	// val between MinZ and MaxZ -> t between 0 and 1
	var t float64
	if h.MaxZ == h.MinZ {
		t = 0.5 // Avoid division by zero, return midpoint or similar
	} else {
		t = (val - h.MinZ) / (h.MaxZ - h.MinZ)
	}

	return lerpColor(h.StartColor.StartColor, h.EndColor.EndColor, t)
}

// NewHeatmap creates a new Heatmap pattern.
func NewHeatmap(f HeatmapFunc, ops ...func(any)) image.Image {
	h := &Heatmap{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Func: f,
		// Default ranges
		MinX: -1.0, MaxX: 1.0,
		MinY: -1.0, MaxY: 1.0,
		MinZ: -1.0, MaxZ: 1.0,
	}
	// Default colors
	h.StartColor.StartColor = color.Black
	h.EndColor.EndColor = color.White

	for _, op := range ops {
		op(h)
	}
	return h
}

// SetXRange sets the logical X range for the heatmap.
func SetXRange(min, max float64) func(any) {
	return func(i any) {
		if h, ok := i.(*Heatmap); ok {
			h.MinX = min
			h.MaxX = max
		}
	}
}

// SetYRange sets the logical Y range for the heatmap.
func SetYRange(min, max float64) func(any) {
	return func(i any) {
		if h, ok := i.(*Heatmap); ok {
			h.MinY = min
			h.MaxY = max
		}
	}
}

// SetZRange sets the expected Z range (value range) for the heatmap color mapping.
func SetZRange(min, max float64) func(any) {
	return func(i any) {
		if h, ok := i.(*Heatmap); ok {
			h.MinZ = min
			h.MaxZ = max
		}
	}
}
