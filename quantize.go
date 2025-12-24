package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Quantize implements the image.Image interface.
var _ image.Image = (*Quantize)(nil)

// Quantize reduces the number of colors in an image by quantizing each channel to a specified number of levels.
type Quantize struct {
	Null
	img    image.Image
	levels int
}

func (q *Quantize) ColorModel() color.Model {
	if q.img != nil {
		return q.img.ColorModel()
	}
	return color.RGBAModel
}

func (q *Quantize) Bounds() image.Rectangle {
	return q.bounds
}

func (q *Quantize) At(x, y int) color.Color {
	if q.img == nil {
		return color.RGBA{}
	}
	c := q.img.At(x, y)
	return quantizeColor(c, q.levels)
}

func quantizeColor(c color.Color, levels int) color.Color {
	if levels < 2 {
		return c
	}

	// Convert to NRGBA64 to handle non-premultiplied alpha.
	// We want to quantize the color channels independently of alpha, then premultiply.
	nrgba := color.NRGBA64Model.Convert(c).(color.NRGBA64)

	rq := quantizeChannel(nrgba.R, levels)
	gq := quantizeChannel(nrgba.G, levels)
	bq := quantizeChannel(nrgba.B, levels)
	// Preserve alpha as is
	aq := nrgba.A

	// Create NRGBA64 with quantized values
	qNRGBA := color.NRGBA64{
		R: uint16(rq),
		G: uint16(gq),
		B: uint16(bq),
		A: uint16(aq),
	}

	// Convert back to RGBA64 (premultiplied)
	return color.RGBA64Model.Convert(qNRGBA)
}

func quantizeChannel(v uint16, levels int) uint32 {
	// v is 0..65535
	f := float64(v) / 65535.0
	f = f * float64(levels - 1)
	f = math.Round(f)
	f = f / float64(levels - 1)
	return uint32(f * 65535.0)
}

// NewQuantize creates a new Quantize pattern.
// levels is the number of levels per channel. e.g. 2 means 2 levels (0 and 255).
func NewQuantize(img image.Image, levels int, ops ...func(any)) image.Image {
	if levels < 2 {
		levels = 2
	}
	b := image.Rect(0, 0, 100, 100)
	if img != nil {
		b = img.Bounds()
	}
	p := &Quantize{
		img:    img,
		levels: levels,
		Null: Null{
			bounds: b,
		},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}
