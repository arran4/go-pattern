package pattern

import (
	"image"
	"image/color"
)

var (
	WorleyTilesOutputFilename = "tile_worley.png"
	WorleyTilesZoomLevels     = []int{}
	WorleyTilesBaseLabel      = "WorleyTiles"
)

func init() {
	RegisterGenerator(WorleyTilesBaseLabel, GenerateWorleyTiles)
}

// ExampleNewWorleyTiles tiles Worley/Voronoi stones with rounded edges and mortar.
// Parameters:
//   - stone size (pixels): SetTileStoneSize
//   - gap width (0-1): SetTileGapWidth
//   - color palette spread (0-1): SetTilePaletteSpread
func ExampleNewWorleyTiles() image.Image {
	baseTile := NewWorleyTiles(
		SetBounds(image.Rect(0, 0, 160, 160)),
		SetTileStoneSize(52),
		SetTileGapWidth(0.1),
		SetTilePaletteSpread(0.18),
		SetTilePalette(
			color.RGBA{128, 116, 106, 255},
			color.RGBA{146, 132, 118, 255},
			color.RGBA{112, 102, 96, 255},
		),
		WithSeed(2024),
	)

	// Tile the base stone field over a larger canvas so seams repeat cleanly.
	return NewTile(baseTile, image.Rect(0, 0, 320, 320))
}

// GenerateWorleyTiles wires the example into the registry for bootstrapping.
func GenerateWorleyTiles(bounds image.Rectangle) image.Image {
	baseTile := NewWorleyTiles(
		SetBounds(image.Rect(0, 0, 160, 160)),
		SetTileStoneSize(48),
		SetTileGapWidth(0.09),
		SetTilePaletteSpread(0.15),
		WithSeed(77),
	)
	return NewTile(baseTile, bounds)
}
