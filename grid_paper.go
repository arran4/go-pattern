package pattern

import (
	"image"
	"image/color"
)

// Ensure GridPaper implements the image.Image interface.
var _ image.Image = (*GridPaper)(nil)

// GridPaper renders a grid paper background with major and minor grid lines.
type GridPaper struct {
	Null
	BackgroundColor color.Color
	MinorLineColor  color.Color
	MajorLineColor  color.Color
	CellSize        int
	MajorCellSize   int
	LineWidth       int
	MajorLineWidth  int
}

func (g *GridPaper) ColorModel() color.Model {
	return color.RGBAModel
}

func (g *GridPaper) Bounds() image.Rectangle {
	return g.bounds
}

func (g *GridPaper) At(x, y int) color.Color {
	localX := x - g.bounds.Min.X
	localY := y - g.bounds.Min.Y

	cellSize := g.CellSize
	if cellSize <= 0 {
		cellSize = 10
	}

	majorCellSize := g.MajorCellSize
	if majorCellSize <= 0 {
		majorCellSize = 50
	}

	lineWidth := g.LineWidth
	if lineWidth <= 0 {
		lineWidth = 1
	}

	majorLineWidth := g.MajorLineWidth
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
		return g.MajorLineColor
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
		return g.MinorLineColor
	}

	return g.BackgroundColor
}

// Options

type hasGridPaperBackgroundColor interface {
	SetGridPaperBackgroundColor(color.Color)
}

func SetGridPaperBackgroundColor(c color.Color) func(any) {
	return func(i any) {
		if v, ok := i.(hasGridPaperBackgroundColor); ok {
			v.SetGridPaperBackgroundColor(c)
		}
	}
}

func (g *GridPaper) SetGridPaperBackgroundColor(c color.Color) {
	g.BackgroundColor = c
}

type hasGridPaperMinorLineColor interface {
	SetGridPaperMinorLineColor(color.Color)
}

func SetGridPaperMinorLineColor(c color.Color) func(any) {
	return func(i any) {
		if v, ok := i.(hasGridPaperMinorLineColor); ok {
			v.SetGridPaperMinorLineColor(c)
		}
	}
}

func (g *GridPaper) SetGridPaperMinorLineColor(c color.Color) {
	g.MinorLineColor = c
}

type hasGridPaperMajorLineColor interface {
	SetGridPaperMajorLineColor(color.Color)
}

func SetGridPaperMajorLineColor(c color.Color) func(any) {
	return func(i any) {
		if v, ok := i.(hasGridPaperMajorLineColor); ok {
			v.SetGridPaperMajorLineColor(c)
		}
	}
}

func (g *GridPaper) SetGridPaperMajorLineColor(c color.Color) {
	g.MajorLineColor = c
}

type hasGridPaperCellSize interface {
	SetGridPaperCellSize(int)
}

func SetGridPaperCellSize(size int) func(any) {
	return func(i any) {
		if v, ok := i.(hasGridPaperCellSize); ok {
			v.SetGridPaperCellSize(size)
		}
	}
}

func (g *GridPaper) SetGridPaperCellSize(size int) {
	g.CellSize = size
}

type hasGridPaperMajorCellSize interface {
	SetGridPaperMajorCellSize(int)
}

func SetGridPaperMajorCellSize(size int) func(any) {
	return func(i any) {
		if v, ok := i.(hasGridPaperMajorCellSize); ok {
			v.SetGridPaperMajorCellSize(size)
		}
	}
}

func (g *GridPaper) SetGridPaperMajorCellSize(size int) {
	g.MajorCellSize = size
}

type hasGridPaperLineWidth interface {
	SetGridPaperLineWidth(int)
}

func SetGridPaperLineWidth(width int) func(any) {
	return func(i any) {
		if v, ok := i.(hasGridPaperLineWidth); ok {
			v.SetGridPaperLineWidth(width)
		}
	}
}

func (g *GridPaper) SetGridPaperLineWidth(width int) {
	g.LineWidth = width
}

type hasGridPaperMajorLineWidth interface {
	SetGridPaperMajorLineWidth(int)
}

func SetGridPaperMajorLineWidth(width int) func(any) {
	return func(i any) {
		if v, ok := i.(hasGridPaperMajorLineWidth); ok {
			v.SetGridPaperMajorLineWidth(width)
		}
	}
}

func (g *GridPaper) SetGridPaperMajorLineWidth(width int) {
	g.MajorLineWidth = width
}

// NewGridPaper creates a new GridPaper grid.
func NewGridPaper(ops ...func(any)) image.Image {
	g := &GridPaper{
		Null: Null{
			bounds: image.Rect(0, 0, 512, 512),
		},
		BackgroundColor: color.RGBA{R: 245, G: 245, B: 245, A: 255},
		MajorLineColor:  color.RGBA{R: 80, G: 120, B: 180, A: 255},
		MinorLineColor:  color.RGBA{R: 140, G: 180, B: 230, A: 255},
		CellSize:        10,
		MajorCellSize:   50,
		LineWidth:       1,
		MajorLineWidth:  2,
	}

	for _, op := range ops {
		op(g)
	}
	return g
}
