package pattern

import (
	"image/color"
	"testing"
)

func TestEdgeDetect_At(t *testing.T) {
	// Create a simple image with a vertical edge
	// Left half black, right half white.
	// 4x4 image.
	// 0,0 1,0 2,0 3,0
	// B   B   W   W

	// Coordinates:
	// x=0,1 are Black. x=2,3 are White.

	// Explicitly set SpaceSize to 1 to match legacy 1x1 checker behavior expected by this test.
	src := NewChecker(color.Black, color.White, SetSpaceSize(1))
	// Checker is 1x1.
	// (0,0): 0==0 -> Black
	// (1,0): 1!=0 -> White
	// (2,0): 0==0 -> Black
	// So Checker is 1 pixel wide alternating. This is very high frequency.

	// NewSimpleZoom(src, 2).
	// Pixels 0,0 and 1,0 map to checker 0,0 (Black).
	// Pixels 2,0 and 3,0 map to checker 1,0 (White).

	zm := NewSimpleZoom(src, 2)
	ed := NewEdgeDetect(zm)

	// At (0,0): Neighbors include (-1,-1)...(1,1).
	// (-1,-1) maps to (0,0) of zoom -> (0,0) checker -> Black.
	// (1,1) maps to (0,0) zoom -> Black.
	// So (0,0) is surrounded by Black. Should be 0 magnitude.

	c00 := ed.At(0, 0)
	g00 := color.GrayModel.Convert(c00).(color.Gray)
	if g00.Y != 0 {
		t.Errorf("Expected black at 0,0 (flat region), got %d", g00.Y)
	}

	// At (1,0):
	// x=1. neighbors x=0,1,2.
	// x=0 -> Black
	// x=1 -> Black
	// x=2 -> White
	// So we have an edge transition on the right side.
	// Gx calculation:
	// grid:
	// B B W
	// B B W
	// B B W
	// (assuming y neighbors are same)

	// grid values (Luminance):
	// 0 0 1
	// 0 0 1
	// 0 0 1

	// Gx:
	// -1(0) + 1(1) = 1
	// -2(0) + 2(1) = 2
	// -1(0) + 1(1) = 1
	// Sum = 4.

	// Gy:
	// -1(0) -2(0) -1(1) = -1
	// 1(0) + 2(0) + 1(1) = 1
	// Sum = 0.

	// Mag = sqrt(16) = 4.
	// We clamp at 1.0 (value 255).

	c10 := ed.At(1, 0)
	g10 := color.GrayModel.Convert(c10).(color.Gray)
	if g10.Y < 200 {
		t.Errorf("Expected bright edge at 1,0, got %d", g10.Y)
	}
}
