package pattern

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// DiffusionKernel represents the error diffusion kernel.
// It consists of a list of weights applied to neighboring pixels.
type DiffusionKernel struct {
	Items   []DiffusionItem
	Divisor float64
}

// DiffusionItem represents a single weight in the kernel.
type DiffusionItem struct {
	DX, DY int
	Weight float64
}

// Predefined kernels.
var (
	// FloydSteinberg is the classic error diffusion kernel (7, 3, 5, 1) / 16.
	FloydSteinberg = DiffusionKernel{
		Items: []DiffusionItem{
			{1, 0, 7},
			{-1, 1, 3},
			{0, 1, 5},
			{1, 1, 1},
		},
		Divisor: 16,
	}

	// JarvisJudiceNinke is a larger kernel (5x3) for smoother gradients.
	JarvisJudiceNinke = DiffusionKernel{
		Items: []DiffusionItem{
			{1, 0, 7}, {2, 0, 5},
			{-2, 1, 3}, {-1, 1, 5}, {0, 1, 7}, {1, 1, 5}, {2, 1, 3},
			{-2, 2, 1}, {-1, 2, 3}, {0, 2, 5}, {1, 2, 3}, {2, 2, 1},
		},
		Divisor: 48,
	}

	// Stucki is similar to Jarvis but with different weights for sharper edges.
	Stucki = DiffusionKernel{
		Items: []DiffusionItem{
			{1, 0, 8}, {2, 0, 4},
			{-2, 1, 2}, {-1, 1, 4}, {0, 1, 8}, {1, 1, 4}, {2, 1, 2},
			{-2, 2, 1}, {-1, 2, 2}, {0, 2, 4}, {1, 2, 2}, {2, 2, 1},
		},
		Divisor: 42,
	}

	// Atkinson is a 3x2 kernel designed for small icons.
	Atkinson = DiffusionKernel{
		Items: []DiffusionItem{
			{1, 0, 1}, {2, 0, 1},
			{-1, 1, 1}, {0, 1, 1}, {1, 1, 1},
			{0, 2, 1},
		},
		Divisor: 8,
	}

	// Burkes is an efficient 5x2 kernel.
	Burkes = DiffusionKernel{
		Items: []DiffusionItem{
			{1, 0, 8}, {2, 0, 4},
			{-2, 1, 2}, {-1, 1, 4}, {0, 1, 8}, {1, 1, 4}, {2, 1, 2},
		},
		Divisor: 32,
	}

	// SierraLite is a compact 3x2 kernel.
	SierraLite = DiffusionKernel{
		Items: []DiffusionItem{
			{1, 0, 2},
			{-1, 1, 1}, {0, 1, 1},
		},
		Divisor: 4,
	}

	// Sierra2 (Two-row Sierra) is a balanced 5x2 kernel.
	Sierra2 = DiffusionKernel{
		Items: []DiffusionItem{
			{1, 0, 4}, {2, 0, 3},
			{-2, 1, 1}, {-1, 1, 2}, {0, 1, 3}, {1, 1, 2}, {2, 1, 1},
		},
		Divisor: 16,
	}

	// Sierra3 (Full Sierra) is a 5x3 kernel.
	Sierra3 = DiffusionKernel{
		Items: []DiffusionItem{
			{1, 0, 5}, {2, 0, 3},
			{-2, 1, 2}, {-1, 1, 4}, {0, 1, 5}, {1, 1, 4}, {2, 1, 2},
			{-1, 2, 2}, {0, 2, 3}, {1, 2, 2},
		},
		Divisor: 32,
	}

	// StevensonArce diffuses error in a dispersed pattern.
	StevensonArce = DiffusionKernel{
		Items: []DiffusionItem{
			{2, 0, 32},
			{-3, 1, 12}, {-1, 1, 26}, {1, 1, 30},
			{-1, 2, 16}, {1, 2, 12},
		},
		Divisor: 128, // Sum of weights
	}
)

// ErrorDiffusion applies error diffusion dithering to an image.
type ErrorDiffusion struct {
	Null
	img          image.Image
	kernel       DiffusionKernel
	palette      color.Palette
	result       *image.RGBA
	once         sync.Once
	serpentine   bool
	gammaCorrection float64
	edgeAwareness float64 // 0..1 factor. 1.0 = full blocking of error diffusion across edges.
}

// Serpentine configures the error diffusion to use serpentine scanning.
type Serpentine struct {
	Serpentine bool
}

func (s *Serpentine) SetSerpentine(v bool) {
	s.Serpentine = v
}

type hasSerpentine interface {
	SetSerpentine(bool)
}

// SetSerpentine creates an option to set serpentine scanning.
func SetSerpentine(v bool) func(any) {
	return func(i any) {
		if h, ok := i.(hasSerpentine); ok {
			h.SetSerpentine(v)
		}
	}
}

// Gamma configures gamma correction.
type Gamma struct {
	Gamma float64
}

func (g *Gamma) SetGamma(v float64) {
	g.Gamma = v
}

type hasGamma interface {
	SetGamma(float64)
}

// SetGamma creates an option to set gamma.
func SetGamma(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasGamma); ok {
			h.SetGamma(v)
		}
	}
}

// EdgeAwareness configures the strength of edge-aware diffusion.
type EdgeAwareness struct {
    EdgeAwareness float64
}

func (e *EdgeAwareness) SetEdgeAwareness(v float64) {
    e.EdgeAwareness = v
}

type hasEdgeAwareness interface {
    SetEdgeAwareness(float64)
}

// SetEdgeAwareness creates an option to set edge awareness.
func SetEdgeAwareness(v float64) func(any) {
    return func(i any) {
        if h, ok := i.(hasEdgeAwareness); ok {
            h.SetEdgeAwareness(v)
        }
    }
}

// NewErrorDiffusion creates a new ErrorDiffusion pattern.
// If palette is nil, it defaults to Black and White (1-bit).
// Supports SetSerpentine(bool) and SetGamma(float64) options.
func NewErrorDiffusion(img image.Image, kernel DiffusionKernel, p color.Palette, ops ...func(any)) image.Image {
	if p == nil {
		p = color.Palette{color.Black, color.White}
	}
	b := image.Rect(0, 0, 100, 100)
	if img != nil {
		b = img.Bounds()
	}
	ed := &ErrorDiffusion{
		img:        img,
		kernel:     kernel,
		palette:    p,
		serpentine: true, // Default to true as per best practice and request
		gammaCorrection: 1.0, // Default no gamma
		edgeAwareness: 0.0,
		Null: Null{
			bounds: b,
		},
	}
	for _, op := range ops {
		op(ed)
	}
	return ed
}

func (e *ErrorDiffusion) SetSerpentine(v bool) {
	e.serpentine = v
}

func (e *ErrorDiffusion) SetGamma(v float64) {
	e.gammaCorrection = v
}

func (e *ErrorDiffusion) SetEdgeAwareness(v float64) {
    e.edgeAwareness = v
}

func (e *ErrorDiffusion) At(x, y int) color.Color {
	e.once.Do(e.compute)
	if e.result == nil {
		return color.Black
	}
	return e.result.At(x, y)
}

func (e *ErrorDiffusion) compute() {
	if e.img == nil {
		return
	}
	bounds := e.img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	e.result = image.NewRGBA(image.Rect(0, 0, w, h))

	pixels := make([]float64, w*h*4) // R, G, B, A

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := e.img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			idx := (y*w + x) * 4
			pixels[idx] = float64(r) / 257.0
			pixels[idx+1] = float64(g) / 257.0
			pixels[idx+2] = float64(b) / 257.0
			pixels[idx+3] = float64(a) / 257.0

			// Apply gamma if needed
			if e.gammaCorrection != 1.0 && e.gammaCorrection > 0 {
				pixels[idx] = 255 * math.Pow(pixels[idx]/255, e.gammaCorrection)
				pixels[idx+1] = 255 * math.Pow(pixels[idx+1]/255, e.gammaCorrection)
				pixels[idx+2] = 255 * math.Pow(pixels[idx+2]/255, e.gammaCorrection)
			}
		}
	}

	for y := 0; y < h; y++ {
		// Determine direction
		direction := 1
		startX := 0
		endX := w
		if e.serpentine && (y%2 == 1) {
			direction = -1
			startX = w - 1
			endX = -1
		}

		for x := startX; x != endX; x += direction {
			idx := (y*w + x) * 4

			// Current pixel value
			oldR := pixels[idx]
			oldG := pixels[idx+1]
			oldB := pixels[idx+2]
			oldA := pixels[idx+3]

			c := color.RGBA{
				R: uint8(clamp(oldR)),
				G: uint8(clamp(oldG)),
				B: uint8(clamp(oldB)),
				A: uint8(clamp(oldA)),
			}

			nc := e.palette.Convert(c)
			nr, ng, nb, _ := nc.RGBA()

			newR := float64(nr) / 257.0
			newG := float64(ng) / 257.0
			newB := float64(nb) / 257.0

			// Set result pixel
			e.result.Set(x, y, nc)

			// Calculate error
			errR := oldR - newR
			errG := oldG - newG
			errB := oldB - newB

            // Luminance of current pixel for edge detection
            lum := 0.299*oldR + 0.587*oldG + 0.114*oldB

            // First pass: Calculate total valid weight sum to conserve energy
            totalWeight := 0.0
            type neighbor struct {
                idx int
                weight float64
            }
            var neighbors []neighbor

            for _, item := range e.kernel.Items {
				dx := item.DX * direction
				dy := item.DY
				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < w && ny >= 0 && ny < h {
					nidx := (ny*w + nx) * 4
                    weight := item.Weight

                    if e.edgeAwareness > 0 {
                        nLum := 0.299*pixels[nidx] + 0.587*pixels[nidx+1] + 0.114*pixels[nidx+2]
                        diff := math.Abs(lum - nLum)
                        edgeStrength := diff * 2.0 // Scale diff (0-1) to be more aggressive
                        if edgeStrength > 1 { edgeStrength = 1 }
                        weight = weight * (1.0 - (e.edgeAwareness * edgeStrength))
                    }
                    totalWeight += weight
                    neighbors = append(neighbors, neighbor{nidx, weight})
				}
            }

            // Second pass: Distribute error normalized
            if totalWeight > 0 {
                for _, n := range neighbors {
                    // Normalized weight
                    w := n.weight / totalWeight
                    pixels[n.idx] += errR * w
                    pixels[n.idx+1] += errG * w
                    pixels[n.idx+2] += errB * w
                }
            } else {
                 // If totalWeight is 0 (all edges blocked completely),
                 // we cannot distribute error. It is lost.
                 // This usually only happens if edgeAwareness is 1.0 and all neighbors are very different.
            }
		}
	}
}

func clamp(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
