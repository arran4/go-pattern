package pattern

import (
	"image"
	"image/color"
	"math"
)

// Globe renders a 3D sphere projected onto 2D.
// It supports configurable latitude and longitude grid lines,
// 3D rotation via Angle (Y-axis) and Tilt (X-axis/Z-axis),
// and texture mapping via FillImageSource with UV coordinates.
type Globe struct {
	Null
	FillImageSource
	LatitudeLines  // Number of lines along latitude
	LongitudeLines // Number of lines along longitude
	Angle          // Rotation around Y axis (in degrees)
	Tilt           // Rotation around X axis (in degrees)
	LineColor
	SpaceColor
}

func (p *Globe) At(x, y int) color.Color {
	b := p.Bounds()
	width := float64(b.Dx())
	height := float64(b.Dy())
	cx := float64(b.Min.X) + width/2
	cy := float64(b.Min.Y) + height/2

	minDim := width
	if height < minDim {
		minDim = height
	}
	radius := minDim / 2

	// Normalized coordinates [-1, 1]
	nx := (float64(x) - cx) / radius
	ny := (float64(y) - cy) / radius

	// Check if outside circle
	distSq := nx*nx + ny*ny
	if distSq > 1 {
		if p.SpaceColor.SpaceColor != nil {
			return p.SpaceColor.SpaceColor
		}
		return color.RGBA{}
	}

	// Calculate Z (point on sphere surface, facing viewer)
	z := math.Sqrt(1 - distSq)

	// 3D Point P = (nx, ny, z)
	// Coordinate system: X right, Y down (screen), Z out (towards viewer)
	// Wait, for standard math, let's map screen Y (down) to 3D Y (up) by inverting.
	// py = -ny.
	// px = nx.
	// pz = z.
	px := nx
	py := -ny
	pz := z

	// Rotations.
	// Tilt (X-axis rotation). Positive tilt moves top away or towards?
	// Let's use standard rotation.
	tiltRad := p.Tilt.Tilt * math.Pi / 180
	cosT := math.Cos(tiltRad)
	sinT := math.Sin(tiltRad)

	// Rotate around X
	// y' = y*cos - z*sin
	// z' = y*sin + z*cos
	py2 := py*cosT - pz*sinT
	pz2 := py*sinT + pz*cosT
	px2 := px

	// Angle (Y-axis rotation).
	angleRad := p.Angle.Angle * math.Pi / 180
	cosA := math.Cos(angleRad)
	sinA := math.Sin(angleRad)

	// Rotate around Y
	// x'' = x*cos + z*sin
	// z'' = -x*sin + z*cos
	px3 := px2*cosA + pz2*sinA
	pz3 := -px2*sinA + pz2*cosA
	py3 := py2

	// Spherical Coords
	// Lat (phi): -pi/2 to pi/2 (South to North)
	// Lon (theta): -pi to pi

	// Clamp Y to [-1, 1]
	if py3 > 1 {
		py3 = 1
	}
	if py3 < -1 {
		py3 = -1
	}

	phi := math.Asin(py3)
	theta := math.Atan2(px3, pz3)

	// Texture Mapping (UV)
	// u = (theta + pi) / (2pi) -> 0 to 1
	// v = 0.5 - phi/pi -> 0 to 1 (North Pole at 0)
	u := (theta + math.Pi) / (2 * math.Pi)
	v := 0.5 - phi/math.Pi

	// If we have a texture, sample it.
	if p.FillImageSource.FillImageSource != nil {
		// Map UV to source bounds
		src := p.FillImageSource.FillImageSource
		sb := src.Bounds()
		sx := int(u * float64(sb.Dx()))
		sy := int(v * float64(sb.Dy()))
		// Add Min
		sx += sb.Min.X
		sy += sb.Min.Y
		// Clamp to bounds (just in case)
		if sx >= sb.Max.X {
			sx = sb.Max.X - 1
		}
		if sy >= sb.Max.Y {
			sy = sb.Max.Y - 1
		}
		return src.At(sx, sy)
	}

	// If no texture, check for Grid Lines.
	latLines := p.LatitudeLines.LatitudeLines
	lonLines := p.LongitudeLines.LongitudeLines
	if latLines > 0 || lonLines > 0 {
		lineColor := p.LineColor.LineColor
		if lineColor == nil {
			lineColor = color.Black
		}

		// Line thickness in angle radians?
		// A constant screen-space thickness is hard without derivatives.
		// Let's use a fixed angular thickness for simplicity, or approximate screen derivative.
		// Use 0.02 radians ~ 1 degree as base thickness?
		thickness := 0.02

		// Latitude Lines
		if latLines > 0 {
			// phi ranges -pi/2 to pi/2.
			// Normalized lat: (phi + pi/2) / pi * Lines
			val := (phi + math.Pi/2) / math.Pi * float64(latLines)
			_, frac := math.Modf(val)
			// Check distance to integer
			if frac < thickness || frac > 1-thickness {
				return lineColor
			}
		}

		// Longitude Lines
		if lonLines > 0 {
			// theta ranges -pi to pi
			// Normalized lon: (theta + pi) / (2*pi) * Lines
			val := (theta + math.Pi) / (2 * math.Pi) * float64(lonLines)
			_, frac := math.Modf(val)
			if frac < thickness/2 || frac > 1-thickness/2 {
				// thickness/2 because circumference is 2pi?
				return lineColor
			}
		}
	}

	// Default fill if no texture and not on line
	return color.White // Default sphere color? Or Transparent?
}

// NewGlobe creates a new Globe pattern.
func NewGlobe(ops ...func(any)) image.Image {
	p := &Globe{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	p.LineColor.LineColor = color.Black

	for _, op := range ops {
		op(p)
	}
	return p
}
