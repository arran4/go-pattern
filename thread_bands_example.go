package pattern

import (
	"image"
	"image/color"
)

func ExampleNewThreadBands() {
	// Example body intentionally empty. The bootstrap tool inspects the signature and
	// metadata variables below to render documentation assets.
}

var (
	ThreadBandsOutputFilename = "thread_bands.png"
	ThreadBandsZoomLevels     = []int{}
	ThreadBandsBaseLabel      = "ThreadBands"
)

func init() {
	GlobalGenerators[ThreadBandsBaseLabel] = GenerateThreadBands
	GlobalReferences[ThreadBandsBaseLabel] = GenerateThreadBandsReferences
}

func GenerateThreadBands(rect image.Rectangle) image.Image {
	return NewThreadBands(
		SetBounds(rect),
		SetThreadWidth(14),
		SetCrossShadowDepth(0.35),
		SetColorVariance(0.08),
		SetLightThreadColor(color.RGBA{225, 218, 205, 255}),
		SetDarkThreadColor(color.RGBA{142, 134, 121, 255}),
		WithSeed(7),
	)
}

func GenerateThreadBandsReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	refs := map[string]func(image.Rectangle) image.Image{
		"Fine Weave": func(r image.Rectangle) image.Image {
			return NewThreadBands(
				SetBounds(r),
				SetThreadWidth(8),
				SetCrossShadowDepth(0.28),
				SetColorVariance(0.06),
				WithSeed(11),
			)
		},
		"Deep Shadow": func(r image.Rectangle) image.Image {
			return NewThreadBands(
				SetBounds(r),
				SetThreadWidth(16),
				SetCrossShadowDepth(0.55),
				SetColorVariance(0.1),
				SetLightThreadColor(color.RGBA{236, 228, 214, 255}),
				SetDarkThreadColor(color.RGBA{118, 108, 94, 255}),
				WithSeed(19),
			)
		},
	}

	return refs, []string{"Fine Weave", "Deep Shadow"}
}
