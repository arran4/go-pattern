package pattern

import (
	"image"
	"image/color"
	"math/rand"
	"sync"
)

// Ensure BlueNoise implements the image.Image interface.
var _ image.Image = (*BlueNoise)(nil)

// BlueNoise generates a blue noise texture approximation.
// We use Mitchell's Best-Candidate Algorithm to generate points,
// then rasterize them with a slight Gaussian blur (conceptually) or return high-frequency noise.
//
// Blue noise is characterized by minimal low-frequency components and no spectral peaks.
//
// For this implementation, we will generate a set of points that are well-separated (Poisson Disk-like)
// and then return a value based on the distance to the nearest point, or simply
// fill pixels.
//
// A better "Blue Noise Mask" for dithering is usually a dense array of values 0-255
// arranged such that thresholding at any level produces a blue noise distribution.
// Generating such a mask (e.g. via Void-and-Cluster) is expensive (O(N^2) or O(N^3)).
//
// We will implement a simplified generator that tries to maintain separation.
type BlueNoise struct {
	Null
	Seed int64
	Values []uint8
	once      sync.Once
}

func (p *BlueNoise) SetSeed(v int64) {
	p.Seed = v
}

func (p *BlueNoise) generate() {
	p.once.Do(func() {
		w, h := p.bounds.Dx(), p.bounds.Dy()
		p.Values = make([]uint8, w*h)

		// Initialize with random noise
		rnd := rand.New(rand.NewSource(p.Seed))
		for i := range p.Values {
			p.Values[i] = uint8(rnd.Intn(256))
		}

		// Apply a simple high-pass filter to remove low frequencies (clumps)
		// Blue noise = White Noise - Low Pass Filtered White Noise.
		// We can simulate this by subtracting a blurred version from the original.

		temp := make([]int, w*h)

		// 3x3 Box Blur (Low Pass)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				sum := 0
				count := 0
				for ky := -1; ky <= 1; ky++ {
					for kx := -1; kx <= 1; kx++ {
						px := x + kx
						py := y + ky
						// Wrap
						if px < 0 { px += w }
						if px >= w { px -= w }
						if py < 0 { py += h }
						if py >= h { py -= h }

						sum += int(p.Values[py*w + px])
						count++
					}
				}
				avg := sum / count
				// High Pass = Original - Low Pass + 128 (to center)
				val := int(p.Values[y*w + x]) - avg + 128
				if val < 0 { val = 0 }
				if val > 255 { val = 255 }
				temp[y*w + x] = val
			}
		}

		for i := range p.Values {
			p.Values[i] = uint8(temp[i])
		}
	})
}

func (p *BlueNoise) At(x, y int) color.Color {
	p.generate()
	w := p.bounds.Dx()
	h := p.bounds.Dy()

	gx := x % w; if gx < 0 { gx += w }
	gy := y % h; if gy < 0 { gy += h }

	v := p.Values[gy*w + gx]
	return color.Gray{Y: v}
}

// NewBlueNoise creates a new BlueNoise pattern.
func NewBlueNoise(ops ...func(any)) image.Image {
	p := &BlueNoise{
		Null: Null{
			bounds: image.Rect(0, 0, 64, 64), // Default size
		},
		Seed: 1, // Default seed
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

