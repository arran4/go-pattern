package pattern

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var SpeedLinesOutputFilename = "speedlines.png"
var SpeedLinesZoomLevels = []int{}

const SpeedLinesOrder = 25
const SpeedLinesBaseLabel = "SpeedLines"

// SpeedLines Pattern
// Basic radial speed lines.
func ExampleNewSpeedLines() {
	i := NewSpeedLines(
		SetDensity(150),
		SetMinRadius(30),
		SetMaxRadius(80),
	)
	f, err := os.Create(SpeedLinesOutputFilename)
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

// ExampleNewSpeedLines_GopherBurst shows speed lines coming out of the Gopher's mouth.
func ExampleNewSpeedLines_GopherBurst() {
	gopher := NewGopher()

	// Gopher mouth is approx at (50, 75).
	// We want the lines to originate there.
	// We use Gopher as SpaceImageSource (background).
	// Lines are drawn on top.

	i := NewSpeedLines(
		func(i any) {
			if n, ok := i.(*Null); ok {
				// Match gopher bounds
				n.bounds = gopher.Bounds()
			}
		},
		SetCenter(50, 75),
		SetMinRadius(10), // Start close to mouth
		SetMaxRadius(30), // Variance
		SetDensity(200),
		SetSpaceImageSource(gopher),
		SetLineColor(color.Black),
	)

	f, err := os.Create("speedlines_gopher.png")
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

func GenerateSpeedLines(b image.Rectangle) image.Image {
	return NewSpeedLines(
		func(i any) {
			if n, ok := i.(*Null); ok {
				n.bounds = b
			}
		},
		SetDensity(100), SetMinRadius(20), SetMaxRadius(50))
}

func GenerateSpeedLinesReferences() (map[string]func(image.Rectangle) image.Image, []string) {
	return map[string]func(image.Rectangle) image.Image{
		"Radial": func(b image.Rectangle) image.Image {
			return NewSpeedLines(
				func(i any) {
					if n, ok := i.(*Null); ok {
						n.bounds = b
					}
				},
				SetDensity(100), SetMinRadius(20), SetMaxRadius(50))
		},
		"Linear": func(b image.Rectangle) image.Image {
			return NewSpeedLines(
				func(i any) {
					if n, ok := i.(*Null); ok {
						n.bounds = b
					}
				},
				SpeedLinesLinearType(), SetDensity(50))
		},
		"GopherBurst": func(b image.Rectangle) image.Image {
			gopher := NewGopher()
			// For GopherBurst, we want to match Gopher bounds, but if b is provided by generator,
			// usually we should respect it. However, the overlay depends on Gopher image size.
			// Let's stick to Gopher bounds for the demo visual consistency, or center it in b.
			// The generator framework usually passes a standard box size (e.g. 200x200).
			// If we want to show the gopher, we should probably stick to gopher size or center the gopher.
			// I'll stick to the previous logic for GopherBurst as it's a specific "Reference" image,
			// but I'll ensure it returns a valid image.
			// Actually, the `b` passed to references is usually the size of the box in the readme table.
			// If I return a smaller/larger image, it might be scaled or cropped.
			// I'll leave GopherBurst as fixed size because it depends on the asset.
			return NewSpeedLines(
				func(i any) {
					if n, ok := i.(*Null); ok {
						n.bounds = gopher.Bounds()
					}
				},
				SetCenter(50, 75),
				SetMinRadius(15),
				SetMaxRadius(40),
				SetDensity(150),
				SetSpaceImageSource(gopher),
			)
		},
	}, []string{"Radial", "Linear", "GopherBurst"}
}

func init() {
	RegisterGenerator("SpeedLines", GenerateSpeedLines)
	RegisterReferences("SpeedLines", GenerateSpeedLinesReferences)
}
