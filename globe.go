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
	FillImageSource
	LatitudeLines
	LongitudeLines
	Angle // Rotation around Y axis in degrees
	Tilt  // Rotation around X axis in degrees
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

	// Outside the sphere
	if distSq > r*r {
		if p.SpaceColor.SpaceColor != nil {
			return p.SpaceColor.SpaceColor
		}
		return color.RGBA{} // Transparent
	}

	// Calculate front-facing point
	z := math.Sqrt(r*r - distSq)

	c, hitLine := p.sampleSpherePoint(dx, dy, z, r)
	if hitLine {
		return c
	}

	// If we have a fill image, map it
	if p.FillImageSource.FillImageSource != nil {
		// Calculate lat/long for UV mapping
		// We need consistent UVs regardless of rotation?
		// Usually UVs are attached to the sphere surface.
		// So we use the rotated coordinates? Yes.

		// Re-calculate rotation to get surface coordinates
		// (Same as in sampleSpherePoint but we need the coordinates)
		_, _, _, phi, lambda := p.getSphereCoords(dx, dy, z, r)

		// Map phi [-pi/2, pi/2] -> v [0, 1]
		// Map lambda [-pi, pi] -> u [0, 1]

		u := (lambda + math.Pi) / (2 * math.Pi)
		v := (phi + math.Pi/2) / math.Pi

		// Sample image
		src := p.FillImageSource.FillImageSource
		sb := src.Bounds()
		sx := int(float64(sb.Min.X) + u*float64(sb.Dx()))
		sy := int(float64(sb.Min.Y) + (1.0-v)*float64(sb.Dy())) // Flip V? usually maps top-down

		// Handle wrapping/clamping?
		// float to int truncates.
		if sx >= sb.Max.X { sx = sb.Max.X - 1 }
		if sy >= sb.Max.Y { sy = sb.Max.Y - 1 }

		return src.At(sx, sy)
	}

	if p.FillColor.FillColor != nil {
		return p.FillColor.FillColor
	}

	// Transparent sphere body (Wireframe mode).
	// Check back face.
	// Back face z is -z (relative to sphere center 0,0,0)
	zBack := -z
	cBack, hitLineBack := p.sampleSpherePoint(dx, dy, zBack, r)
	if hitLineBack {
		return cBack
	}

	// Transparent
	return color.RGBA{}
}

// getSphereCoords rotates the point (dx,dy,dz) and returns rotated (rx, ry, rz) and spherical (phi, lambda).
func (p *Globe) getSphereCoords(dx, dy, dz, r float64) (rx, ry, rz, phi, lambda float64) {
	// Normal vector
	nx := dx / r
	ny := dy / r
	nz := dz / r

	// 1. Tilt (Rotation around X axis)
	tiltRad := p.Tilt.Tilt * math.Pi / 180.0
	cosT := math.Cos(tiltRad)
	sinT := math.Sin(tiltRad)

	// y' = y cos - z sin
	// z' = y sin + z cos
	ny2 := ny*cosT - nz*sinT
	nz2 := ny*sinT + nz*cosT
	nx2 := nx

	// 2. Spin (Rotation around Y axis)
	spinRad := p.Angle.Angle * math.Pi / 180.0
	cosS := math.Cos(spinRad)
	sinS := math.Sin(spinRad)

	// x'' = x' cos - z' sin
	// z'' = x' sin + z' cos
	rx = nx2*cosS - nz2*sinS
	rz = nx2*sinS + nz2*cosS
	ry = ny2

	// Spherical coordinates
	// Latitude phi: arcsin(ry)
	// Longitude lambda: atan2(rz, rx)
	// Clamp ry to [-1, 1] to avoid NaN from float precision errors
	if ry > 1.0 {
		ry = 1.0
	} else if ry < -1.0 {
		ry = -1.0
	}
	phi = math.Asin(ry)
	lambda = math.Atan2(rz, rx)
	return
}

// sampleSpherePoint checks if point (dx, dy, dz) on sphere surface hits a grid line.
func (p *Globe) sampleSpherePoint(dx, dy, dz, r float64) (color.Color, bool) {
	_, _, _, phi, lambda := p.getSphereCoords(dx, dy, dz, r)

	lineSizeVal := float64(p.LineSize.LineSize)
	if lineSizeVal <= 0 {
		return nil, false
	}

	halfLineWidthRad := (lineSizeVal / 2) / r

	// Check Latitude Lines
	if p.LatitudeLines.LatitudeLines > 0 {
		nLat := p.LatitudeLines.LatitudeLines
		dPhi := math.Pi / float64(nLat+1)

		shiftedPhi := phi + math.Pi/2
		k := math.Round(shiftedPhi / dPhi)

		if k > 0 && k < float64(nLat+1) {
			targetPhi := -math.Pi/2 + k*dPhi
			if math.Abs(phi - targetPhi) < halfLineWidthRad {
				return p.LineColor.LineColor, true
			}
		}
	}

	// Check Longitude Lines
	if p.LongitudeLines.LongitudeLines > 0 {
		nLong := p.LongitudeLines.LongitudeLines
		dLambda := 2 * math.Pi / float64(nLong)

		k := math.Round(lambda / dLambda)
		targetLambda := k * dLambda

		diff := math.Abs(lambda - targetLambda)
		if diff > math.Pi {
			diff = 2*math.Pi - diff
		}

		if diff*math.Cos(phi) < halfLineWidthRad {
			return p.LineColor.LineColor, true
		}
	}

	return nil, false
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
