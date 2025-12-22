package pattern

import (
	"image"

	"golang.org/x/image/draw"
)

// We need a wrapper that holds the Result image.
type ScaledImage struct {
	image.Image
}

// ScaleOption defines options for the Scale pattern
type ScaleOption func(*scaleConfig)

type scaleConfig struct {
	width, height int
	scaleX, scaleY float64
	scaler         draw.Scaler
}

func ScaleSize(w, h int) ScaleOption {
	return func(c *scaleConfig) {
		c.width = w
		c.height = h
	}
}

func ScaleFactor(f float64) ScaleOption {
	return func(c *scaleConfig) {
		c.scaleX = f
		c.scaleY = f
	}
}

func ScaleAlg(s draw.Scaler) ScaleOption {
	return func(c *scaleConfig) {
		c.scaler = s
	}
}

// NewScale creates a new scaled image.
// Note: This eagerly computes the scaled image because advanced interpolation requires neighborhood access.
func NewScale(img image.Image, opts ...ScaleOption) image.Image {
	cfg := &scaleConfig{
		scaler: draw.CatmullRom, // Default
		scaleX: 1.0,
		scaleY: 1.0,
	}
	for _, o := range opts {
		o(cfg)
	}

	bounds := img.Bounds()
	dstW, dstH := cfg.width, cfg.height

	if dstW == 0 && dstH == 0 {
		dstW = int(float64(bounds.Dx()) * cfg.scaleX)
		dstH = int(float64(bounds.Dy()) * cfg.scaleY)
	}

	dstRect := image.Rect(0, 0, dstW, dstH)
	dst := image.NewRGBA(dstRect)

	cfg.scaler.Scale(dst, dstRect, img, bounds, draw.Over, nil)

	return &ScaledImage{Image: dst}
}
