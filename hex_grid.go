package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure HexGrid implements the image.Image interface.
var _ image.Image = (*HexGrid)(nil)

// HexGrid renders an axial-coordinate hexagonal grid with alternating colors
// and a soft inner shadow near the cell edges.
type HexGrid struct {
	Null
	Radius
	Palette    color.Palette
	BevelDepth float64
}

func (h *HexGrid) ColorModel() color.Model {
	return color.RGBAModel
}

func (h *HexGrid) Bounds() image.Rectangle {
	return h.bounds
}

func (h *HexGrid) SetHexPalette(pal color.Palette) {
	h.Palette = pal
}

func (h *HexGrid) SetHexBevelDepth(depth float64) {
	h.BevelDepth = depth
}

func (h *HexGrid) At(x, y int) color.Color {
	radius := h.Radius.Radius
	if radius <= 0 {
		radius = 24
	}
	size := float64(radius)

	bevel := h.BevelDepth
	if bevel <= 0 {
		bevel = size * 0.35
	}

	palette := h.Palette
	if len(palette) == 0 {
		palette = color.Palette{
			color.NRGBA{R: 36, G: 66, B: 86, A: 255},
			color.NRGBA{R: 147, G: 197, B: 253, A: 255},
		}
	}

	// Center of bounds for symmetric layout.
	b := h.Bounds()
	cx := float64(b.Min.X+b.Max.X) / 2.0
	cy := float64(b.Min.Y+b.Max.Y) / 2.0

	// Relative coordinates from center.
	rx := float64(x) - cx
	ry := float64(y) - cy

	// Convert pixel coordinate to axial coordinate (pointy-top orientation).
	q := (math.Sqrt(3)/3*rx - ry/3.0) / size
	r := (2.0 / 3.0 * ry) / size

	aq, ar := axialRound(q, r)

	// Base color alternating by axial position.
	idx := posMod(aq+ar, len(palette))
	base := color.NRGBAModel.Convert(palette[idx]).(color.NRGBA)

	// Calculate the center of the target hex for shading.
	hx := size * (math.Sqrt(3) * (float64(aq) + float64(ar)/2.0))
	hy := size * (1.5 * float64(ar))
	lx := rx - hx
	ly := ry - hy

	dist := distanceToHexEdge(lx, ly, size)

	shadow := 0.0
	if bevel > 0 && dist < bevel {
		t := 1.0 - dist/bevel
		shadow = t * t
	}

	darken := 1.0 - 0.45*shadow

	return color.NRGBA{
		R: uint8(clampFloatRange(float64(base.R)*darken, 0, 255)),
		G: uint8(clampFloatRange(float64(base.G)*darken, 0, 255)),
		B: uint8(clampFloatRange(float64(base.B)*darken, 0, 255)),
		A: base.A,
	}
}

// NewHexGrid creates a new HexGrid pattern.
func NewHexGrid(ops ...func(any)) image.Image {
	p := &HexGrid{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Radius: Radius{Radius: 24},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// axialRound converts fractional axial coordinates to the nearest integer axial coordinates.
func axialRound(q, r float64) (int, int) {
	x := q
	z := r
	y := -x - z

	rx := math.Round(x)
	ry := math.Round(y)
	rz := math.Round(z)

	xDiff := math.Abs(rx - x)
	yDiff := math.Abs(ry - y)
	zDiff := math.Abs(rz - z)

	if xDiff > yDiff && xDiff > zDiff {
		rx = -ry - rz
	} else {
		rz = -rx - ry
	}

	return int(rx), int(rz)
}

type floatPoint struct {
	X, Y float64
}

func hexCorners(size float64) []floatPoint {
	corners := make([]floatPoint, 6)
	for i := 0; i < 6; i++ {
		angle := math.Pi/6.0 + math.Pi/3.0*float64(i)
		corners[i] = floatPoint{
			X: size * math.Cos(angle),
			Y: size * math.Sin(angle),
		}
	}
	return corners
}

func distanceToHexEdge(x, y, size float64) float64 {
	p := floatPoint{X: x, Y: y}
	corners := hexCorners(size)
	minDist := math.MaxFloat64
	for i := 0; i < 6; i++ {
		a := corners[i]
		b := corners[(i+1)%6]
		if d := pointSegmentDistance(p, a, b); d < minDist {
			minDist = d
		}
	}
	return minDist
}

func pointSegmentDistance(p, a, b floatPoint) float64 {
	vx := b.X - a.X
	vy := b.Y - a.Y
	wx := p.X - a.X
	wy := p.Y - a.Y

	lenSq := vx*vx + vy*vy
	t := 0.0
	if lenSq > 0 {
		t = (wx*vx + wy*vy) / lenSq
	}
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	dx := wx - t*vx
	dy := wy - t*vy
	return math.Hypot(dx, dy)
}

// SetHexPalette creates an option to set the palette used by the HexGrid.
func SetHexPalette(pal color.Palette) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetHexPalette(color.Palette) }); ok {
			h.SetHexPalette(pal)
		}
	}
}

// SetHexBevelDepth creates an option to set the depth of the soft inner shadow.
func SetHexBevelDepth(depth float64) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetHexBevelDepth(float64) }); ok {
			h.SetHexBevelDepth(depth)
		}
	}
}
