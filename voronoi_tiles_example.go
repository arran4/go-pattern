package pattern

import (
	"image"
	"image/png"
	"os"
)

var VoronoiTilesOutputFilename = "voronoi_tiles.png"
var VoronoiTilesZoomLevels = []int{}

const VoronoiTilesBaseLabel = "Voronoi_tiles"

// Voronoi Tiles
// Uses Voronoi cells to define tiles, raises the centers, darkens the gaps, and sprinkles dust noise.
func ExampleNewVoronoiTiles() {
	img := NewVoronoiTiles(image.Rect(0, 0, 255, 255), defaultVoronoiTileCellSize, defaultVoronoiTileGapWidth, defaultVoronoiTileHeightImpact, 2024)

	f, err := os.Create(VoronoiTilesOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}

func GenerateVoronoiTiles(b image.Rectangle) image.Image {
	return NewVoronoiTiles(b, defaultVoronoiTileCellSize, defaultVoronoiTileGapWidth, defaultVoronoiTileHeightImpact, 2024)
}

func init() {
	RegisterGenerator("VoronoiTiles", GenerateVoronoiTiles)
}
