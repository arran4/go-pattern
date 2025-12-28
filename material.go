package pattern

import (
	"image"
	"image/color"
	"math"
)

// AmbientOcclusion calculates AO from a height map using a sampling kernel.
type AmbientOcclusion struct {
	Source image.Image
	Radius int
}

func (ao *AmbientOcclusion) At(x, y int) color.Color {
	if ao.Radius <= 0 {
		return color.White
	}

	// This is a simplified 2D AO approximation.
	// We check if neighbors are higher than the current pixel (plus a bias)
	// and if so, how much they "occlude".

	centerHeight := getLuminanceForMaterial(ao.Source.At(x, y))

	totalOcclusion := 0.0
	samples := 0.0

	// We can't iterate too many pixels if Radius is large, but for 2D patterns
	// a small radius (1-5) is usually intended, or we can use a sparse sampling.
	// For "Radius", we can just iterate the box.

	for j := -ao.Radius; j <= ao.Radius; j++ {
		for i := -ao.Radius; i <= ao.Radius; i++ {
			if i == 0 && j == 0 {
				continue
			}

			// Distance check for circular kernel
			distSq := float64(i*i + j*j)
			rSq := float64(ao.Radius * ao.Radius)
			if distSq > rSq {
				continue
			}

			h := getLuminanceForMaterial(ao.Source.At(x+i, y+j))

			// Logic: If neighbor is higher, it might occlude.
			// The occlusion depends on the height difference and distance.
			// Higher neighbor closer = more occlusion.

			diff := h - centerHeight
			if diff > 0 {
				dist := math.Sqrt(distSq)
				// Basic attenuation
				occlusion := diff / (1.0 + dist)
				totalOcclusion += occlusion
			}
			samples++
		}
	}

	// Normalize
	// The max occlusion per sample is < 1.0 (since diff < 1.0 and dist >= 1.0).
	// We average it.

	aoVal := 1.0
	if samples > 0 {
		// Just a heuristic.
		// totalOcclusion sum can be up to samples * 1.0 roughly.
		// We want 0.0 (fully occluded) to 1.0 (open).
		// Currently totalOcclusion is "how much stuff is blocking".
		// We can try to normalize it.
		// Let's multiply by a factor to make it visible.

		avgOcclusion := totalOcclusion / samples
		// Amplify it a bit
		aoVal = 1.0 - (avgOcclusion * 4.0)
	}

	if aoVal < 0 {
		aoVal = 0
	}
	if aoVal > 1 {
		aoVal = 1
	}

	c := uint8(aoVal * 255)
	return color.Gray{Y: c}
}

func (ao *AmbientOcclusion) Bounds() image.Rectangle {
	if ao.Source == nil {
		return image.Rect(0, 0, 255, 255)
	}
	return ao.Source.Bounds()
}

func (ao *AmbientOcclusion) ColorModel() color.Model {
	return color.GrayModel
}

func NewAmbientOcclusion(source image.Image, ops ...func(any)) image.Image {
	ao := &AmbientOcclusion{
		Source: source,
		Radius: 2,
	}
	for _, op := range ops {
		op(ao)
	}
	return ao
}

func AmbientOcclusionRadius(radius int) func(any) {
	return func(i any) {
		if ao, ok := i.(*AmbientOcclusion); ok {
			ao.Radius = radius
		}
	}
}

// Curvature calculates curvature (convex/concave) from a height map.
type Curvature struct {
	Source image.Image
}

func (c *Curvature) At(x, y int) color.Color {
	// Laplacian operator
	//  0 -1  0
	// -1  4 -1
	//  0 -1  0

	h := getLuminanceForMaterial(c.Source.At(x, y))
	hUp := getLuminanceForMaterial(c.Source.At(x, y-1))
	hDown := getLuminanceForMaterial(c.Source.At(x, y+1))
	hLeft := getLuminanceForMaterial(c.Source.At(x-1, y))
	hRight := getLuminanceForMaterial(c.Source.At(x+1, y))

	laplacian := 4*h - hUp - hDown - hLeft - hRight

	// Laplacian is 0 for flat, positive for convex (peak), negative for concave (valley) relative to neighbors?
	// Actually:
	// Peak: Center is high (1), neighbors low (0). 4*1 - 0 = 4. Positive.
	// Valley: Center is low (0), neighbors high (1). 4*0 - 4 = -4. Negative.

	// Usually Curvature maps are Gray where 0.5 is flat.
	// Convex (peaks) -> White (1.0)
	// Concave (valleys) -> Black (0.0)

	val := (laplacian * 0.5) + 0.5

	if val < 0 {
		val = 0
	}
	if val > 1 {
		val = 1
	}

	g := uint8(val * 255)
	return color.Gray{Y: g}
}

func (c *Curvature) Bounds() image.Rectangle {
	if c.Source == nil {
		return image.Rect(0, 0, 255, 255)
	}
	return c.Source.Bounds()
}

func (c *Curvature) ColorModel() color.Model {
	return color.GrayModel
}

func NewCurvature(source image.Image, ops ...func(any)) image.Image {
	curv := &Curvature{
		Source: source,
	}
	for _, op := range ops {
		op(curv)
	}
	return curv
}

// Material is a collection of images representing different PBR channels.
type Material struct {
	Albedo    image.Image
	Normal    image.Image
	Roughness image.Image
	Metalness image.Image
	AO        image.Image
	Height    image.Image
}

// getLuminanceForMaterial is a helper to calculate luminance.
func getLuminanceForMaterial(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0
}

// HeightToNormal is a utility function that wraps NewNormalMap to creates a normal map from a height map with the given strength.
func HeightToNormal(source image.Image, strength float64) image.Image {
	return NewNormalMap(source, NormalMapStrength(strength))
}

// AOFromHeight is a utility function that wraps NewAmbientOcclusion to creates an ambient occlusion map from a height map with the given radius.
func AOFromHeight(source image.Image, radius int) image.Image {
	return NewAmbientOcclusion(source, AmbientOcclusionRadius(radius))
}

// CurvatureFromHeight is a utility function that wraps NewCurvature to creates a curvature map from a height map.
func CurvatureFromHeight(source image.Image) image.Image {
	return NewCurvature(source)
}
