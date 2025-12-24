package pattern

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type TextOption func(*textConfig)

type textConfig struct {
	fontSize float64
	dpi      float64
	fg       image.Image
	bg       image.Image
}

func TextSize(size float64) TextOption {
	return func(c *textConfig) {
		c.fontSize = size
	}
}

func TextDPI(dpi float64) TextOption {
	return func(c *textConfig) {
		c.dpi = dpi
	}
}

func TextColor(fg image.Image) TextOption {
	return func(c *textConfig) {
		c.fg = fg
	}
}

func TextColorColor(fg color.Color) TextOption {
	return func(c *textConfig) {
		c.fg = image.NewUniform(fg)
	}
}

func TextBackgroundColor(bg image.Image) TextOption {
	return func(c *textConfig) {
		c.bg = bg
	}
}

func TextBackgroundColorColor(bg color.Color) TextOption {
	return func(c *textConfig) {
		c.bg = image.NewUniform(bg)
	}
}

func NewText(s string, opts ...TextOption) image.Image {
	cfg := &textConfig{
		fontSize: 24,
		dpi:      72,
		fg:       image.NewUniform(color.Black),
		bg:       nil,
	}
	for _, o := range opts {
		o(cfg)
	}

	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    cfg.fontSize,
		DPI:     cfg.dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}

	// Measure text
	d := &font.Drawer{
		Face: face,
	}
	width := d.MeasureString(s).Ceil()
	height := int(cfg.fontSize * 1.5) // Approximate height

	// Create image
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	// Fill background
	if cfg.bg != nil {
		draw.Draw(img, rect, cfg.bg, image.Point{}, draw.Src)
	}

	// Draw text
	d.Dst = img
	d.Src = cfg.fg
	// Base line position
	d.Dot = fixed.P(0, int(cfg.fontSize))
	d.DrawString(s)

	return img
}
