package pattern

import (
	"fmt"
	"image"
	"image/color"
)

var (
	AmbientOcclusionOutputFilename = "ambient_occlusion.png"
	CurvatureOutputFilename        = "curvature.png"
)

func ExampleMaterial_basic() {
	// Base height map: A simple noise or shape
	// Let's use a Radial Gradient to simulate a sphere bump
	base := NewRadialGradient(
		GradientCenter(0.5, 0.5),
		SetStartColor(color.White), // High in center
		SetEndColor(color.Black),   // Low at edges
		// GradientRadius(0.4), // This option does not exist in RadialGradient, it uses bounds.
		// But we can simulate radius by scaling or bounds if we wanted.
		// RadialGradient defaults to filling the bounds.
	)

	// Create Material channels
	mat := Material{
		Height:    base,
		Albedo:    NewColorMap(base, ColorStop{0, color.RGBA{0, 0, 100, 255}}, ColorStop{1, color.RGBA{100, 200, 255, 255}}),
		Normal:    HeightToNormal(base, 2.0),
		Roughness: NewUniform(color.Gray{128}), // Uniform roughness
		Metalness: NewUniform(color.Black),     // Non-metal
		AO:        AOFromHeight(base, 4),
	}

	// We can also generate Curvature
	curv := CurvatureFromHeight(base)

	fmt.Println("Material created")
	_ = mat
	_ = curv

	// Output:
	// Material created
}

func ExampleNewAmbientOcclusion() {
	// This function is for documentation reference
	_ = GenerateAmbientOcclusion(image.Rect(0, 0, 200, 200))
}

func GenerateAmbientOcclusion(rect image.Rectangle) image.Image {
	base := NewWorleyNoise(
		SetFrequency(0.1),
	)
	return AOFromHeight(base, 2)
}

func ExampleNewCurvature() {
	// This function is for documentation reference
	_ = GenerateCurvature(image.Rect(0, 0, 200, 200))
}

func GenerateCurvature(rect image.Rectangle) image.Image {
	base := NewWorleyNoise(
		SetFrequency(0.1),
	)
	return CurvatureFromHeight(base)
}

func init() {
	RegisterGenerator("AmbientOcclusion", GenerateAmbientOcclusion)
	RegisterGenerator("Curvature", GenerateCurvature)
}
