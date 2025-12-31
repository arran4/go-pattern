package pattern

import (
	"image"
	"image/color"
	"testing"
)

func TestRectColor(t *testing.T) {
	bgCol := color.RGBA{0, 60, 20, 255}
	r := NewRect(SetFillColor(bgCol))

	c := r.At(50, 50)
	rn, gn, bn, an := c.RGBA()

	t.Logf("Rect Color: %d %d %d %d", rn, gn, bn, an)

	if gn < 10000 {
		t.Errorf("Rect FillColor not set correctly! Got Green=%d", gn)
	}
}

func TestGenerateCircuit_Output(t *testing.T) {
	rect := image.Rect(0, 0, 100, 100)
	img := GenerateCircuitImpl(rect)

	nonBlack := 0
	bgCount := 0
	traceCount := 0

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			if a == 0 { continue }

			if r == 0 && g == 0 && b == 0 {
				// Black
			} else {
				nonBlack++
			}

			// Background Dark Green: R=0, G=60(high), B=20(mid)
			// 60 -> 15420
			if r < 1000 && g > 10000 && b < 10000 {
				bgCount++
			}
			// Trace Light Green: 100, 200, 100
			if r > 20000 && g > 40000 {
				traceCount++
			}
		}
	}

	if nonBlack == 0 {
		t.Errorf("Generated circuit image is completely black!")
	}
	if bgCount == 0 {
		t.Errorf("No background detected!")
	}
	if traceCount == 0 {
		t.Errorf("No traces detected!")
	}
	t.Logf("Non-Black: %d, BG: %d, Trace: %d", nonBlack, bgCount, traceCount)
}
