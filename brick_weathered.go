package pattern

import (
	"image"
	"image/color"
	"math"
)

// ChippedBrick renders a brick wall with hue jitter per brick, chipped edges, and uneven mortar shading.
type ChippedBrick struct {
	Null
	BrickWidth, BrickHeight int
	MortarSize              int
	ChipIntensity           float64
	MortarDepth             float64
	HueJitter               float64
	Offset                  float64
	Seed
}

var _ image.Image = (*ChippedBrick)(nil)

func (b *ChippedBrick) At(x, y int) color.Color {
	width := b.BrickWidth
	if width <= 0 {
		width = 48
	}
	height := b.BrickHeight
	if height <= 0 {
		height = 22
	}
	mortar := b.MortarSize
	if mortar < 0 {
		mortar = 3
	}

	chip := clamp01(b.ChipIntensity)
	depth := clamp01(b.MortarDepth)
	hueJitter := b.HueJitter
	if hueJitter <= 0 {
		hueJitter = 0.15
	}

	cellW := width + mortar
	cellH := height + mortar

	row := int(math.Floor(float64(y) / float64(cellH)))
	localY := (y%cellH + cellH) % cellH

	xOffset := 0.0
	offset := b.Offset
	if offset == 0 {
		offset = 0.5
	}
	if row%2 != 0 {
		xOffset = offset * float64(cellW)
	}

	effX := float64(x) - xOffset
	col := int(math.Floor(effX / float64(cellW)))
	floorEffX := int(math.Floor(effX))
	localX := (floorEffX%cellW + cellW) % cellW

	baseMargin := float64(mortar) / 2
	jitterScale := chip * (baseMargin + 0.5)

	lx := float64(localX)
	ly := float64(localY)

	leftMargin := clampFloatRange(baseMargin+b.edgeJitter(col, row, localY, 0x11, jitterScale), 0, float64(mortar))
	rightMargin := clampFloatRange(baseMargin+b.edgeJitter(col, row, localY, 0x23, jitterScale), 0, float64(mortar))
	topMargin := clampFloatRange(baseMargin+b.edgeJitter(col, row, localX, 0x31, jitterScale), 0, float64(mortar))
	bottomMargin := clampFloatRange(baseMargin+b.edgeJitter(col, row, localX, 0x47, jitterScale), 0, float64(mortar))

	brickMinX := leftMargin
	brickMaxX := float64(cellW) - rightMargin
	brickMinY := topMargin
	brickMaxY := float64(cellH) - bottomMargin

	inBrick := lx >= brickMinX && lx < brickMaxX && ly >= brickMinY && ly < brickMaxY

	if !inBrick {
		return mortarColor(lx, ly, brickMinX, brickMaxX, brickMinY, brickMaxY, depth, uint64(b.Seed.Seed))
	}

	brickHash := StableHash(col, row, uint64(b.Seed.Seed))
	tint := (float64(brickHash&0xffff) / 65535.0) - 0.5
	mix := clamp01(0.5 + tint*hueJitter)

	baseA := color.RGBA{170, 70, 55, 255}
	baseB := color.RGBA{200, 95, 70, 255}

	pixelNoise := ((float64(StableHash(x, y, uint64(b.Seed.Seed)^0x9e3779b97f4a7c15)&0xffff) / 65535.0) - 0.5) * 6
	edgeSoftness := math.Min(
		math.Min(lx-brickMinX, brickMaxX-lx),
		math.Min(ly-brickMinY, brickMaxY-ly),
	)
	edgeSoftness = math.Max(0, edgeSoftness)
	chipShade := 0.0
	if edgeSoftness < 1.5 {
		chipShade = -chip * 25 * (1.5 - edgeSoftness) / 1.5
	}

	r := lerpBrick(float64(baseA.R), float64(baseB.R), mix) + pixelNoise + chipShade
	g := lerpBrick(float64(baseA.G), float64(baseB.G), mix) + pixelNoise*0.7 + chipShade*0.7
	bc := lerpBrick(float64(baseA.B), float64(baseB.B), mix) + pixelNoise*0.5 + chipShade*0.6

	return color.RGBA{
		R: clampColor(r),
		G: clampColor(g),
		B: clampColor(bc),
		A: 255,
	}
}

func (b *ChippedBrick) ColorModel() color.Model {
	return color.RGBAModel
}

// NewChippedBrick creates a weathered brick wall.
func NewChippedBrick(ops ...func(any)) image.Image {
	p := &ChippedBrick{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		BrickWidth:    48,
		BrickHeight:   22,
		MortarSize:    3,
		ChipIntensity: 0.35,
		MortarDepth:   0.7,
		HueJitter:     0.15,
		Offset:        0.5,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// Options and setters

type ChipIntensity struct{ Amount float64 }

func (c *ChipIntensity) SetChipIntensity(v float64) { c.Amount = v }

type MortarDepth struct{ Depth float64 }

func (m *MortarDepth) SetMortarDepth(v float64) { m.Depth = v }

func SetChipIntensity(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetChipIntensity(float64) }); ok {
			h.SetChipIntensity(v)
		}
	}
}

func SetMortarDepth(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetMortarDepth(float64) }); ok {
			h.SetMortarDepth(v)
		}
	}
}

type hasHueJitter interface{ SetHueJitter(float64) }

func SetHueJitter(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasHueJitter); ok {
			h.SetHueJitter(v)
		}
	}
}

func (b *ChippedBrick) SetBounds(r image.Rectangle) { b.bounds = r }
func (b *ChippedBrick) SetBrickSize(w, h int)       { b.BrickWidth, b.BrickHeight = w, h }
func (b *ChippedBrick) SetMortarSize(v int)         { b.MortarSize = v }
func (b *ChippedBrick) SetBrickOffset(v float64)    { b.Offset = v }
func (b *ChippedBrick) SetChipIntensity(v float64)  { b.ChipIntensity = v }
func (b *ChippedBrick) SetMortarDepth(v float64)    { b.MortarDepth = v }
func (b *ChippedBrick) SetHueJitter(v float64)      { b.HueJitter = v }
func (b *ChippedBrick) SetSeedUint64(v uint64)      { b.Seed.Seed = int64(v) }
func (b *ChippedBrick) SetSeed(v int64)             { b.Seed.Seed = v }

func (b *ChippedBrick) edgeJitter(col, row, coord int, salt uint64, scale float64) float64 {
	if scale == 0 {
		return 0
	}
	h := StableHash(col*131+coord, row*197+coord, uint64(b.Seed.Seed)^salt)
	n := float64(h&0xffff) / 65535.0
	return (n - 0.5) * scale
}

func mortarColor(x, y, brickMinX, brickMaxX, brickMinY, brickMaxY, depth float64, seed uint64) color.RGBA {
	base := 190.0
	noise := ((float64(StableHash(int(x), int(y), seed^0xabcdef)&0xffff) / 65535.0) - 0.5) * 12

	distX := math.Min(math.Abs(x-brickMinX), math.Abs(x-brickMaxX))
	distY := math.Min(math.Abs(y-brickMinY), math.Abs(y-brickMaxY))
	edgeDist := math.Min(distX, distY)
	edgeFalloff := clamp01(edgeDist / 3.0)
	shade := -depth * 28 * (1 - edgeFalloff)

	v := base + noise + shade - depth*6
	return color.RGBA{
		R: clampColor(v),
		G: clampColor(v - 2),
		B: clampColor(v - 4),
		A: 255,
	}
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func clampColor(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

func clampFloatRange(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func lerpBrick(a, b, t float64) float64 {
	return a + (b-a)*t
}
