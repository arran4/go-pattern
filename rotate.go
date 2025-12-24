package pattern

import (
	"image"
	"image/color"
)

// Ensure Rotate implements the image.Image interface.
var _ image.Image = (*Rotate)(nil)

// Rotate is a pattern that rotates an underlying image by 90, 180, or 270 degrees.
type Rotate struct {
	img     image.Image
	degrees int
}

func (r *Rotate) ColorModel() color.Model {
	return r.img.ColorModel()
}

func (r *Rotate) Bounds() image.Rectangle {
	b := r.img.Bounds()
	w := b.Dx()
	h := b.Dy()

	switch r.degrees {
	case 90, 270:
		// Swap width and height.
		return image.Rect(b.Min.X, b.Min.Y, b.Min.X+h, b.Min.Y+w)
	case 180:
		return b
	default:
		return b
	}
}

func (r *Rotate) At(x, y int) color.Color {
	b := r.img.Bounds()
	// We need to map dest(x,y) to src(sx, sy).
	// Dest coords relative to Min:
	dx := x - r.Bounds().Min.X
	dy := y - r.Bounds().Min.Y
	w := b.Dx()
	h := b.Dy()

	var sx, sy int

	switch r.degrees {
	case 90:
		// 90 deg CW.
		// h is source height, which is dest width.
		// dx ranges 0..h-1.
		sx = dy
		sy = h - 1 - dx

	case 180:
		// (0,0) -> (w-1, h-1)
		sx = w - 1 - dx
		sy = h - 1 - dy

	case 270:
		// 270 CW (or 90 CCW).
		// w is source width, which is dest height.
		// dy ranges 0..w-1.
		sx = w - 1 - dy
		sy = dx

	default:
		sx = dx
		sy = dy
	}

	return r.img.At(b.Min.X+sx, b.Min.Y+sy)
}

// NewRotate creates a new Rotate from an existing image.
// degrees: 90, 180, 270 (values are normalized to these).
func NewRotate(img image.Image, degrees int, ops ...func(any)) image.Image {
	d := degrees % 360
	if d < 0 {
		d += 360
	}
	r := &Rotate{
		img:     img,
		degrees: d,
	}
	for _, op := range ops {
		op(r)
	}
	return r
}
