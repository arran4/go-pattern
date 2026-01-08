package pattern

import (
	"image"
)

var MetalPanelGrungeOutputFilename = "metal_panel_grunge.png"

// ExampleNewMetalPanelGrunge showcases the brushed metal panel with grunge seams.
func ExampleNewMetalPanelGrunge() image.Image {
	return NewMetalPanelGrunge()
}

// GenerateMetalPanelGrunge honors the provided bounds while keeping the default look.
func GenerateMetalPanelGrunge(rect image.Rectangle) image.Image {
	return NewMetalPanelGrunge(SetBounds(rect))
}

func init() {
	GlobalGenerators["MetalPanelGrunge"] = GenerateMetalPanelGrunge
}
