package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var VoronoiOutputFilename = "voronoi.png"
var VoronoiZoomLevels = []int{}

const VoronoiBaseLabel = "Voronoi"

// Voronoi Pattern
// Generates Voronoi cells.
func ExampleNewVoronoi() {
	// Define some points and colors
	points := []image.Point{
		{50, 50}, {200, 50}, {125, 125}, {50, 200}, {200, 200},
	}
	colors := []color.Color{
		color.RGBA{255, 100, 100, 255},
		color.RGBA{100, 255, 100, 255},
		color.RGBA{100, 100, 255, 255},
		color.RGBA{255, 255, 100, 255},
		color.RGBA{100, 255, 255, 255},
	}

	i := NewVoronoi(points, colors)
	f, err := os.Create(VoronoiOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func GenerateVoronoi(b image.Rectangle) image.Image {
	return NewDemoVoronoi(SetBounds(b))
}

func init() {
	RegisterGenerator(VoronoiBaseLabel, GenerateVoronoi)
}
