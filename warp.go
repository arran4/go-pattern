package pattern

import (
	"image"
	"image/color"
)

// Ensure Warp implements the image.Image interface.
var _ image.Image = (*Warp)(nil)

// Warp distorts the coordinates of the Source image using the Distortion image.
// It maps the color intensity of the Distortion image to a coordinate offset.
type Warp struct {
	Null
	Source          image.Image
	Distortion      image.Image
	DistortionX     image.Image
	DistortionY     image.Image
	Scale           float64
	XScale          float64
	YScale          float64
	DistortionScale float64
}

// NewWarp creates a new Warp pattern.
// If only Distortion is provided, it displaces both X and Y using the same map (usually diagonal if not handled).
// However, typically you want different noise for X and Y, so DistortionX and DistortionY are preferred for independent axis warping.
func NewWarp(source image.Image, ops ...func(any)) image.Image {
	p := &Warp{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Source:          source,
		Scale:           20.0, // Default distortion magnitude
		DistortionScale: 1.0,  // Scale of the noise texture coordinates
	}
	for _, op := range ops {
		op(p)
	}
	// Default behavior: if bounds are empty (from Null), we might want to adopt source bounds.
	// But Null defaults to 255x255.
	// Let's check if Source has bounds?
	if p.Source != nil && p.Source.Bounds() != image.Rect(0, 0, 0, 0) {
		p.bounds = p.Source.Bounds()
	}

	return p
}

func (p *Warp) At(x, y int) color.Color {
	if p.Source == nil {
		return color.Transparent
	}

	dx, dy := 0.0, 0.0

	// If uniform Distortion map is provided
	if p.Distortion != nil {
		// Sample distortion map
		// Usually we sample it at (x, y), maybe scaled.
		// If DistortionScale is used, we might need floating point sampling?
		// But existing patterns take Int.
		// So we just sample At(x, y).

		c := p.Distortion.At(x, y)
		gray := color.GrayModel.Convert(c).(color.Gray)
		// Map [0, 255] to [-1, 1] or [0, 1]?
		// Typically [-0.5, 0.5] * Scale
		val := (float64(gray.Y)/255.0 - 0.5) * 2.0
		dx += val * p.Scale
		dy += val * p.Scale
	}

	if p.DistortionX != nil {
		c := p.DistortionX.At(x, y)
		gray := color.GrayModel.Convert(c).(color.Gray)
		val := (float64(gray.Y)/255.0 - 0.5) * 2.0
		scale := p.XScale
		if scale == 0 {
			scale = p.Scale
		}
		dx += val * scale
	}

	if p.DistortionY != nil {
		c := p.DistortionY.At(x, y)
		gray := color.GrayModel.Convert(c).(color.Gray)
		val := (float64(gray.Y)/255.0 - 0.5) * 2.0
		scale := p.YScale
		if scale == 0 {
			scale = p.Scale
		}
		dy += val * scale
	}

	// Sample source at displaced coordinates
	srcX := int(float64(x) + dx)
	srcY := int(float64(y) + dy)

	return p.Source.At(srcX, srcY)
}

// WarpScale sets the global distortion scale (magnitude).
func WarpScale(scale float64) func(any) {
	return func(i any) {
		if p, ok := i.(*Warp); ok {
			p.Scale = scale
		}
	}
}

// WarpXScale sets the X distortion scale (magnitude).
func WarpXScale(scale float64) func(any) {
	return func(i any) {
		if p, ok := i.(*Warp); ok {
			p.XScale = scale
		}
	}
}

// WarpYScale sets the Y distortion scale (magnitude).
func WarpYScale(scale float64) func(any) {
	return func(i any) {
		if p, ok := i.(*Warp); ok {
			p.YScale = scale
		}
	}
}

// WarpDistortion sets the distortion map (affects both X and Y).
func WarpDistortion(img image.Image) func(any) {
	return func(i any) {
		if p, ok := i.(*Warp); ok {
			p.Distortion = img
		}
	}
}

// WarpDistortionX sets the distortion map for the X axis.
func WarpDistortionX(img image.Image) func(any) {
	return func(i any) {
		if p, ok := i.(*Warp); ok {
			p.DistortionX = img
		}
	}
}

// WarpDistortionY sets the distortion map for the Y axis.
func WarpDistortionY(img image.Image) func(any) {
	return func(i any) {
		if p, ok := i.(*Warp); ok {
			p.DistortionY = img
		}
	}
}
