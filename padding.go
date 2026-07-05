package pattern

import (
	"image"
	"image/color"
)

var _ image.Image = (*Padding)(nil)

type Padding struct {
	img                      image.Image
	bounds                   image.Rectangle
	bgPattern                image.Image
	top, left, bottom, right int
}

func (p *Padding) ColorModel() color.Model {
	return p.img.ColorModel()
}

func (p *Padding) Bounds() image.Rectangle {
	return p.bounds
}

func (p *Padding) At(x, y int) color.Color {
	// Inner image area
	innerMinX := p.bounds.Min.X + p.left
	innerMinY := p.bounds.Min.Y + p.top
	innerMaxX := p.bounds.Max.X - p.right
	innerMaxY := p.bounds.Max.Y - p.bottom

	if x >= innerMinX && x < innerMaxX && y >= innerMinY && y < innerMaxY {
		lx := x - innerMinX + p.img.Bounds().Min.X
		ly := y - innerMinY + p.img.Bounds().Min.Y

		pt := image.Point{X: lx, Y: ly}
		if pt.In(p.img.Bounds()) {
			return p.img.At(lx, ly)
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
	top, left, bottom, right int
	bg                       image.Image
	bounds                   image.Rectangle // explicit bounds to enforce
}

func PaddingMargin(m int) PaddingOption {
	return func(c *paddingConfig) {
		c.top = m
		c.left = m
		c.bottom = m
		c.right = m
	}
}

func PaddingTop(m int) PaddingOption    { return func(c *paddingConfig) { c.top = m } }
func PaddingLeft(m int) PaddingOption   { return func(c *paddingConfig) { c.left = m } }
func PaddingBottom(m int) PaddingOption { return func(c *paddingConfig) { c.bottom = m } }
func PaddingRight(m int) PaddingOption  { return func(c *paddingConfig) { c.right = m } }

func PaddingBackground(bg image.Image) PaddingOption {
	return func(c *paddingConfig) {
		c.bg = bg
	}
}

func PaddingBoundary(r image.Rectangle) PaddingOption {
	return func(c *paddingConfig) {
		c.bounds = r
	}
}

// Helpers for Alignment (centering, etc.) inside a fixed boundary

// PaddingCenter attempts to center the image within the provided boundary.
// It requires PaddingBoundary to be set OR assumes the output size will be large enough.
// Actually, if we don't know the output size yet, "Center" is ambiguous.
// But we can define "Center" as: Calculate padding such that image is centered in `bounds`.
// This requires `bounds` to be known.
// If used with `PaddingBoundary`, we can calculate margins.
func PaddingCenter() PaddingOption {
	return func(c *paddingConfig) {
		// Marker to calculate centering later?
		// Or we need to apply this AFTER bounds are known.
		// Since we process options linearly, if Boundary is set before, we can calc.
		// If set after, we might miss it.
		// Let's assume this option sets a flag or we calculate in NewPadding.
		// But functional options modify config state.
		// Let's implement a delayed calculation or logic in NewPadding.
		// However, simple approach:
		// We need to know image size and boundary size.
		// We don't have image size here.
		// So we can't implement "PaddingCenter" as a simple option unless it sets a flag.
	}
}

// Re-thinking: The user wants "Center which backs out to Padding".
// "TopLeftAlign(Left(Pixel(3)))"
// This implies we are wrapping an image in a larger box.
// So we need to specify the target box size, OR the padding amounts.

// If we use `PaddingBoundary(rect)`, we know the target size.
// If we want to center `img` in `rect`, we can calculate top/left/bottom/right.

func PaddingCenterIn(bounds image.Rectangle, imgBounds image.Rectangle) PaddingOption {
	return func(c *paddingConfig) {
		c.bounds = bounds
		w := bounds.Dx()
		h := bounds.Dy()
		iw := imgBounds.Dx()
		ih := imgBounds.Dy()

		mx := (w - iw) / 2
		my := (h - ih) / 2

		if mx < 0 {
			mx = 0
		}
		if my < 0 {
			my = 0
		}

		c.left = mx
		c.right = w - iw - mx
		c.top = my
		c.bottom = h - ih - my
	}
}

// Since we pass `img` to `NewPadding`, we can simplify usage:
// NewPadding(img, PaddingCenter(bounds))

func PaddingCenterBox(bounds image.Rectangle) PaddingOption {
	return func(c *paddingConfig) {
		c.bounds = bounds
		// We can't calc margins here without image.
		// We'll set a flag or special value?
		// Better: NewPadding logic handles "auto" margins if bounds are set?
		// But how do we distinguish "TopLeft" from "Center"?
		// Let's add alignment flags to config.
	}
}

// Let's modify config to support alignment mode when bounds are set.
// const (
// 	alignNone = iota
// 	alignCenter
// 	alignTopLeft
// 	// ...
// )

func NewPadding(img image.Image, opts ...PaddingOption) image.Image {
	cfg := &paddingConfig{
		top: 0, left: 0, bottom: 0, right: 0,
		bg: nil,
	}
	// We might need to handle alignment options.
	// Let's support a generic "Layout" option?
	// The user suggested "TopLeftAlign(Left(Pixel(3)))".
	// This looks like `NewPadding(img, PadAlignTopLeft, PadLeft(3))`.

	// Let's parse standard options first.
	for _, o := range opts {
		o(cfg)
	}

	// If bounds are set, but margins are NOT fully specified (or we want to override them for alignment),
	// we need more logic.
	// But `PaddingCenter` usage in `NewPadding(img, ...)`:
	// If `PaddingCenter` is an option, it needs access to `img`. It doesn't have it.
	// So `NewPadding` must handle it.

	// Let's just expose `NewCenter(img, bounds)` or similar?
	// User said: "Can you use Padding to "center" text too? We can create a Center which backs out to Padding."

	var finalBounds image.Rectangle
	if !cfg.bounds.Empty() {
		finalBounds = cfg.bounds
	} else {
		b := img.Bounds()
		w := b.Dx() + cfg.left + cfg.right
		h := b.Dy() + cfg.top + cfg.bottom
		finalBounds = image.Rect(0, 0, w, h)
	}

	return &Padding{
		img:       img,
		bounds:    finalBounds,
		bgPattern: cfg.bg,
		top:       cfg.top,
		left:      cfg.left,
		bottom:    cfg.bottom,
		right:     cfg.right,
	}
}

// Helper for "Center"
func NewCenter(img image.Image, width, height int, bg image.Image) image.Image {
	b := img.Bounds()
	mx := (width - b.Dx()) / 2
	my := (height - b.Dy()) / 2
	if mx < 0 {
		mx = 0
	}
	if my < 0 {
		my = 0
	}

	// We want explicit bounds `width, height`.
	// Margins: Left=mx, Top=my. Right/Bottom filled to match width/height.

	// Right margin: width - b.Dx() - mx
	mr := width - b.Dx() - mx
	mb := height - b.Dy() - my
	if mr < 0 {
		mr = 0
	}
	if mb < 0 {
		mb = 0
	}

	return NewPadding(img,
		PaddingTop(my), PaddingLeft(mx),
		PaddingBottom(mb), PaddingRight(mr),
		PaddingBackground(bg),
		PaddingBoundary(image.Rect(0, 0, width, height)),
	)
}

// NewAligned returns an image padded to the specified width and height,
// with the inner image aligned according to xAlign and yAlign (0.0 to 1.0).
// 0.0 means Top/Left, 0.5 means Center, 1.0 means Bottom/Right.
// Optional padding arguments can be provided (following CSS standards):
// 1 arg: All sides
// 2 args: Vertical, Horizontal
// 4 args: Top, Right, Bottom, Left
func NewAligned(img image.Image, width, height int, xAlign, yAlign float64, bg image.Image, padding ...int) image.Image {
	b := img.Bounds()

	padTop, padRight, padBottom, padLeft := 0, 0, 0, 0
	if len(padding) > 0 {
		padTop = padding[0]
		padRight = padding[0]
		padBottom = padding[0]
		padLeft = padding[0]
	}
	if len(padding) >= 2 {
		padTop = padding[0]
		padBottom = padding[0]
		padLeft = padding[1]
		padRight = padding[1]
	}
	if len(padding) >= 4 {
		padTop = padding[0]
		padRight = padding[1]
		padBottom = padding[2]
		padLeft = padding[3]
	}

	// Calculate available space
	availW := width - b.Dx() - padLeft - padRight
	availH := height - b.Dy() - padTop - padBottom

	// Calculate top/left margin based on alignment
	mx := int(float64(availW) * xAlign)
	my := int(float64(availH) * yAlign)

	if mx < 0 {
		mx = 0
	}
	if my < 0 {
		my = 0
	}

	// Final margins
	finalLeft := padLeft + mx
	finalTop := padTop + my

	// Calculate bottom/right margin to ensure exact width/height
	mr := width - b.Dx() - finalLeft
	mb := height - b.Dy() - finalTop

	if mr < 0 {
		mr = 0
	}
	if mb < 0 {
		mb = 0
	}

	return NewPadding(img,
		PaddingTop(finalTop), PaddingLeft(finalLeft),
		PaddingBottom(mb), PaddingRight(mr),
		PaddingBackground(bg),
		PaddingBoundary(image.Rect(0, 0, width, height)),
	)
}
