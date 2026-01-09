package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure ChunkyBands implements the image.Image interface.
var _ image.Image = (*ChunkyBands)(nil)

// ChunkyBands composes chunky pixel bands at a configurable angle.
type ChunkyBands struct {
	Null
	Angle
	BlockSize
	Palette
}

func (p *ChunkyBands) At(x, y int) color.Color {
	block := p.BlockSize.BlockSize
	if block <= 0 {
		block = 8
	}

	// Default palette if none provided.
	pal := p.Palette.Palette
	if len(pal) == 0 {
		pal = []color.Color{
			color.RGBA{10, 20, 40, 255},
			color.RGBA{80, 30, 120, 255},
			color.RGBA{200, 120, 40, 255},
			color.RGBA{240, 210, 120, 255},
		}
	}

	// Quantize coordinates to chunky blocks.
	qx := quantize(x, block)
	qy := quantize(y, block)

	// Project coordinates along the band direction.
	rad := p.Angle.Angle * math.Pi / 180.0
	proj := float64(qx)*math.Cos(rad) + float64(qy)*math.Sin(rad)

	band := int(math.Floor(proj / float64(block)))
	idx := posMod(band, len(pal))

	return pal[idx]
}

func quantize(v, step int) int {
	f := math.Floor(float64(v) / float64(step))
	return int(f * float64(step))
}


// NewChunkyBands creates a new ChunkyBands pattern.
func NewChunkyBands(ops ...func(any)) image.Image {
	p := &ChunkyBands{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
	}
	p.BlockSize.BlockSize = 8
	p.Angle.Angle = 0

	for _, op := range ops {
		op(p)
	}
	return p
}

// NewDemoChunkyBands produces a demo variant for readme.md pre-populated values.
func NewDemoChunkyBands(ops ...func(any)) image.Image {
	return NewChunkyBands(ops...)
}
