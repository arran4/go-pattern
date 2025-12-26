package pattern

import (
	"image"
	"image/color"
	"math/rand"
	"sync"
)

// Ensure Plasma implements the image.Image interface.
var _ image.Image = (*Plasma)(nil)

// Plasma generates a plasma noise texture using Diamond-Square algorithm.
// It supports RGB (independent channels) or Grayscale.
type Plasma struct {
	Null
	Seed       int64
	Roughness  float64
	Color      bool // If true, generates RGB plasma. If false, grayscale.
	gridR      [][]float64
	gridG      [][]float64
	gridB      [][]float64
	once       sync.Once
}

func (p *Plasma) SetSeed(v int64) {
	p.Seed = v
}

func (p *Plasma) generate() {
	p.once.Do(func() {
		w := p.bounds.Dx()
		h := p.bounds.Dy()
		size := w
		if h > size {
			size = h
		}

		n := 1
		for n < size {
			n *= 2
		}
		gridSize := n + 1

		rnd := rand.New(rand.NewSource(p.Seed))

		// Generate R (or Gray)
		p.gridR = p.generateGrid(gridSize, n, rnd)

		if p.Color {
			// Seed offset for other channels
			rndG := rand.New(rand.NewSource(p.Seed + 1))
			p.gridG = p.generateGrid(gridSize, n, rndG)

			rndB := rand.New(rand.NewSource(p.Seed + 2))
			p.gridB = p.generateGrid(gridSize, n, rndB)
		}
	})
}

func (p *Plasma) generateGrid(gridSize, n int, rnd *rand.Rand) [][]float64 {
	grid := make([][]float64, gridSize)
	for i := range grid {
		grid[i] = make([]float64, gridSize)
	}

	// Initialize corners
	grid[0][0] = rnd.Float64()
	grid[0][n] = rnd.Float64()
	grid[n][0] = rnd.Float64()
	grid[n][n] = rnd.Float64()

	p.divide(grid, 0, 0, n, p.Roughness, rnd)
	return grid
}

func (p *Plasma) divide(grid [][]float64, x, y, size int, stdDev float64, rnd *rand.Rand) {
	half := size / 2
	if half < 1 {
		return
	}

	scale := stdDev * float64(size) // Tuning factor

	// Square
	avg := (grid[x][y] + grid[x+size][y] + grid[x][y+size] + grid[x+size][y+size]) / 4.0
	grid[x+half][y+half] = avg + (rnd.Float64()*2 - 1) * scale

	c := grid[x+half][y+half]

	// Diamond
	grid[x+half][y] = (grid[x][y] + grid[x+size][y] + c) / 3.0 + (rnd.Float64()*2 - 1) * scale
	grid[x+half][y+size] = (grid[x][y+size] + grid[x+size][y+size] + c) / 3.0 + (rnd.Float64()*2 - 1) * scale
	grid[x][y+half] = (grid[x][y] + grid[x][y+size] + c) / 3.0 + (rnd.Float64()*2 - 1) * scale
	grid[x+size][y+half] = (grid[x+size][y] + grid[x+size][y+size] + c) / 3.0 + (rnd.Float64()*2 - 1) * scale

	newRoughness := stdDev / 2.0
	p.divide(grid, x, y, half, newRoughness, rnd)
	p.divide(grid, x+half, y, half, newRoughness, rnd)
	p.divide(grid, x, y+half, half, newRoughness, rnd)
	p.divide(grid, x+half, y+half, half, newRoughness, rnd)
}


func (p *Plasma) At(x, y int) color.Color {
	p.generate()

	gw := len(p.gridR) - 1
	gh := len(p.gridR[0]) - 1

	gx := x % gw
	if gx < 0 { gx += gw }
	gy := y % gh
	if gy < 0 { gy += gh }

	r := clampFloat(p.gridR[gx][gy])

	if p.Color {
		g := clampFloat(p.gridG[gx][gy])
		b := clampFloat(p.gridB[gx][gy])
		return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
	}

	v := uint8(r * 255)
	return color.RGBA{v, v, v, 255}
}

func clampFloat(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}


// NewPlasma creates a new Plasma pattern.
func NewPlasma(ops ...func(any)) image.Image {
	p := &Plasma{
		Null: Null{
			bounds: image.Rect(0, 0, 256, 256),
		},
		Seed: 1,
		Roughness: 1.0,
		Color: true, // Default to Color as per user request "full 0-255 per channel"
	}
	for _, op := range ops {
		op(p)
	}
	return p
}
