package pattern

import (
	"image"
	"image/color"
)

var _ image.Image = (*Padding)(nil)

type Padding struct {
	img      image.Image
	bounds   image.Rectangle
	bgPattern image.Image
	padding  int
}

func (p *Padding) ColorModel() color.Model {
	return p.img.ColorModel()
}

func (p *Padding) Bounds() image.Rectangle {
	return p.bounds
}

func (p *Padding) At(x, y int) color.Color {
	// Check if we are inside the inner image area
	innerMinX := p.bounds.Min.X + p.padding
	innerMinY := p.bounds.Min.Y + p.padding
	innerMaxX := p.bounds.Max.X - p.padding
	innerMaxY := p.bounds.Max.Y - p.padding

	// If the inner image has its own bounds relative to 0,0 (or whatever), we map it.
	// But p.img might be unbounded or have arbitrary bounds.
	// We want to map the inner rectangle [innerMin, innerMax] to the p.img's coordinate space.
	// But how?
	// If p.img is just "placed" there.
	// Usually "Padding" implies the image is centered or top-left aligned inside the padding?
	// NewPadding(img, padding) -> image is at (padding, padding).
	// So (x, y) corresponds to (x - padding, y - padding) in image space (assuming image starts at 0,0).
	// But if image starts at (MinX, MinY), then (x,y) -> (x - padding + MinX, y - padding + MinY).

	if x >= innerMinX && x < innerMaxX && y >= innerMinY && y < innerMaxY {
		lx := x - innerMinX + p.img.Bounds().Min.X
		ly := y - innerMinY + p.img.Bounds().Min.Y

		pt := image.Point{X: lx, Y: ly}
		// If image is unbounded (Bounds are huge or 0,0,1,1 but effectively infinite), In() might fail if not careful?
		// But for unbounded patterns like Checker, Bounds might be default.
		// If we want to "bound an unbounded image", we should just call At(lx, ly).
		// Checking In() is optimization/safety for bounded images.
		// If we know we are "bounding" it, we might skip In() check or rely on the fact that we are inside the "PaddingBounds".

		if pt.In(p.img.Bounds()) {
			return p.img.At(lx, ly)
		} else {
			// If unbounded logic, maybe we still return At?
			// But standard Image behavior is zero outside bounds.
			// If we want to support Unbounded, the input image should probably report large bounds.
			// Or we assume `At` works.
			// Let's stick to standard behavior: If point is in bounds, return color.
			// If p.img is "Unbounded" (infinite), it should return large bounds.
		}
	}

	// Outside inner area (or transparent part of inner image?), draw background
	if p.bgPattern != nil {
		return p.bgPattern.At(x, y)
	}
	return color.RGBA{}
}

type PaddingOption func(*paddingConfig)

type paddingConfig struct {
	margin int
	bg     image.Image
	bounds image.Rectangle // explicit bounds to enforce
}

func PaddingMargin(m int) PaddingOption {
	return func(c *paddingConfig) {
		c.margin = m
	}
}

func PaddingBackground(bg image.Image) PaddingOption {
	return func(c *paddingConfig) {
		c.bg = bg
	}
}

// PaddingBoundary forces the output size of the padding image.
// This effectively bounds the content if the content is unbounded.
// The content will be placed inside these bounds (minus padding).
func PaddingBoundary(r image.Rectangle) PaddingOption {
	return func(c *paddingConfig) {
		c.bounds = r
	}
}

func NewPadding(img image.Image, opts ...PaddingOption) image.Image {
	cfg := &paddingConfig{
		margin: 0,
		bg:     nil,
	}
	for _, o := range opts {
		o(cfg)
	}

	var finalBounds image.Rectangle

	if !cfg.bounds.Empty() {
		finalBounds = cfg.bounds
	} else {
		// Calculate based on image bounds + margin
		b := img.Bounds()
		// If image is "unbounded" (e.g. 0-sized or huge), this might be weird.
		// But defaulting to image bounds is standard.
		w := b.Dx() + 2*cfg.margin
		h := b.Dy() + 2*cfg.margin
		finalBounds = image.Rect(0, 0, w, h)
	}

	return &Padding{
		img:       img,
		bounds:    finalBounds,
		bgPattern: cfg.bg,
		padding:   cfg.margin,
	}
}
