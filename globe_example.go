package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	GlobeOutputFilename           = "globe.png"
	Globe_SimpleOutputFilename    = "globe_simple.png"
	Globe_ProjectedOutputFilename = "globe_projected.png"
	Globe_GridOutputFilename      = "globe_grid.png"
)

const GlobeBaseLabel = "Globe"

// GenerateEarthTexture creates a simple earth-like texture map
func GenerateEarthTexture(b image.Rectangle) image.Image {
	// Water base
	// Continents using Worley Noise (F1 Euclidean)
	// Vegetation using Perlin Noise

	// Use slightly different seeds/params than Islands to look global
	continents := NewWorleyNoise(
		SetBounds(b),
		SetFrequency(0.015),
		SetSeed(1001),
		SetWorleyOutput(OutputF1),
		SetWorleyMetric(MetricEuclidean),
	)

	// Detail for coastline/terrain
	detail := NewNoise(
		SetBounds(b),
		SetNoiseAlgorithm(&PerlinNoise{
			Seed:        2002,
			Frequency:   0.08,
			Octaves:     5,
			Persistence: 0.55,
			Lacunarity:  2.0,
		}),
	)

	mixed := NewBlend(continents, detail, BlendOverlay)

	// Map to Earth colors
	// 0.0 - 0.55: Ocean
	// 0.55 - 0.6: Coast/Sand
	// 0.6 - 0.8: Land/Forest
	// 0.8 - 0.9: Mountain
	// 0.9 - 1.0: Ice/Snow

	return NewColorMap(mixed,
		ColorStop{Position: 0.0, Color: color.RGBA{0, 0, 100, 255}},      // Deep Ocean
		ColorStop{Position: 0.50, Color: color.RGBA{0, 50, 150, 255}},    // Ocean
		ColorStop{Position: 0.55, Color: color.RGBA{30, 100, 180, 255}},  // Shallow Water
		ColorStop{Position: 0.58, Color: color.RGBA{210, 190, 150, 255}}, // Sand
		ColorStop{Position: 0.62, Color: color.RGBA{34, 139, 34, 255}},   // Forest
		ColorStop{Position: 0.80, Color: color.RGBA{85, 107, 47, 255}},   // Dark Green/Mountain Base
		ColorStop{Position: 0.90, Color: color.RGBA{100, 100, 100, 255}}, // Rock
		ColorStop{Position: 0.95, Color: color.RGBA{240, 240, 255, 255}}, // Snow
		ColorStop{Position: 1.0, Color: color.RGBA{255, 255, 255, 255}},  // Ice
	)
}

// ExampleNewGlobe_Simple demonstrates the "Circle and Texture" technique requested.
// It uses a flat circular mask over a terrain texture.
func ExampleNewGlobe_Simple() {
	g := GenerateGlobe_Simple(image.Rect(0, 0, 300, 300))
	saveImage(Globe_SimpleOutputFilename, g)
}

// ExampleNewGlobe_Projected demonstrates the true "Globe" pattern with spherical projection.
// It maps the same texture onto a sphere.
func ExampleNewGlobe_Projected() {
	g := GenerateGlobe_Projected(image.Rect(0, 0, 300, 300))
	saveImage(Globe_ProjectedOutputFilename, g)
}

// ExampleNewGlobe_Grid demonstrates the wireframe/grid mode of the Globe pattern.
func ExampleNewGlobe_Grid() {
	g := GenerateGlobe_Grid(image.Rect(0, 0, 300, 300))
	saveImage(Globe_GridOutputFilename, g)
}

// ExampleNewGlobe is the default example for the documentation.
func ExampleNewGlobe() {
	ExampleNewGlobe_Projected()
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func init() {
	RegisterGenerator(GlobeBaseLabel, GenerateGlobe)
	RegisterGenerator("Globe_Simple", GenerateGlobe_Simple)
	RegisterGenerator("Globe_Projected", GenerateGlobe_Projected)
	RegisterGenerator("Globe_Grid", GenerateGlobe_Grid)
}

func GenerateGlobe(b image.Rectangle) image.Image {
	return GenerateGlobe_Projected(b)
}

func GenerateGlobe_Simple(b image.Rectangle) image.Image {
	texture := GenerateEarthTexture(b)

	// Use Circle to mask it
	return NewCircle(
		SetBounds(b),
		SetFillImageSource(texture),
		SetSpaceColor(color.Transparent), // Transparent background
	)
}

func GenerateGlobe_Projected(b image.Rectangle) image.Image {
	// Texture needs to be equirectangular (2:1 aspect ratio ideally) for correct mapping,
	// but we'll use a square one and it will stretch at poles, which is fine for demo.
	// We generate a larger texture for better resolution
	texture := GenerateEarthTexture(image.Rect(0, 0, b.Dx()*2, b.Dy()))

	return NewGlobe(
		SetBounds(b),
		SetFillImageSource(texture),
		SetAngle(45), // Rotate to see some features
		SetTilt(20),  // Tilt axis
		SetSpaceColor(color.Transparent),
	)
}

func GenerateGlobe_Grid(b image.Rectangle) image.Image {
	return NewGlobe(
		SetBounds(b),
		SetLatitudeLines(12),
		SetLongitudeLines(24),
		SetAngle(30),
		SetTilt(15),
		SetLineColor(color.RGBA{0, 0, 255, 255}),
		SetSpaceColor(color.White),
	)
}
