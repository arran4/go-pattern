package pattern

import (
	"image"
	"image/color"
	"sort"
)

// ColorStop defines a color at a specific position (0.0 to 1.0) in the map.
type ColorStop struct {
	Position float64
	Color    color.Color
}

// ColorMap applies a color ramp to the luminance of the source image.
type ColorMap struct {
	Null
	Source image.Image
	Stops  []ColorStop
}

func (c *ColorMap) Bounds() image.Rectangle {
	if c.Source != nil {
		return c.Source.Bounds()
	}
	return c.Null.Bounds()
}

// At returns the color at (x, y).
func (c *ColorMap) At(x, y int) color.Color {
	if c.Source == nil {
		return color.Black
	}

	// Get source color
	srcColor := c.Source.At(x, y)

	// Convert to grayscale intensity [0, 1]
	// We use 16-bit grayscale for better precision to avoid banding
	gray := color.Gray16Model.Convert(srcColor).(color.Gray16)
	t := float64(gray.Y) / 65535.0

	// Find the stops surrounding t
	if len(c.Stops) == 0 {
		return srcColor // No mapping, return original
	}

	// Handle edge cases
	if t <= c.Stops[0].Position {
		return c.Stops[0].Color
	}
	if t >= c.Stops[len(c.Stops)-1].Position {
		return c.Stops[len(c.Stops)-1].Color
	}

	// Binary search or linear scan? Linear is fine for small number of stops.
	for i := 0; i < len(c.Stops)-1; i++ {
		s1 := c.Stops[i]
		s2 := c.Stops[i+1]
		if t >= s1.Position && t <= s2.Position {
			// Interpolate
			ratio := (t - s1.Position) / (s2.Position - s1.Position)
			return lerpColor(s1.Color, s2.Color, ratio)
		}
	}

	return c.Stops[len(c.Stops)-1].Color
}

// NewColorMap creates a new ColorMap pattern.
// It sorts the stops by position.
func NewColorMap(source image.Image, stops ...ColorStop) image.Image {
	// Sort stops
	sortedStops := make([]ColorStop, len(stops))
	copy(sortedStops, stops)
	sort.Slice(sortedStops, func(i, j int) bool {
		return sortedStops[i].Position < sortedStops[j].Position
	})

	b := image.Rect(0, 0, 255, 255)
	if source != nil {
		b = source.Bounds()
	}

	return &ColorMap{
		Null: Null{
			bounds: b,
		},
		Source: source,
		Stops:  sortedStops,
	}
}
