package pattern

import (
	"image"
	"image/color"
	"math"
)

// Ensure VHS implements the image.Image interface.
var _ image.Image = (*VHS)(nil)

// VHS applies a retro VHS effect with scanlines, chromatic aberration, and noise.
type VHS struct {
	Null
	Image             image.Image
	ScanlineFrequency float64 // Frequency of the scanlines (e.g. 0.5 for every other line)
	ScanlineIntensity float64 // Intensity of the scanline darkening (0.0 to 1.0)
	ColorOffset       int     // Pixel offset for the red/blue channels
	NoiseIntensity    float64 // Intensity of the noise (0.0 to 1.0)
	Seed              int64   // Seed for the deterministic noise
	noise             *HashNoise
}

func (p *VHS) At(x, y int) color.Color {
	if p.Image == nil {
		return color.RGBA{}
	}

	// 1. Chromatic Aberration (Channel Offset)
	// Get source colors
	rCol := p.Image.At(x-p.ColorOffset, y)
	gCol := p.Image.At(x, y)
	bCol := p.Image.At(x+p.ColorOffset, y)

	// Convert to NRGBA to work with straight alpha
	rN := color.NRGBAModel.Convert(rCol).(color.NRGBA)
	gN := color.NRGBAModel.Convert(gCol).(color.NRGBA)
	bN := color.NRGBAModel.Convert(bCol).(color.NRGBA)

	// Use Green's alpha as the main alpha
	alpha := gN.A
	// Or maybe max alpha? For "bad signal" usually alpha is consistent if it's a layer.
	// If the image has transparency, shifting channels might shift alpha too.
	// Let's stick to using the center pixel's alpha to define the shape.

	r8 := rN.R
	g8 := gN.G
	b8 := bN.B

	// 2. Scanlines
	s := math.Sin(float64(y) * p.ScanlineFrequency)
	scanFactor := 1.0
	if p.ScanlineIntensity > 0 {
		normSine := (s + 1.0) / 2.0
		scanFactor = 1.0 - (p.ScanlineIntensity * (1.0 - normSine))
	}

	// 3. Noise
	nCol := p.noise.At(x, y)
	ng := color.GrayModel.Convert(nCol).(color.Gray).Y

	noiseFactor := float64(ng) / 255.0

	rf := float64(r8) * scanFactor
	gf := float64(g8) * scanFactor
	bf := float64(b8) * scanFactor

	if p.NoiseIntensity > 0 {
		// Additive noise or blend?
		// Blend towards noise color (gray)
		// Or modulation?
		// "TV Static" is usually additive/subtractive.
		// Let's use interpolation:
		// val = val * (1 - k) + noise * k * 255
		rf = rf*(1.0-p.NoiseIntensity) + noiseFactor*255.0*p.NoiseIntensity
		gf = gf*(1.0-p.NoiseIntensity) + noiseFactor*255.0*p.NoiseIntensity
		bf = bf*(1.0-p.NoiseIntensity) + noiseFactor*255.0*p.NoiseIntensity
	}

	// Clamp
	if rf > 255 { rf = 255 }
	if gf > 255 { gf = 255 }
	if bf > 255 { bf = 255 }

	// Re-apply alpha?
	// If we are modifying RGB of a transparent pixel, we must be careful.
	// We are working with NRGBA, so RGB values are independent of Alpha.
	// However, usually we can't emit light where Alpha is 0.
	// If the user wants the "VHS noise" to appear even in transparent areas,
	// they should put a black background first.
	// Assuming the output should respect the input's alpha mask:
	// We return NRGBA, so we don't premultiply here. The conversion to RGBA will handle it.
	// BUT, if Alpha is 0, RGB is discarded in RGBA.

	return color.NRGBA{
		R: uint8(rf),
		G: uint8(gf),
		B: uint8(bf),
		A: alpha,
	}
}

// NewVHS creates a new VHS pattern.
func NewVHS(img image.Image, ops ...func(any)) image.Image {
	p := &VHS{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Image:             img,
		ScanlineFrequency: math.Pi, // Default to PI (period of 2 pixels)
		ScanlineIntensity: 0.5,
		ColorOffset:       3,
		NoiseIntensity:    0.2,
		Seed:              1,
	}
	// Copy bounds from image if available
	if img != nil {
		p.bounds = img.Bounds()
	}

	for _, op := range ops {
		op(p)
	}

	// Initialize noise after ops, so Seed is set
	p.noise = &HashNoise{Seed: p.Seed}

	return p
}

// Configuration helpers

type hasScanlineFrequency interface {
	SetScanlineFrequency(float64)
}

func SetScanlineFrequency(f float64) func(any) {
	return func(i any) {
		if v, ok := i.(hasScanlineFrequency); ok {
			v.SetScanlineFrequency(f)
		}
	}
}

func (p *VHS) SetScanlineFrequency(f float64) {
	p.ScanlineFrequency = f
}

type hasScanlineIntensity interface {
	SetScanlineIntensity(float64)
}

func SetScanlineIntensity(f float64) func(any) {
	return func(i any) {
		if v, ok := i.(hasScanlineIntensity); ok {
			v.SetScanlineIntensity(f)
		}
	}
}

func (p *VHS) SetScanlineIntensity(f float64) {
	p.ScanlineIntensity = f
}

type hasColorOffset interface {
	SetColorOffset(int)
}

func SetColorOffset(i int) func(any) {
	return func(p any) {
		if v, ok := p.(hasColorOffset); ok {
			v.SetColorOffset(i)
		}
	}
}

func (p *VHS) SetColorOffset(i int) {
	p.ColorOffset = i
}

type hasNoiseIntensity interface {
	SetNoiseIntensity(float64)
}

func SetNoiseIntensity(f float64) func(any) {
	return func(i any) {
		if v, ok := i.(hasNoiseIntensity); ok {
			v.SetNoiseIntensity(f)
		}
	}
}

func (p *VHS) SetNoiseIntensity(f float64) {
	p.NoiseIntensity = f
}

type hasSeed interface {
	SetSeed(int64)
}

func SetSeed(s int64) func(any) {
	return func(i any) {
		if v, ok := i.(hasSeed); ok {
			v.SetSeed(s)
		}
	}
}

func (p *VHS) SetSeed(s int64) {
	p.Seed = s
}
