package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure WorleyTiles implements image.Image.
var _ image.Image = (*WorleyTiles)(nil)

// WorleyTiles generates tiled Voronoi/Worley stones with mortar.
// Cells are rounded with a smooth transition near the gaps and each cell
// receives a small palette jitter for natural variation.
type WorleyTiles struct {
	Null
	Seed
	Frequency
	StoneSize     float64
	GapWidth      float64
	PaletteSpread float64
	Jitter        float64
	Palette       []color.RGBA
}

// NewWorleyTiles creates a new WorleyTiles pattern.
func NewWorleyTiles(ops ...func(any)) image.Image {
	w := &WorleyTiles{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		StoneSize:     72,
		GapWidth:      0.08,
		PaletteSpread: 0.12,
		Jitter:        0.85,
		Palette: []color.RGBA{
			{134, 125, 116, 255},
			{112, 104, 98, 255},
			{148, 138, 126, 255},
			{120, 112, 104, 255},
		},
	}
	for _, op := range ops {
		op(w)
	}
	return w
}

// SetTileStoneSize sets the desired average stone size (in pixels).
func SetTileStoneSize(v float64) func(any) {
	return func(i any) {
		if w, ok := i.(*WorleyTiles); ok {
			w.StoneSize = v
		}
	}
}

// SetTileGapWidth configures the mortar gap width in Worley distance space (0-1).
func SetTileGapWidth(v float64) func(any) {
	return func(i any) {
		if w, ok := i.(*WorleyTiles); ok {
			w.GapWidth = v
		}
	}
}

// SetTilePaletteSpread configures how much each stone jitters its color palette (0-1).
func SetTilePaletteSpread(v float64) func(any) {
	return func(i any) {
		if w, ok := i.(*WorleyTiles); ok {
			w.PaletteSpread = v
		}
	}
}

// SetTilePalette overrides the default palette used for stones.
func SetTilePalette(colors ...color.RGBA) func(any) {
	return func(i any) {
		if w, ok := i.(*WorleyTiles); ok && len(colors) > 0 {
			w.Palette = colors
		}
	}
}

// SetTileJitter sets the Worley point jitter inside each cell (0-1).
func SetTileJitter(j float64) func(any) {
	return func(i any) {
		if w, ok := i.(*WorleyTiles); ok {
			w.Jitter = j
		}
	}
}

func (w *WorleyTiles) palette() []color.RGBA {
	if len(w.Palette) == 0 {
		return []color.RGBA{
			{134, 125, 116, 255},
			{112, 104, 98, 255},
			{148, 138, 126, 255},
			{120, 112, 104, 255},
		}
	}
	return w.Palette
}

func (w *WorleyTiles) ColorModel() color.Model {
	return color.RGBAModel
}

func (w *WorleyTiles) Bounds() image.Rectangle {
	return w.bounds
}

func (w *WorleyTiles) At(x, y int) color.Color {
	freq := w.Frequency.Frequency
	if w.StoneSize > 0 {
		freq = 1.0 / w.StoneSize
	}
	if freq <= 0 {
		freq = 1.0 / 72.0
	}
	jitter := w.Jitter
	if jitter <= 0 {
		jitter = 1.0
	}

	nx, ny := float64(x)*freq, float64(y)*freq
	ix, iy := math.Floor(nx), math.Floor(ny)
	fx, fy := nx-ix, ny-iy

	minDist := math.MaxFloat64
	secondMinDist := math.MaxFloat64
	closestHash := uint64(0)
	closestCellX, closestCellY := 0, 0

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			cellX := int(ix) + dx
			cellY := int(iy) + dy

			h := StableHash(cellX, cellY, uint64(w.Seed.Seed))
			rX := float64(h&0xFFFF) / 65535.0
			rY := float64((h>>16)&0xFFFF) / 65535.0

			pointX := float64(dx) + rX*jitter
			pointY := float64(dy) + rY*jitter

			diffX := pointX - fx
			diffY := pointY - fy
			dist := math.Sqrt(diffX*diffX + diffY*diffY)

			if dist < minDist {
				secondMinDist = minDist
				minDist = dist
				closestHash = h
				closestCellX = cellX
				closestCellY = cellY
			} else if dist < secondMinDist {
				secondMinDist = dist
			}
		}
	}

	const maxCellDistance = math.Sqrt2
	minNorm := clamp01(minDist / maxCellDistance)
	borderDistance := clamp01((secondMinDist - minDist) / maxCellDistance)

	mortarWidth := w.GapWidth
	if mortarWidth <= 0 {
		mortarWidth = 0.08
	}
	stoneMask := smoothStep(mortarWidth, mortarWidth*1.7, borderDistance)

	palette := w.palette()
	baseColor := palette[int(closestHash%uint64(len(palette)))]
	jittered := tileJitterColor(baseColor, w.PaletteSpread, closestHash, closestCellX, closestCellY)

	edgeLight := 0.75 + 0.25*smoothStep(mortarWidth*0.5, 1.0, borderDistance)
	centerLight := 0.8 + 0.2*(1.0-minNorm)
	brightness := edgeLight * centerLight
	shaded := scaleColor(jittered, brightness)

	mortar := color.RGBA{30, 28, 32, 255}
	final := lerpRGBA(mortar, shaded, stoneMask)
	return final
}

func tileJitterColor(c color.RGBA, spread float64, h uint64, cellX, cellY int) color.RGBA {
	if spread <= 0 {
		return c
	}
	// Use multiple hashes to decorrelate channels and edges.
	seedR := StableHash(cellX+17, cellY-9, h)
	seedG := StableHash(cellX-11, cellY+23, h>>1)
	seedB := StableHash(cellX+5, cellY+5, h<<1)

	r := tileJitterChannel(float64(c.R), spread, seedR)
	g := tileJitterChannel(float64(c.G), spread, seedG)
	b := tileJitterChannel(float64(c.B), spread, seedB)

	return color.RGBA{uint8(r), uint8(g), uint8(b), c.A}
}

func tileJitterChannel(base float64, spread float64, h uint64) float64 {
	rand := float64(h&0xFFFF) / 65535.0
	scale := 1 + (rand-0.5)*2*spread
	return clamp255(base * scale)
}

func scaleColor(c color.RGBA, scale float64) color.RGBA {
	return color.RGBA{
		R: uint8(clamp255(float64(c.R) * scale)),
		G: uint8(clamp255(float64(c.G) * scale)),
		B: uint8(clamp255(float64(c.B) * scale)),
		A: c.A,
	}
}

func smoothStep(edge0, edge1, x float64) float64 {
	if edge1 == edge0 {
		if x < edge0 {
			return 0
		}
		return 1
	}
	t := clamp01((x - edge0) / (edge1 - edge0))
	return t * t * (3 - 2*t)
}
