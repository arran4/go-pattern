package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Voronoi implements the image.Image interface.
var _ image.Image = (*Voronoi)(nil)

// Voronoi is a pattern that generates Voronoi cells based on a set of points and colors.
type Voronoi struct {
	Null
	Points []image.Point
	Colors []color.Color
}

func (v *Voronoi) At(x, y int) color.Color {
	if len(v.Points) == 0 {
		return color.Transparent
	}

	minDist := math.MaxFloat64
	closestIndex := 0

	for i, p := range v.Points {
		dx := float64(x - p.X)
		dy := float64(y - p.Y)
		dist := dx*dx + dy*dy

		if dist < minDist {
			minDist = dist
			closestIndex = i
		}
	}

	if len(v.Colors) > 0 {
		return v.Colors[closestIndex%len(v.Colors)]
	}
	return color.Black
}

// NewVoronoi creates a new Voronoi pattern.
func NewVoronoi(points []image.Point, colors []color.Color, ops ...func(any)) image.Image {
	p := &Voronoi{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Points: points,
		Colors: colors,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoVoronoi produces a demo variant for readme.md pre-populated values.
func NewDemoVoronoi(ops ...func(any)) image.Image {
	points := []image.Point{
		{50, 50}, {205, 50},
		{127, 127},
		{50, 205}, {205, 205},
	}
	colors := []color.Color{
		color.RGBA{255, 100, 100, 255},
		color.RGBA{100, 255, 100, 255},
		color.RGBA{100, 100, 255, 255},
		color.RGBA{255, 255, 100, 255},
		color.RGBA{100, 255, 255, 255},
	}
	return NewVoronoi(points, colors, ops...)
}
