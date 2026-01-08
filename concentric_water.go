package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure ConcentricWater implements the image.Image interface.
var _ image.Image = (*ConcentricWater)(nil)

// ConcentricWater renders concentric distance-field ripples with sine-driven height
// that modulates both tint and inferred normals for a water-like look.
type ConcentricWater struct {
	Null
	Center

	RingSpacing      float64
	Amplitude        float64
	AmplitudeFalloff float64
	BaseTint         color.RGBA
	NormalStrength   float64
	lightDir         [3]float64
	normalizedLight  [3]float64
	lightInitialized bool
}

func (cw *ConcentricWater) ColorModel() color.Model {
	return color.RGBAModel
}

func (cw *ConcentricWater) Bounds() image.Rectangle {
	return cw.Null.Bounds()
}

func (cw *ConcentricWater) At(x, y int) color.Color {
	h := cw.heightAt(float64(x), float64(y))
	nx, ny, nz := cw.normalAt(x, y)

	light := cw.dotLight(nx, ny, nz)
	foam := math.Pow(math.Max(0, h), 3.0)

	shade := 0.55 + 0.45*(0.5+h*0.5)
	highlight := 0.35 + 0.65*light

	base := cw.BaseTint
	br := float64(base.R) / 255.0
	bg := float64(base.G) / 255.0
	bb := float64(base.B) / 255.0

	r := clamp01((br*shade + foam*0.25 + light*0.1) * highlight)
	g := clamp01((bg*shade + foam*0.20 + light*0.1) * highlight)
	b := clamp01((bb*shade + foam*0.30 + light*0.1) * highlight)

	return color.RGBA{
		R: uint8(r*255 + 0.5),
		G: uint8(g*255 + 0.5),
		B: uint8(b*255 + 0.5),
		A: 255,
	}
}

func (cw *ConcentricWater) heightAt(x, y float64) float64 {
	spacing := cw.RingSpacing
	if spacing <= 0 {
		spacing = 16
	}
	amplitude := cw.Amplitude
	if amplitude == 0 {
		amplitude = 1
	}
	falloff := cw.AmplitudeFalloff
	if falloff <= 0 {
		falloff = 0.012
	}

	dx := x - float64(cw.CenterX)
	dy := y - float64(cw.CenterY)
	r := math.Hypot(dx, dy)

	wave := math.Sin((r / spacing) * 2 * math.Pi)
	atten := math.Exp(-falloff * r)

	return amplitude * wave * atten
}

func (cw *ConcentricWater) normalAt(x, y int) (float64, float64, float64) {
	strength := cw.NormalStrength
	if strength == 0 {
		strength = 3.0
	}

	hL := cw.heightAt(float64(x-1), float64(y))
	hR := cw.heightAt(float64(x+1), float64(y))
	hU := cw.heightAt(float64(x), float64(y-1))
	hD := cw.heightAt(float64(x), float64(y+1))

	dx := (hR - hL) * strength
	dy := (hD - hU) * strength
	dz := 1.0

	length := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if length == 0 {
		return 0, 0, 1
	}

	return -dx / length, -dy / length, dz / length
}

func (cw *ConcentricWater) initLightDir() {
	if cw.lightInitialized {
		return
	}
	dir := cw.lightDir
	if dir == [3]float64{} {
		dir = [3]float64{0.3, -0.4, 0.85}
	}
	length := math.Sqrt(dir[0]*dir[0] + dir[1]*dir[1] + dir[2]*dir[2])
	if length == 0 {
		dir = [3]float64{0, 0, 1}
		length = 1
	}
	cw.normalizedLight = [3]float64{dir[0] / length, dir[1] / length, dir[2] / length}
	cw.lightInitialized = true
}

func (cw *ConcentricWater) dotLight(nx, ny, nz float64) float64 {
	cw.initLightDir()
	dl := cw.normalizedLight
	d := nx*dl[0] + ny*dl[1] + nz*dl[2]
	if d < 0 {
		return 0
	}
	if d > 1 {
		return 1
	}
	return d
}

// NewConcentricWater creates a concentric ripple height field rendered with water tinting.
func NewConcentricWater(ops ...func(any)) image.Image {
	p := &ConcentricWater{
		Null:             Null{bounds: image.Rect(0, 0, 255, 255)},
		Center:           Center{CenterX: 127, CenterY: 127},
		RingSpacing:      16,
		Amplitude:        1,
		AmplitudeFalloff: 0.012,
		BaseTint:         color.RGBA{30, 110, 170, 255},
		NormalStrength:   3.0,
		lightDir:         [3]float64{0.3, -0.4, 0.85},
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// ConcentricWaterRingSpacing sets the spacing between ripples.
func ConcentricWaterRingSpacing(spacing float64) func(any) {
	return func(i any) {
		if cw, ok := i.(*ConcentricWater); ok {
			cw.RingSpacing = spacing
		}
	}
}

// ConcentricWaterAmplitude sets the sine wave amplitude.
func ConcentricWaterAmplitude(amplitude float64) func(any) {
	return func(i any) {
		if cw, ok := i.(*ConcentricWater); ok {
			cw.Amplitude = amplitude
		}
	}
}

// ConcentricWaterAmplitudeFalloff controls how quickly the ripple amplitude decays with distance.
func ConcentricWaterAmplitudeFalloff(falloff float64) func(any) {
	return func(i any) {
		if cw, ok := i.(*ConcentricWater); ok {
			cw.AmplitudeFalloff = falloff
		}
	}
}

// ConcentricWaterBaseTint sets the base water color.
func ConcentricWaterBaseTint(tint color.RGBA) func(any) {
	return func(i any) {
		if cw, ok := i.(*ConcentricWater); ok {
			cw.BaseTint = tint
		}
	}
}

// ConcentricWaterNormalStrength sets the slope scaling used for normal estimation.
func ConcentricWaterNormalStrength(strength float64) func(any) {
	return func(i any) {
		if cw, ok := i.(*ConcentricWater); ok {
			cw.NormalStrength = strength
		}
	}
}

// ConcentricWaterLightDirection sets the light vector used for shading.
func ConcentricWaterLightDirection(x, y, z float64) func(any) {
	return func(i any) {
		if cw, ok := i.(*ConcentricWater); ok {
			cw.lightDir = [3]float64{x, y, z}
			cw.lightInitialized = false
		}
	}
}
