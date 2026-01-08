package pattern

import (
	"image"
	"image/png"
	"os"
)

var MudTracksOutputFilename = "mud_tracks.png"

const MudTracksBaseLabel = "MudTracks"

// ExampleNewMudTracks lays down multiple compacted bands with embedded pebble noise.
func ExampleNewMudTracks() {
	img := NewMudTracks()

	f, err := os.Create(MudTracksOutputFilename)
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

// GenerateMudTracks builds the pattern for registry-driven generation.
func GenerateMudTracks(b image.Rectangle) image.Image {
	return NewMudTracks(
		func(cfg *mudTracksConfig) {
			cfg.bounds = b
		},
	)
}

func init() {
	RegisterGenerator(MudTracksBaseLabel, GenerateMudTracks)
}
