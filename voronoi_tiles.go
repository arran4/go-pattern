package pattern

import (
	"image"
	"image/color"
	"math"
	mrand "math/rand"
)

// Default values for the Voronoi tile material.
const (
	defaultVoronoiTileCellSize     = 48.0
	defaultVoronoiTileGapWidth     = 0.08
	defaultVoronoiTileHeightImpact = 0.35
)

// NewVoronoiTiles builds a tileable material using Voronoi cells as tiles.
// Cell centers are lifted slightly, edges are darkened, and subtle dust noise is overlaid.
// Parameters:
//   - cellSize:    Distance in pixels between Voronoi sites.
//   - gapWidth:    How wide (0-1) the dark gap between tiles appears. Higher widens the shadow band.
//   - heightBoost: How much to lighten the tile centers and deepen the edges.
func NewVoronoiTiles(bounds image.Rectangle, cellSize, gapWidth, heightBoost float64, seed int64) image.Image {
	if bounds == (image.Rectangle{}) {
		bounds = image.Rect(0, 0, 255, 255)
	}
	if cellSize <= 0 {
		cellSize = defaultVoronoiTileCellSize
	}
	if gapWidth < 0 {
		gapWidth = 0
	}
	if gapWidth > 1 {
		gapWidth = 1
	}
	if heightBoost <= 0 {
		heightBoost = defaultVoronoiTileHeightImpact
	}

	points := jitteredVoronoiPoints(bounds, cellSize, seed)
	tilePalette := []color.Color{
		color.RGBA{182, 160, 140, 255},
		color.RGBA{174, 146, 126, 255},
		color.RGBA{188, 168, 148, 255},
		color.RGBA{170, 142, 122, 255},
	}
	tiles := NewVoronoi(points, tilePalette, SetBounds(bounds))

	// F1 gives distance to closest point; use it as a height ramp to lift cell centers.
	freq := 1.0 / cellSize
	centerDistance := NewWorleyNoise(
		SetSeed(seed),
		SetFrequency(freq),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyOutput(OutputF1),
		SetWorleyJitter(0.85),
		SetBounds(bounds),
	)

	highlight := clampByte(164 + int(90*heightBoost))
	midtone := clampByte(148)
	edgeTone := clampByte(132 - int(55*heightBoost))

	edgeStart := math.Max(0.65, 1.0-(gapWidth*0.65))
	heightMap := NewColorMap(centerDistance,
		ColorStop{Position: 0.0, Color: color.Gray{Y: highlight}},
		ColorStop{Position: 0.45, Color: color.Gray{Y: midtone}},
		ColorStop{Position: edgeStart, Color: color.Gray{Y: edgeTone}},
		ColorStop{Position: 1.0, Color: color.Gray{Y: edgeTone}},
	)

	// F2-F1 is thin near edges; map it to a dark mortar/gap band.
	gapSignal := NewWorleyNoise(
		SetSeed(seed+7),
		SetFrequency(freq),
		SetWorleyMetric(MetricEuclidean),
		SetWorleyOutput(OutputF2MinusF1),
		SetWorleyJitter(0.85),
		SetBounds(bounds),
	)

	gapTone := clampByte(96 - int(50*heightBoost))
	gapStops := []ColorStop{
		{Position: 0.0, Color: color.Gray{Y: gapTone}},
		{Position: math.Min(1, gapWidth*1.2+0.02), Color: color.Gray{Y: edgeTone}},
		{Position: 1.0, Color: color.White},
	}
	gaps := NewColorMap(gapSignal, gapStops...)

	raisedTiles := NewBlend(tiles, heightMap, BlendOverlay, SetBounds(bounds))
	shadedTiles := NewBlend(raisedTiles, gaps, BlendMultiply, SetBounds(bounds))

	// Dust noise to break up flat areas.
	dustNoise := NewNoise(
		NoiseSeed(seed+19),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        seed + 19,
			Frequency:   freq * 3.2,
			Octaves:     3,
			Persistence: 0.6,
		}),
		SetBounds(bounds),
	)
	dust := NewColorMap(dustNoise,
		ColorStop{Position: 0.0, Color: color.RGBA{225, 220, 212, 255}},
		ColorStop{Position: 0.4, Color: color.RGBA{236, 232, 222, 255}},
		ColorStop{Position: 1.0, Color: color.RGBA{250, 246, 238, 255}},
	)

	return NewBlend(shadedTiles, dust, BlendOverlay, SetBounds(bounds))
}

func jitteredVoronoiPoints(bounds image.Rectangle, cellSize float64, seed int64) []image.Point {
	rng := mrand.New(mrand.NewSource(seed))
	jitter := cellSize * 0.35

	var pts []image.Point
	for y := float64(bounds.Min.Y); y < float64(bounds.Max.Y); y += cellSize {
		for x := float64(bounds.Min.X); x < float64(bounds.Max.X); x += cellSize {
			cx := x + cellSize*0.5 + (rng.Float64()-0.5)*2*jitter
			cy := y + cellSize*0.5 + (rng.Float64()-0.5)*2*jitter

			cx = math.Max(float64(bounds.Min.X), math.Min(cx, float64(bounds.Max.X-1)))
			cy = math.Max(float64(bounds.Min.Y), math.Min(cy, float64(bounds.Max.Y-1)))
			pts = append(pts, image.Pt(int(cx), int(cy)))
		}
	}

	if len(pts) == 0 {
		pts = append(pts, image.Pt(bounds.Min.X+bounds.Dx()/2, bounds.Min.Y+bounds.Dy()/2))
	}

	return pts
}

func clampByte(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > math.MaxUint8 {
		return math.MaxUint8
	}
	return uint8(v)
}
