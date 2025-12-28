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
		Width:        120,
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

// smoothstep for anti-aliasing
func smoothstep(edge0, edge1, x float64) float64 {
	t := (x - edge0) / (edge1 - edge0)
	if t < 0.0 {
		return 0.0
	}
	if t > 1.0 {
		return 1.0
	}
	return t * t * (3.0 - 2.0*t)
}

// lerp mixes colors
func mixColors(c1, c2 color.Color, t float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	r := float64(r1)*(1.0-t) + float64(r2)*t
	g := float64(g1)*(1.0-t) + float64(g2)*t
	b := float64(b1)*(1.0-t) + float64(b2)*t
	a := float64(a1)*(1.0-t) + float64(a2)*t

	return color.RGBA64{
		R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a),
	}
}

func (r *Road) At(x, y int) color.Color {
	// 1. Transform p to local coordinates (relative to Center, rotated)
	// For rotation: standard 2D rotation matrix.
	// We want to rotate the *world* by -Direction so the road aligns with X axis (for straight).

	dx := float64(x - r.CenterX)
	dy := float64(y - r.CenterY)

	rad := r.Direction * math.Pi / 180.0
	cos := math.Cos(-rad)
	sin := math.Sin(-rad)

	px := dx*cos - dy*sin
	py := dx*sin + dy*cos

	// 2. Calculate Signed Distance (sd) to the road shape.
	// sd < 0 is inside, sd > 0 is outside.
	// We use half-width.
	halfWidth := r.Width / 2.0

	var sd float64
	var u, v float64 // Texture/Marking coordinates (v is lateral distance from center)

	// Valid flags for where markings should be drawn
	drawMarkings := false

	switch r.Shape {
	case RoadStraight:
		// Infinite strip along X axis.
		// Distance is |py| - halfWidth.
		sd = math.Abs(py) - halfWidth
		u = px
		v = py
		drawMarkings = true

	case RoadIntersection:
		// Union of Horizontal (along X) and Vertical (along Y).
		// sdH = |py| - halfWidth
		// sdV = |px| - halfWidth
		// sd = min(sdH, sdV)
		sdH := math.Abs(py) - halfWidth
		sdV := math.Abs(px) - halfWidth
		sd = math.Min(sdH, sdV)

		// Markings logic for intersection:
		// Only draw markings if we are clearly in one arm and NOT in the central box.
		// Central box is where both |px| < halfWidth AND |py| < halfWidth.
		inCenter := (math.Abs(px) < halfWidth) && (math.Abs(py) < halfWidth)

		if !inCenter {
			if sdH < sdV {
				// Horizontal arm
				u = px
				v = py
				drawMarkings = true
			} else {
				// Vertical arm
				u = py // Swap for markings to run along road
				v = px
				drawMarkings = true
			}
		}

	case RoadTIntersection:
		// Horizontal (Top bar) + Vertical (Stem, usually from bottom up to 0)
		// Let's say Horizontal is along X.
		// Vertical is along Y, for y > 0 (or y < 0 depending on coord sys).
		// In image coords, +y is down. Let's make the stem go "down" (positive y) or "up" (negative y).
		// Let's assume stem is for py > 0 (Bottom T).

		sdH := math.Abs(py) - halfWidth

		// Vertical segment: |px| - halfWidth, but only for py > -halfWidth
		// Actually, let's just do union of Line(y=0) and Ray(x=0, y>0).
		// sdV = max(|px| - halfWidth, -py) ? No.
		// SDF for a generic segment or ray is cleaner.
		// But keeping it simple: Union of H-Road and V-Road(clipped).

		sdV := math.Abs(px) - halfWidth
		// Clip sdV: It's only valid if we are "below" the top edge of the horizontal road?
		// Actually, T-intersection implies they merge.
		// sd = min(sdH, sdV intersected with y > 0 area)
		// A simple way: Union of Box(Horizontal) and Box(Vertical, starts at 0).
		// Vertical box: x in [-w/2, w/2], y in [-w/2, infinity].

		// Let's refine:
		// Horizontal Strip: |py| <= w/2
		// Vertical Strip: |px| <= w/2 AND py >= 0

		// Combined SDF:
		// d1 = |py| - w/2
		// d2 = max(|px| - w/2, -py) (Ray starting at 0 going +y) -> No, this is for finite segment.
		// Let's just use:
		// d2 = |px| - w/2.
		// But we need to cut d2 at py = -something?
		// Effectively, if py < -halfWidth, we are far from the vertical road stem.
		// So sd = min(sdH, max(sdV, - (py + halfWidth))) ?
		// If py is very negative (top of image), -(py+halfWidth) is large positive -> sdV is masked.

		// Let's simplify:
		// Just min(sdH, sdV) is a full cross.
		// If we want T (removing top arm of vertical):
		// If py < -halfWidth, distance is to the horizontal strip.

		if py < -halfWidth {
			sd = sdH
		} else {
			sd = math.Min(sdH, sdV)
		}

		// Markings
		inCenter := (math.Abs(px) < halfWidth) && (math.Abs(py) < halfWidth)
		if !inCenter {
			if math.Abs(py) <= halfWidth {
				// Horizontal arm
				u = px
				v = py
				drawMarkings = true
			} else if math.Abs(px) <= halfWidth && py > 0 {
				// Vertical arm (stem)
				u = py
				v = px
				drawMarkings = true
			}
		}

	case RoadCurved:
		// Arc.
		// Center of arc is at (0, CurveRadius) in transformed space?
		// Let's put the pivot at (0,0).
		// Road is an annulus segment.
		// Distance to circle of radius R.
		// p is (px, py).
		// We want the road to curve "away".
		// Let's assume the road starts at (0,0) going +X, and curves towards +Y.
		// Then the center of curvature is at (0, R).
		// d = | length(p - (0, R)) - R | - w/2.

		cx, cy := 0.0, r.CurveRadius
		// But wait, if we rotate by Direction, (0,0) is our anchor.
		// If we want the road to pass through (0,0), then yes, distance from center (0, R) is R.

		distToCenter := math.Sqrt((px-cx)*(px-cx) + (py-cy)*(py-cy))
		sd = math.Abs(distToCenter - r.CurveRadius) - halfWidth

		// Angle check for segment?
		// Vector from center to p: (px, py-R).
		// Angle: atan2(py-R, px).
		// At (0,0), vector is (0, -R). Angle is -90 (or 270).
		// We want to limit the arc length.
		// Let's assume full circle if Angle is large, or clipped.
		// For simplicity in this demo, let's do a full ring or half ring.
		// To clip properly, we'd need sdBox logic on the angle.
		// Let's just use the infinite ring for "Curved" to ensure it looks good and "takes up space".
		// Or clamp it.
		// Let's use the Ring.

		// u, v for markings.
		// v = distToCenter - Radius
		// u = Angle * Radius
		angle := math.Atan2(py-cy, px-cx)
		// Normalize angle so u is continuous-ish?
		u = angle * r.CurveRadius
		v = distToCenter - r.CurveRadius
		drawMarkings = true
	}

	// 3. Rendering
	// Anti-aliasing factor
	aaWidth := 1.0
	alpha := 1.0 - smoothstep(0.0, aaWidth, sd)

	if alpha <= 0.0 {
		// Pure Background
		if r.Background != nil {
			return r.Background.At(x, y)
		}
		// Default Grass
		return color.RGBA{34, 139, 34, 255}
	}

	// We are on road (or edge).
	// Sample Asphalt
	var roadColor color.Color
	if r.Asphalt != nil {
		// Use world coordinates for texture to avoid warping/seams
		roadColor = r.Asphalt.At(x, y)
	} else {
		roadColor = color.RGBA{60, 60, 60, 255}
	}

	// Apply Markings
	if drawMarkings {
		// Lane logic
		laneWidth := r.Width / float64(r.Lanes)
		lineWidth := 3.0

		// Shift v to be 0 at left edge
		// v is centered at 0. range [-w/2, w/2]
		// v' = v + w/2
		vp := v + halfWidth

		// Which lane line?
		// i = round(vp / laneWidth)
		// Location of line i: i * laneWidth
		i := math.Round(vp / laneWidth)
		linePos := i * laneWidth

		distToLine := math.Abs(vp - linePos)

		// Draw line if close
		if distToLine < lineWidth {
			// AA for line
			lineAlpha := 1.0 - smoothstep(lineWidth/2.0 - 0.5, lineWidth/2.0 + 0.5, distToLine)

			if lineAlpha > 0 {
				isEdge := (i == 0 || int(i) == r.Lanes)
				isCenter := (!isEdge && r.Lanes%2 == 0 && int(i) == r.Lanes/2)

				shouldDraw := false

				if isEdge {
					shouldDraw = true
				} else if isCenter {
					// Double line? Or simple solid/dashed.
					// Let's do dashed for center usually, or double solid.
					// Let's do Dashed for center.
					dashLen := 20.0
					if math.Mod(math.Abs(u), dashLen*2) < dashLen {
						shouldDraw = true
					}
				} else {
					// Other dividers: Dashed
					dashLen := 20.0
					if math.Mod(math.Abs(u), dashLen*2) < dashLen {
						shouldDraw = true
					}
				}

				if shouldDraw {
					roadColor = mixColors(roadColor, r.MarkingColor, lineAlpha)
				}
			}
		}
	}

	// Mix Road with Background at edges
	bg := color.RGBA{34, 139, 34, 255}
	if r.Background != nil {
		c := r.Background.At(x, y)
		r1, g1, b1, a1 := c.RGBA()
		bg = color.RGBA{uint8(r1 >> 8), uint8(g1 >> 8), uint8(b1 >> 8), uint8(a1 >> 8)}
	}

	return mixColors(bg, roadColor, alpha)
}

func (r *Road) Bounds() image.Rectangle {
	return r.bounds
}

func (r *Road) ColorModel() color.Model {
	return color.RGBAModel
}
