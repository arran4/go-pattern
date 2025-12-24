package pattern

import (
	"crypto/rand"
	"image"
	"image/color"
)

// Ensure Noise implements the image.Image interface.
var _ image.Image = (*Noise)(nil)

// NoiseAlgorithm defines the source of randomness for the Noise pattern.
type NoiseAlgorithm interface {
	At(x, y int) color.Color
}

// Noise pattern generates random noise based on a selected algorithm.
type Noise struct {
	Null
	algo NoiseAlgorithm
}

func (n *Noise) At(x, y int) color.Color {
	if n.algo != nil {
		return n.algo.At(x, y)
	}
	return color.RGBA{0, 0, 0, 255}
}

// NewNoise creates a new Noise pattern.
func NewNoise(ops ...func(any)) image.Image {
	p := &Noise{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		algo: &CryptoNoise{}, // Default to crypto noise
	}
	for _, op := range ops {
		op(p)
	}
	return p
}

// NoiseAlgorithm configuration
type hasNoiseAlgorithm interface {
	SetNoiseAlgorithm(NoiseAlgorithm)
}

func (n *Noise) SetNoiseAlgorithm(algo NoiseAlgorithm) {
	n.algo = algo
}

func SetNoiseAlgorithm(algo NoiseAlgorithm) func(any) {
	return func(i any) {
		if n, ok := i.(hasNoiseAlgorithm); ok {
			n.SetNoiseAlgorithm(algo)
		}
	}
}

// --- Algorithms ---

// CryptoNoise uses crypto/rand.
type CryptoNoise struct{}

func (c *CryptoNoise) At(x, y int) color.Color {
	b := make([]byte, 1)
	rand.Read(b)
	return color.Gray{Y: b[0]}
}

// UniformNoise uses a simple deterministic pseudo-random number generator based on coordinates.
type UniformNoise struct {
	Seed int64
}

func (u *UniformNoise) At(x, y int) color.Color {
	h := hash(x, y) + u.Seed
	r := simpleRand(h)
	return color.Gray{Y: r}
}

func hash(x, y int) int64 {
	// Simple mixing
	return int64(x)*48271 + int64(y)*69621
}

func simpleRand(seed int64) uint8 {
	// Linear Congruential Generator
	seed = seed * 6364136223846793005 + 1442695040888963407
	return uint8(seed >> 56)
}

// PiNoise uses digits of Pi.
type PiNoise struct {
	Stride int
}

// 1000 digits of Pi
var piDigits = "3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067982148086513282306647093844609550582231725359408128481117450284102701938521105559644622948954930381964428810975665933446128475648233786783165271201909145648566923460348610454326648213393607260249141273724587006606315588174881520920962829254091715364367892590360011330530548820466521384146951941511609433057270365759591953092186117381932611793105118548074462379962749567351885752724891227938183011949129833673362440656643086021394946395224737190702179860943702770539217176293176752384674818467669405132000568127145263560827785771342757789609173637178721468440901224953430146549585371050792279689258923542019956112129021960864034418159813629774771309960518707211349999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989"

func (p *PiNoise) At(x, y int) color.Color {
	stride := p.Stride
	if stride == 0 {
		stride = 255
	}
	// Map 2D to 1D
	idx := (x + y*stride) % len(piDigits)
	if idx < 0 {
		idx += len(piDigits)
	}
	digit := piDigits[idx]
	// '0' -> 0, '9' -> 252
	val := (digit - '0') * 28
	return color.Gray{Y: val}
}
