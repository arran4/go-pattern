package pattern

import (
	"image"
	"image/color"
)

// Ensure Remap implements the image.Image interface.
var _ image.Image = (*Remap)(nil)

// Remap re-maps the coordinates of the Source image using the Map image (UV mapping).
// The Map image Red channel defines the X coordinate (normalized 0..1).
// The Map image Green channel defines the Y coordinate (normalized 0..1).
type Remap struct {
	Null
	Source image.Image
	Map    image.Image
}

func (r *Remap) ColorModel() color.Model {
	if r.Source != nil {
		return r.Source.ColorModel()
	}
	return color.RGBAModel
}

func (r *Remap) Bounds() image.Rectangle {
	if !r.Null.Bounds().Empty() {
		return r.Null.Bounds()
	}
	if r.Map != nil {
		return r.Map.Bounds()
	}
	return image.Rect(0, 0, 0, 0)
}

func (r *Remap) At(x, y int) color.Color {
	if r.Source == nil || r.Map == nil {
		return color.Transparent
	}

	// Get UV from Map
	c := r.Map.At(x, y)
	// We need 16-bit components. RGBA() returns alpha-premultiplied values in [0, 65535].
	cr, cg, _, ca := c.RGBA()

	if ca == 0 {
		return color.Transparent
	}

	// Normalize to 0..1
	// Since values are premultiplied, we should technically divide by alpha if we want "true" color,
	// but usually UV maps are opaque. If it's semi-transparent, what does it mean?
	// If we treat R/G as coordinates, 0 opacity usually means "no mapping".
	// We handle full transparency above.
	// For partial transparency, let's treat the premultiplied value as the coordinate directly (weighted).
	// Or should we unpremultiply?
	// UV map values are usually pure data.
	// If we use R,G directly, they are already scaled by Alpha.
	// If Alpha is 0, R,G are 0.
	// If Alpha is 65535, R,G are full value.
	// Let's assume the user provides opaque maps or we just use the raw value.
	// But RGBA() returns premultiplied.
	// If the map is intended to be a UV map, it should probably be opaque.
	// If it is transparent, we return transparent.

	u := float64(cr) / 65535.0
	v := float64(cg) / 65535.0

	// Map to Source bounds
	sb := r.Source.Bounds()
	w := float64(sb.Dx())
	h := float64(sb.Dy())

	// Nearest neighbor mapping
	// Map 0.0 to Min.X, 1.0 to Max.X (exclusive-ish)
	// Usually 0.0 -> pixel 0. 1.0 -> pixel width.
	// So index = floor(u * w).
	sx := int(u * w) + sb.Min.X
	sy := int(v * h) + sb.Min.Y

	// Clamp to valid range to be safe (e.g. if u=1.0)
	if sx >= sb.Max.X {
		sx = sb.Max.X - 1
	}
	if sy >= sb.Max.Y {
		sy = sb.Max.Y - 1
	}
	if sx < sb.Min.X {
		sx = sb.Min.X
	}
	if sy < sb.Min.Y {
		sy = sb.Min.Y
	}

	return r.Source.At(sx, sy)
}

// NewRemap creates a new Remap pattern.
// source: The image to sample from.
func NewRemap(source image.Image, ops ...func(any)) image.Image {
	r := &Remap{
		Source: source,
	}
	for _, op := range ops {
		op(r)
	}
	return r
}

// RemapMap sets the map image (UV map) for the Remap pattern.
func RemapMap(m image.Image) func(any) {
	return func(i any) {
		if r, ok := i.(*Remap); ok {
			r.Map = m
		}
	}
}
