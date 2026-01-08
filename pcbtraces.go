package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure PCBTraces implements image.Image.
var _ image.Image = (*PCBTraces)(nil)

// PCBTraces renders an orthogonal trace grid with vias and a solder-mask-tinted background.
type PCBTraces struct {
	Null
	LineSize
	Seed

	padDensity float64
	maskTint   color.Color
	copper     color.Color
}

func (p *PCBTraces) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *PCBTraces) Bounds() image.Rectangle {
	return p.Null.bounds
}

func (p *PCBTraces) SetBounds(b image.Rectangle) {
	p.Null.bounds = b
}

// SetPadDensity adjusts how frequently vias appear (0..1 range recommended).
func (p *PCBTraces) SetPadDensity(v float64) {
	p.padDensity = v
}

// SetSolderMaskTint customizes the base solder mask color.
func (p *PCBTraces) SetSolderMaskTint(c color.Color) {
	p.maskTint = c
}

// SetCopperColor customizes the copper color for traces and vias.
func (p *PCBTraces) SetCopperColor(c color.Color) {
	p.copper = c
}

// SetPCBPadDensity creates an option to set pad density for PCBTraces.
func SetPCBPadDensity(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetPadDensity(float64) }); ok {
			h.SetPadDensity(v)
		}
	}
}

// SetSolderMaskTint creates an option to tint the solder mask.
func SetSolderMaskTint(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetSolderMaskTint(color.Color) }); ok {
			h.SetSolderMaskTint(v)
		}
	}
}

// SetCopperColor creates an option to set the copper color used for traces and vias.
func SetCopperColor(v color.Color) func(any) {
	return func(i any) {
		if h, ok := i.(interface{ SetCopperColor(color.Color) }); ok {
			h.SetCopperColor(v)
		}
	}
}

// NewPCBTraces builds the PCB trace pattern with configurable parameters.
func NewPCBTraces(ops ...func(any)) image.Image {
	p := &PCBTraces{
		Null:       Null{bounds: image.Rect(0, 0, 192, 192)},
		LineSize:   LineSize{LineSize: 3},
		Seed:       Seed{Seed: 1337},
		padDensity: 0.14,
		maskTint:   color.RGBA{14, 82, 44, 255},
		copper:     color.RGBA{200, 140, 60, 255},
	}

	for _, op := range ops {
		op(p)
	}

	if p.LineSize.LineSize <= 0 {
		p.LineSize.LineSize = 1
	}
	if p.padDensity < 0 {
		p.padDensity = 0
	}
	if p.padDensity > 1 {
		p.padDensity = 1
	}

	return p
}

type pcbCell struct {
	north, south, east, west bool
	via                      bool
}

func (p *PCBTraces) cellProfile(cx, cy int) pcbCell {
	h := StableHash(cx, cy, uint64(p.Seed.Seed))
	roll := h % 100

	// Encourage orthogonal runs with occasional corners and tees.
	switch {
	case roll < 15:
		return pcbCell{north: true, south: true}
	case roll < 30:
		return pcbCell{east: true, west: true}
	case roll < 45:
		return pcbCell{north: true, east: true}
	case roll < 60:
		return pcbCell{north: true, west: true}
	case roll < 75:
		return pcbCell{south: true, east: true}
	case roll < 90:
		return pcbCell{south: true, west: true}
	default:
		// Occasional cross/tee adds junction interest.
		return pcbCell{north: true, south: true, east: true}
	}
}

func (p *PCBTraces) At(x, y int) color.Color {
	bounds := p.Bounds()
	if !image.Pt(x, y).In(bounds) {
		return color.RGBA{}
	}

	cellSize := int(math.Max(float64(p.LineSize.LineSize*4), 10))
	cx := (x - bounds.Min.X) / cellSize
	cy := (y - bounds.Min.Y) / cellSize
	cell := p.cellProfile(cx, cy)

	// Via placement driven by padDensity and a secondary hash.
	viaHash := StableHash(cx, cy, uint64(p.Seed.Seed)^0x9e3779b97f4a7c15)
	cell.via = float64(viaHash%1000)/1000.0 < p.padDensity

	// Local coordinates centered in cell.
	lx := float64((x-bounds.Min.X)%cellSize) + 0.5
	ly := float64((y-bounds.Min.Y)%cellSize) + 0.5
	half := float64(cellSize) / 2
	ux := lx - half
	uy := ly - half

	width := math.Max(1, float64(p.LineSize.LineSize))
	run := half - width

	nearTrace := math.MaxFloat64
	checkSegment := func(x1, y1, x2, y2 float64) {
		segLen2 := (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
		if segLen2 == 0 {
			d := math.Hypot(ux-x1, uy-y1)
			if d < nearTrace {
				nearTrace = d
			}
			return
		}
		t := ((ux-x1)*(x2-x1) + (uy-y1)*(y2-y1)) / segLen2
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		px := x1 + t*(x2-x1)
		py := y1 + t*(y2-y1)
		d := math.Hypot(ux-px, uy-py)
		if d < nearTrace {
			nearTrace = d
		}
	}

	if cell.north {
		checkSegment(0, 0, 0, -run)
	}
	if cell.south {
		checkSegment(0, 0, 0, run)
	}
	if cell.west {
		checkSegment(0, 0, -run, 0)
	}
	if cell.east {
		checkSegment(0, 0, run, 0)
	}

	viaDist := math.Hypot(ux, uy)

	copperCore := nearTrace < width/2 || (cell.via && viaDist < width*1.2)
	copperHalo := nearTrace < width*1.5 || (cell.via && viaDist < width*1.6)

	mask := jitterColor(toRGBA(p.maskTint), x, y, p.Seed.Seed)
	if copperHalo {
		mask = mixRGBA(mask, toRGBA(color.RGBA{10, 40, 20, 255}), 0.35)
	}

	if copperCore {
		baseCopper := toRGBA(p.copper)
		shiny := float64(StableHash(x, y, uint64(p.Seed.Seed)^0x51edc0de)) / float64(^uint64(0))
		highlight := float64(baseCopper.A) / 255 * 0.15
		copper := color.RGBA{
			R: clampChannel(float64(baseCopper.R) + highlight*255*(shiny-0.5)),
			G: clampChannel(float64(baseCopper.G) + highlight*255*(shiny-0.5)),
			B: clampChannel(float64(baseCopper.B) + highlight*255*(shiny-0.5)),
			A: 255,
		}
		return copper
	}

	return mask
}

func toRGBA(c color.Color) color.RGBA {
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

func mixRGBA(a, b color.RGBA, t float64) color.RGBA {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	inv := 1 - t
	return color.RGBA{
		R: clampChannel(float64(a.R)*inv + float64(b.R)*t),
		G: clampChannel(float64(a.G)*inv + float64(b.G)*t),
		B: clampChannel(float64(a.B)*inv + float64(b.B)*t),
		A: clampChannel(float64(a.A)*inv + float64(b.A)*t),
	}
}
