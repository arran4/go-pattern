package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure Brick implements image.Image
var _ image.Image = (*Brick)(nil)

// Brick is a pattern that simulates a brick wall with running bond layout.
// It supports configurable brick size, mortar size, row offset, and multiple brick textures.
type Brick struct {
	Null
	Width, Height int
	MortarSize    int
	Offset        float64
	BrickImages   []image.Image
	MortarImage   image.Image
	Seed
}

func (b *Brick) At(x, y int) color.Color {
	// Defaults
	width := b.Width
	if width <= 0 {
		width = 40
	}
	height := b.Height
	if height <= 0 {
		height = 20
	}
	mortar := b.MortarSize
	// Allow mortar to be 0
	if mortar < 0 {
		mortar = 2
	}

	cellW := width + mortar
	cellH := height + mortar

	// Calculate row index
	// Handle negative coordinates correctly
	row := int(math.Floor(float64(y) / float64(cellH)))

	// Calculate local Y within the cell (0 to cellH-1)
	localY := (y%cellH + cellH) % cellH

	// Determine row offset
	xOffset := 0.0
	// Standard running bond: offset every odd row
	// We can make this more complex if we want, but for now strict running bond
	// Use math.Abs to handle negative rows parity if needed, but row%2 works fine in Go (result is -1, 0, 1)
	if row%2 != 0 {
		xOffset = b.Offset * float64(cellW)
	}

	// Adjust x by offset
	effX := float64(x) - xOffset
	col := int(math.Floor(effX / float64(cellW)))
	floorEffX := int(math.Floor(effX))
	localX := (floorEffX%cellW + cellW) % cellW

	// Determine if we are in mortar or brick
	// Center the brick in the cell
	// Mortar is split: half on left/top, half on right/bottom
	mortarHalf := mortar / 2
	// If mortar is odd, the extra pixel goes to the end (right/bottom)
	// range: [mortarHalf, mortarHalf + width) is brick

	inMortar := false
	if localX < mortarHalf || localX >= mortarHalf+width {
		inMortar = true
	}
	if localY < mortarHalf || localY >= mortarHalf+height {
		inMortar = true
	}

	if inMortar {
		if b.MortarImage != nil {
			// Sample mortar in world space
			return b.MortarImage.At(x, y)
		}
		// Default mortar color
		return color.RGBA{200, 200, 200, 255}
	}

	// We are in a brick
	// Determine which brick texture to use
	brickIdx := 0
	if len(b.BrickImages) > 1 {
		// Stochastic selection based on row/col
		h := StableHash(col, row, uint64(b.Seed.Seed))
		brickIdx = int(uint(h) % uint(len(b.BrickImages)))
	}

	// Get the image
	var img image.Image
	if len(b.BrickImages) > 0 {
		img = b.BrickImages[brickIdx]
	}

	if img == nil {
		// Default brick color (Reddish)
		return color.RGBA{180, 50, 50, 255}
	}

	// Calculate coordinates inside the brick (0 to width-1, 0 to height-1)
	bx := localX - mortarHalf
	by := localY - mortarHalf

	// Map to image coordinates
	// We want to "stamp" the image onto the brick.
	// We assume the image provided is the texture for the brick face.
	// We should tile it or clamp it?
	// Tiling (Tile pattern logic) is safest if the texture is smaller/larger.
	ib := img.Bounds()
	iw, ih := ib.Dx(), ib.Dy()
	if iw == 0 || ih == 0 {
		return color.Transparent
	}

	// Simple tiling mapping
	ix := (bx%iw + iw) % iw
	iy := (by%ih + ih) % ih

	return img.At(ib.Min.X+ix, ib.Min.Y+iy)
}

// NewBrick creates a new Brick pattern.
func NewBrick(ops ...func(any)) image.Image {
	p := &Brick{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Width:      40,
		Height:     20,
		MortarSize: 4,
		Offset:     0.5,
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// Options

type BrickSize struct{ Width, Height int }

func (p *BrickSize) SetBrickSize(w, h int) { p.Width, p.Height = w, h }

type MortarSize struct{ Size int }

func (p *MortarSize) SetMortarSize(s int) { p.Size = s }

type BrickOffset struct{ Offset float64 }

func (p *BrickOffset) SetBrickOffset(o float64) { p.Offset = o }

// Interface adapters for options

type hasBrickSize interface{ SetBrickSize(w, h int) }

func SetBrickSize(w, h int) func(any) {
	return func(i any) {
		if x, ok := i.(hasBrickSize); ok {
			x.SetBrickSize(w, h)
		}
	}
}

type hasMortarSize interface{ SetMortarSize(int) }

func SetMortarSize(s int) func(any) {
	return func(i any) {
		if h, ok := i.(hasMortarSize); ok {
			h.SetMortarSize(s)
		}
	}
}

type hasBrickOffset interface{ SetBrickOffset(float64) }

func SetBrickOffset(o float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasBrickOffset); ok {
			h.SetBrickOffset(o)
		}
	}
}

type hasBrickImages interface{ SetBrickImages([]image.Image) }

func SetBrickImages(imgs ...image.Image) func(any) {
	return func(i any) {
		if h, ok := i.(hasBrickImages); ok {
			h.SetBrickImages(imgs)
		}
	}
}

type hasMortarImage interface{ SetMortarImage(image.Image) }

func SetMortarImage(img image.Image) func(any) {
	return func(i any) {
		if h, ok := i.(hasMortarImage); ok {
			h.SetMortarImage(img)
		}
	}
}

// Implement setters on Brick

func (b *Brick) SetBrickSize(w, h int)             { b.Width, b.Height = w, h }
func (b *Brick) SetMortarSize(s int)               { b.MortarSize = s }
func (b *Brick) SetBrickOffset(o float64)          { b.Offset = o }
func (b *Brick) SetBrickImages(imgs []image.Image) { b.BrickImages = imgs }
func (b *Brick) SetMortarImage(img image.Image)    { b.MortarImage = img }
