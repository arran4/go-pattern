package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Road implements image.Image
var _ image.Image = (*Road)(nil)

type RoadShape int

const (
	RoadStraight RoadShape = iota
	RoadCurved
	RoadIntersection
	RoadTIntersection
)

type Road struct {
	Null
	Lanes        int
	Width        float64
	Direction    float64 // Degrees
	Shape        RoadShape
	CurveRadius  float64
	CurveAngle   float64 // Degrees
	Asphalt      image.Image
	Background   image.Image
	MarkingColor color.Color
	Center
}

// SetLanes sets the number of lanes.
func SetLanes(l int) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.Lanes = l
		}
	}
}

// SetWidth sets the total width of the road.
func SetWidth(w float64) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.Width = w
		}
	}
}

// SetDirection sets the rotation of the road in degrees.
func SetDirection(d float64) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.Direction = d
		}
	}
}

// SetShape sets the shape of the road.
func SetShape(s RoadShape) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.Shape = s
		}
	}
}

// SetCurveRadius sets the radius for curved roads.
func SetCurveRadius(rad float64) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.CurveRadius = rad
		}
	}
}

// SetCurveAngle sets the angle of the curve in degrees.
func SetCurveAngle(a float64) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.CurveAngle = a
		}
	}
}

// SetAsphalt sets the texture for the road surface.
func SetAsphalt(img image.Image) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.Asphalt = img
		}
	}
}

// SetBackground sets the texture for the non-road area.
func SetBackground(img image.Image) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.Background = img
		}
	}
}

// SetMarkingColor sets the color of the road markings.
func SetMarkingColor(c color.Color) func(any) {
	return func(p any) {
		if r, ok := p.(*Road); ok {
			r.MarkingColor = c
		}
	}
}

func NewRoad(ops ...func(any)) image.Image {
	r := &Road{
		Null:         Null{bounds: image.Rect(0, 0, 255, 255)},
		Lanes:        2,
		Width:        100,
		Direction:    0,
		Shape:        RoadStraight,
		CurveRadius:  100,
		CurveAngle:   90,
		MarkingColor: color.White,
		Center:       Center{127, 127},
	}
	for _, op := range ops {
		op(r)
	}
	return r
}

func (r *Road) At(x, y int) color.Color {
	// 1. Transform to local coordinates relative to Center and rotated by Direction.
	dx := float64(x - r.CenterX)
	dy := float64(y - r.CenterY)

	rad := r.Direction * math.Pi / 180.0
	cos := math.Cos(-rad)
	sin := math.Sin(-rad)

	rx := dx*cos - dy*sin
	ry := dx*sin + dy*cos

	// 2. Determine if we are on the road.
	onRoad := false
	u, v := 0.0, 0.0 // curvilinear coordinates: u = longitudinal, v = lateral (-Width/2 to Width/2)

	halfW := r.Width / 2.0

	switch r.Shape {
	case RoadStraight:
		// Horizontal road along X axis
		if math.Abs(ry) <= halfW {
			onRoad = true
			u = rx
			v = ry
		}
	case RoadIntersection:
		// Union of Horizontal and Vertical
		// Horizontal
		if math.Abs(ry) <= halfW {
			onRoad = true
			u = rx
			v = ry
		} else if math.Abs(rx) <= halfW {
			// Vertical
			onRoad = true
			u = ry // Swap for vertical so markings run along it?
			v = rx
			// Note: Intersection area is tricky for markings.
			// Simple union: The markings might overlap weirdly.
			// For "macro view", overlapping markings is okay or we can prioritize one.
		}
	case RoadTIntersection:
		// Horizontal (Through) + Vertical (Stem) coming from +Y?
		// Let's say Horizontal is the top bar of T. Vertical comes from bottom.
		// Horizontal: abs(ry) <= halfW
		// Vertical: abs(rx) <= halfW AND ry >= 0 (or ry >= -halfW to merge)
		if math.Abs(ry) <= halfW {
			onRoad = true
			u = rx
			v = ry
		} else if math.Abs(rx) <= halfW && ry >= 0 {
			onRoad = true
			u = ry
			v = rx
		}
	case RoadCurved:
		// Arc. Center of curvature at (0, R) relative to the start point?
		// Let's define the curve as: center at (0, CurveRadius).
		// Road is annulus with radius [CurveRadius - halfW, CurveRadius + halfW].
		// Angle range needs handling.
		// Let's assume the curve is centered at (0, r.CurveRadius) in local coords?
		// No, usually "Center" is the pivot.
		// Let's say Center is the center of the arc.
		// Then dist = sqrt(rx^2 + ry^2).
		// if dist in [Radius - halfW, Radius + halfW] -> on road.
		// And check angle.

		// Let's shift so that (0,0) is the center of the road arc.
		// But `r.Center` is the image placement.
		// So `rx, ry` are relative to arc center.

		dist := math.Sqrt(rx*rx + ry*ry)
		if dist >= r.CurveRadius-halfW && dist <= r.CurveRadius+halfW {
			// Check angle
			angle := math.Atan2(ry, rx) * 180.0 / math.Pi
			// Normalize angle to [0, 360) or similar.
			// Let's assume curve is valid within some angle range.
			// If CurveAngle is 90, maybe from -45 to 45? or 0 to 90?
			// Let's assume -CurveAngle/2 to +CurveAngle/2 for symmetry around X axis?
			// Or 0 to CurveAngle.
			// Let's go with symmetry around Y axis (top)?
			// Simpler: 0 to CurveAngle.
			// Need to normalize angle.
			// Atan2 returns -180 to 180.

			// Let's simplify: Full circle if CurveAngle >= 360.
			if r.CurveAngle >= 360 || (angle >= -r.CurveAngle/2 && angle <= r.CurveAngle/2) {
				onRoad = true
				u = angle * r.CurveRadius * (math.Pi / 180.0) // Arc length
				v = dist - r.CurveRadius
			}
		}
	}

	if onRoad {
		// Draw Road
		var c color.Color

		// Surface
		if r.Asphalt != nil {
			c = r.Asphalt.At(x, y)
		} else {
			c = color.RGBA{60, 60, 60, 255} // Default dark grey
		}

		// Markings
		// v is distance from center line.
		// Lanes setup.
		// Lane width = Width / Lanes
		laneWidth := r.Width / float64(r.Lanes)

		// Line width
		lineWidth := 2.0 // Fixed for now, or configurable?

		// Check if we are on a line
		// Lines are at v = -Width/2 + i * laneWidth
		// i goes from 0 to Lanes.

		// Map v to lane coordinates
		// Shift v to start from edge: v' = v + Width/2
		vp := v + halfW

		// Check distance to nearest line
		// i = round(vp / laneWidth)
		// dist = abs(vp - i * laneWidth)

		i := math.Round(vp / laneWidth)
		distToLine := math.Abs(vp - i * laneWidth)

		if distToLine < lineWidth/2 {
			// It is a line.
			// Determine line type.
			// Edge lines (i == 0 or i == Lanes): Solid
			// Middle lines: Dashed?
			// Center line (if Lanes is even, i == Lanes/2): Double Solid? or Dashed?

			isEdge := (i == 0 || int(i) == r.Lanes)

			if isEdge {
				// Solid edge line
				c = r.MarkingColor
			} else {
				// Inner line.
				// Dash pattern. Based on u.
				// Dash length e.g. 20.
				dashLen := 20.0
				gapLen := 20.0
				cycle := dashLen + gapLen

				// Determine phase
				phase := math.Mod(math.Abs(u), cycle)
				if phase < dashLen {
					c = r.MarkingColor
				}
			}
		}

		return c
	}

	// Background
	if r.Background != nil {
		return r.Background.At(x, y)
	}

	// Default grass-like color
	return color.RGBA{34, 139, 34, 255}
}

func (r *Road) Bounds() image.Rectangle {
	return r.bounds
}

func (r *Road) ColorModel() color.Model {
	return color.RGBAModel
}
