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

func NewText(s string, fontSize float64, fg color.Color, bg color.Color) image.Image {
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
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
	height := int(fontSize * 1.5) // Approximate height

	// Create image
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	// Fill background
	if bg != nil {
		draw.Draw(img, rect, image.NewUniform(bg), image.Point{}, draw.Src)
	}

	// Draw text
	d.Dst = img
	d.Src = image.NewUniform(fg)
	// Base line position
	d.Dot = fixed.P(0, int(fontSize))
	d.DrawString(s)

	return img
}
