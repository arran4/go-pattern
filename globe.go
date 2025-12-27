package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Globe implements the image.Image interface.
var _ image.Image = (*Globe)(nil)

// LatitudeLines configures the number of latitude lines.
type LatitudeLines struct {
	LatitudeLines int
}

func (l *LatitudeLines) SetLatitudeLines(v int) {
	l.LatitudeLines = v
}

type hasLatitudeLines interface {
	SetLatitudeLines(int)
}

// SetLatitudeLines creates an option to set the number of latitude lines.
func SetLatitudeLines(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasLatitudeLines); ok {
			h.SetLatitudeLines(v)
		}
	}
}

// LongitudeLines configures the number of longitude lines.
type LongitudeLines struct {
	LongitudeLines int
}

func (l *LongitudeLines) SetLongitudeLines(v int) {
	l.LongitudeLines = v
}

type hasLongitudeLines interface {
	SetLongitudeLines(int)
}

// SetLongitudeLines creates an option to set the number of longitude lines.
func SetLongitudeLines(v int) func(any) {
	return func(i any) {
		if h, ok := i.(hasLongitudeLines); ok {
			h.SetLongitudeLines(v)
		}
	}
}

// Globe represents a sphere pattern with optional grid lines.
type Globe struct {
	Null
	Center
	Radius
	LineColor
	LineSize
	SpaceColor
	FillColor
	LatitudeLines
	LongitudeLines
	Angle // Rotation around Y axis in degrees
}

func (p *Globe) At(x, y int) color.Color {
	b := p.Bounds()

	// Determine Center
	cx := p.CenterX
	cy := p.CenterY
	if cx == 0 && cy == 0 {
		cx = b.Min.X + b.Dx()/2
		cy = b.Min.Y + b.Dy()/2
	}

	// Determine Radius
	r := float64(p.Radius.Radius)
	if r == 0 {
		minDim := b.Dx()
		if b.Dy() < minDim {
			minDim = b.Dy()
		}
		r = float64(minDim) / 2
	}

	// Check against radius
	dx := float64(x - cx)
	dy := float64(y - cy)
	distSq := dx*dx + dy*dy
	if distSq > r*r {
		if p.SpaceColor.SpaceColor != nil {
			return p.SpaceColor.SpaceColor
		}
		return color.RGBA{} // Transparent
	}

	// Z coordinate on sphere
	z := math.Sqrt(r*r - distSq)

	// Normal vector
	nx := dx / r
	ny := dy / r
	nz := z / r

	// Rotate around Y axis
	rad := p.Angle.Angle * math.Pi / 180.0
	cosA := math.Cos(rad)
	sinA := math.Sin(rad)

	// Rotated vector (nx', ny', nz')
	// Rotation: x' = x cos - z sin, z' = x sin + z cos
	rx := nx*cosA - nz*sinA
	ry := ny
	rz := nx*sinA + nz*cosA

	// Convert to spherical coordinates
	// Latitude phi: arcsin(y)
	// Longitude lambda: atan2(z, x)
	// Note: in math package, Asin returns [-pi/2, pi/2].
	// Atan2 returns [-pi, pi].

	phi := math.Asin(ry)
	lambda := math.Atan2(rz, rx)

	lineSizeVal := float64(p.LineSize.LineSize)
	drawLines := lineSizeVal > 0

	if drawLines {
		halfLineWidthRad := (lineSizeVal / 2) / r

		// Check Latitude Lines
		if p.LatitudeLines.LatitudeLines > 0 {
			nLat := p.LatitudeLines.LatitudeLines
			// Spacing. If nLat=1, draw at 0 (Equator).
			// If nLat=2, draw at -pi/6, +pi/6? or -pi/4, +pi/4?
			// Let's divide pi into nLat+1 segments.
			dPhi := math.Pi / float64(nLat+1)

			// We want to find closest k*dPhi to phi + pi/2
			// phi ranges [-pi/2, pi/2].
			// shifted: phi + pi/2 ranges [0, pi].

			shiftedPhi := phi + math.Pi/2
			k := math.Round(shiftedPhi / dPhi)

			// If k=0 or k=nLat+1, these are poles. Usually don't draw point at pole as a line,
			// but meridians meet there anyway.
			if k > 0 && k < float64(nLat+1) {
				targetPhi := -math.Pi/2 + k*dPhi
				if math.Abs(phi - targetPhi) < halfLineWidthRad {
					return p.LineColor.LineColor
				}
			}
		}

		// Check Longitude Lines
		if p.LongitudeLines.LongitudeLines > 0 {
			nLong := p.LongitudeLines.LongitudeLines
			dLambda := 2 * math.Pi / float64(nLong)

			// We want closest multiple of dLambda to lambda.
			// lambda in [-pi, pi].

			k := math.Round(lambda / dLambda)
			targetLambda := k * dLambda

			// Distance on sphere: |lambda - target| * cos(phi)
			diff := math.Abs(lambda - targetLambda)
			// Handle wrap around pi/-pi
			if diff > math.Pi {
				diff = 2*math.Pi - diff
			}

			if diff*math.Cos(phi) < halfLineWidthRad {
				return p.LineColor.LineColor
			}
		}
	}

	if p.FillColor.FillColor != nil {
		return p.FillColor.FillColor
	}

	// If no fill color and we didn't hit a line, return transparent?
	// Or maybe the user expects a solid sphere if lines are not drawn or hit?
	// Circle returns FillColor if set.
	return color.RGBA{}
}

// NewGlobe creates a new Globe pattern.
func NewGlobe(ops ...func(any)) image.Image {
	p := &Globe{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	// Defaults
	p.LineColor.LineColor = color.Black
	p.LineSize.LineSize = 1
	p.LatitudeLines.LatitudeLines = 0
	p.LongitudeLines.LongitudeLines = 0

	for _, op := range ops {
		op(p)
	}
	return p
}
