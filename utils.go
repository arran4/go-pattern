package pattern

import (
	"image/color"
)

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func clamp255(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

func clampFloatRange(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func lerpRGBA(a, b color.RGBA, t float64) color.RGBA {
	t = clamp01(t)
	return color.RGBA{
		R: uint8(clamp255(float64(a.R) + (float64(b.R)-float64(a.R))*t)),
		G: uint8(clamp255(float64(a.G) + (float64(b.G)-float64(a.G))*t)),
		B: uint8(clamp255(float64(a.B) + (float64(b.B)-float64(a.B))*t)),
		A: uint8(clamp255(float64(a.A) + (float64(b.A)-float64(a.A))*t)),
	}
}

func posMod(v, m int) int {
	if m == 0 {
		return 0
	}
	r := v % m
	if r < 0 {
		r += m
	}
	return r
}

func jitterColor(c color.RGBA, x, y int, seed int64) color.RGBA {
	noise := float64(StableHash(x, y, uint64(seed)^0xdecafbad)) / float64(^uint64(0))
	// Center noise around 0 and give it a small amplitude.
	delta := (noise - 0.5) * 14
	return color.RGBA{
		R: clampChannel(float64(c.R) + delta),
		G: clampChannel(float64(c.G) + delta*0.8),
		B: clampChannel(float64(c.B) + delta*0.6),
		A: c.A,
	}
}

func clampChannel(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v + 0.5)
}
