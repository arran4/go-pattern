package pattern

import (
	"image"
	"image/color"
)

// Ensure Blueprint implements the image.Image interface.
var _ image.Image = (*Blueprint)(nil)

// Blueprint renders a technical blueprint background with major and minor grid lines.
type Blueprint struct {
	Null
	BackgroundColor color.Color
	MinorLineColor  color.Color
	MajorLineColor  color.Color
	CellSize        int
	MajorCellSize   int
	LineWidth       int
	MajorLineWidth  int
}

func (b *Blueprint) ColorModel() color.Model {
	return color.RGBAModel
}

func (b *Blueprint) Bounds() image.Rectangle {
	return b.bounds
}

func (b *Blueprint) At(x, y int) color.Color {
	localX := x - b.bounds.Min.X
	localY := y - b.bounds.Min.Y

	cellSize := b.CellSize
	if cellSize <= 0 {
		cellSize = 10
	}

	majorCellSize := b.MajorCellSize
	if majorCellSize <= 0 {
		majorCellSize = 50
	}

	lineWidth := b.LineWidth
	if lineWidth <= 0 {
		lineWidth = 1
	}

	majorLineWidth := b.MajorLineWidth
	if majorLineWidth <= 0 {
		majorLineWidth = 2
	}

	// Check if we are on a major line
	isMajorX := false
	modXMajor := posMod(localX, majorCellSize)
	if modXMajor < majorLineWidth {
		isMajorX = true
	}

	isMajorY := false
	modYMajor := posMod(localY, majorCellSize)
	if modYMajor < majorLineWidth {
		isMajorY = true
	}

	if isMajorX || isMajorY {
		return b.MajorLineColor
	}

	// Check if we are on a minor line
	isMinorX := false
	modXMinor := posMod(localX, cellSize)
	if modXMinor < lineWidth {
		isMinorX = true
	}

	isMinorY := false
	modYMinor := posMod(localY, cellSize)
	if modYMinor < lineWidth {
		isMinorY = true
	}

	if isMinorX || isMinorY {
		return b.MinorLineColor
	}

	return b.BackgroundColor
}

// Options

type hasBlueprintBackgroundColor interface {
	SetBlueprintBackgroundColor(color.Color)
}

func SetBlueprintBackgroundColor(c color.Color) func(any) {
	return func(i any) {
		if v, ok := i.(hasBlueprintBackgroundColor); ok {
			v.SetBlueprintBackgroundColor(c)
		}
	}
}

func (b *Blueprint) SetBlueprintBackgroundColor(c color.Color) {
	b.BackgroundColor = c
}

type hasBlueprintMinorLineColor interface {
	SetBlueprintMinorLineColor(color.Color)
}

func SetBlueprintMinorLineColor(c color.Color) func(any) {
	return func(i any) {
		if v, ok := i.(hasBlueprintMinorLineColor); ok {
			v.SetBlueprintMinorLineColor(c)
		}
	}
}

func (b *Blueprint) SetBlueprintMinorLineColor(c color.Color) {
	b.MinorLineColor = c
}

type hasBlueprintMajorLineColor interface {
	SetBlueprintMajorLineColor(color.Color)
}

func SetBlueprintMajorLineColor(c color.Color) func(any) {
	return func(i any) {
		if v, ok := i.(hasBlueprintMajorLineColor); ok {
			v.SetBlueprintMajorLineColor(c)
		}
	}
}

func (b *Blueprint) SetBlueprintMajorLineColor(c color.Color) {
	b.MajorLineColor = c
}

type hasBlueprintCellSize interface {
	SetBlueprintCellSize(int)
}

func SetBlueprintCellSize(size int) func(any) {
	return func(i any) {
		if v, ok := i.(hasBlueprintCellSize); ok {
			v.SetBlueprintCellSize(size)
		}
	}
}

func (b *Blueprint) SetBlueprintCellSize(size int) {
	b.CellSize = size
}

type hasBlueprintMajorCellSize interface {
	SetBlueprintMajorCellSize(int)
}

func SetBlueprintMajorCellSize(size int) func(any) {
	return func(i any) {
		if v, ok := i.(hasBlueprintMajorCellSize); ok {
			v.SetBlueprintMajorCellSize(size)
		}
	}
}

func (b *Blueprint) SetBlueprintMajorCellSize(size int) {
	b.MajorCellSize = size
}

type hasBlueprintLineWidth interface {
	SetBlueprintLineWidth(int)
}

func SetBlueprintLineWidth(width int) func(any) {
	return func(i any) {
		if v, ok := i.(hasBlueprintLineWidth); ok {
			v.SetBlueprintLineWidth(width)
		}
	}
}

func (b *Blueprint) SetBlueprintLineWidth(width int) {
	b.LineWidth = width
}

type hasBlueprintMajorLineWidth interface {
	SetBlueprintMajorLineWidth(int)
}

func SetBlueprintMajorLineWidth(width int) func(any) {
	return func(i any) {
		if v, ok := i.(hasBlueprintMajorLineWidth); ok {
			v.SetBlueprintMajorLineWidth(width)
		}
	}
}

func (b *Blueprint) SetBlueprintMajorLineWidth(width int) {
	b.MajorLineWidth = width
}


// NewBlueprint creates a new Blueprint grid.
func NewBlueprint(ops ...func(any)) image.Image {
	b := &Blueprint{
		Null: Null{
			bounds: image.Rect(0, 0, 512, 512),
		},
		BackgroundColor: color.RGBA{R: 21, G: 67, B: 122, A: 255},
		MajorLineColor:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
		MinorLineColor:  color.RGBA{R: 124, G: 167, B: 214, A: 255},
		CellSize:        10,
		MajorCellSize:   50,
		LineWidth:       1,
		MajorLineWidth:  2,
	}

	for _, op := range ops {
		op(b)
	}
	return b
}
