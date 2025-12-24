package pattern

import (
	"image"
)

var (
	TileOutputFilename = "tile.png"
	TileZoomLevels     = []int{}
	TileOrder          = 6
	TileBaseLabel      = "Tile"
)

func init() {
	RegisterGenerator("Tile", func(bounds image.Rectangle) image.Image {
		return ExampleNewTile(SetBounds(bounds))
	})
	RegisterReferences("Tile", BootstrapTileReferences)
}

func BootstrapTileReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Source": func(bounds image.Rectangle) image.Image {
			return NewGopher()
		},
	}, []string{"Source"}
}

func ExampleNewTile(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleToRatio(0.25))
	// Tile the gopher in a 200x200 area
	return NewTile(gopher, image.Rect(0, 0, 200, 200))
}
