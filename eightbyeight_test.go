package pattern

import (
	"image/color"
	"testing"
)

func TestEightByEight_Logic(t *testing.T) {
	// Test Mode 1 (south = {1, 0, 0, 0})
	// idx 0 (x%8 == 0 or 4): south=1 -> sv=2^(4-1)=8.
	// idx 1 (x%8 == 1 or 5): south=0 -> sv=2^4=16.
	// idx 2 (x%8 == 2 or 6): south=0 -> sv=16.
	// idx 3 (x%8 == 3 or 7): south=0 -> sv=16.

	p := NewEightByEight(1, SetLineColor(color.Black), SetSpaceColor(color.White))

	// x=2. xpRaw=2. Mode=1. 1 >= 2 False. -> White.
	if p.At(2, 0) != color.White {
		t.Errorf("Expected White at 2,0 for Mode 1")
	}

	// x=0, y=0. xpRaw=0. Mode=1. 1>=0 True.
	// idx=0. sv=8.
	// dp = (0 - 0) % 8 = 0. 0%8==0. Expect Black.
	if p.At(0, 0) != color.Black {
		t.Errorf("Expected Black at 0,0 for Mode 1")
	}

	// x=0, y=1. dp=(1-0)%8 = 1. 1%8!=0. Expect White.
	if p.At(0, 1) != color.White {
		t.Errorf("Expected White at 0,1 for Mode 1")
	}

	// x=0, y=8. dp=(8-0)%8 = 0. Expect Black.
	if p.At(0, 8) != color.Black {
		t.Errorf("Expected Black at 0,8 for Mode 1")
	}
}

func TestEightByEight_NegativeCoordinates(t *testing.T) {
	// Mode 1.
	p := NewEightByEight(1, SetLineColor(color.Black), SetSpaceColor(color.White))

	// x = -8 (same as x=0).
	// xpRaw = 0. Mode >= 0.
	// idx = 0. sv = 8.
	// dp = (y - 0) % 8.
	// y=0. dp=0. Black.

	if p.At(-8, 0) != color.Black {
		t.Errorf("Expected Black at -8,0")
	}

	// x = -7. xpRaw = -7.
	// Mode >= -7. True.
	// xpAbs = 7. idx = 3. svVal = 0 (default). sv = 16.
	// dp = (0 - (-7)) % 8 = 7 % 8 = 7.
	// 7 % 16 != 0. White.

	if p.At(-7, 0) != color.White {
		t.Errorf("Expected White at -7,0")
	}
}

func TestEightByEight_Palette(t *testing.T) {
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	// Test 2-color palette (replaces default Black/White)
	p := NewEightByEight(1, SetPalette(blue, red)) // Space=Blue(0), Line=Red(1)

	// Same logic as TestLogic, but Red/Blue.
	// x=0, y=0. Expect Line (Red).
	if p.At(0, 0) != red {
		t.Errorf("Expected Red at 0,0 for Palette Mode 1")
	}
	// x=2, y=0. Expect Space (Blue).
	if p.At(2, 0) != blue {
		t.Errorf("Expected Blue at 2,0 for Palette Mode 1")
	}
}

func TestEightByEight_MultiColor(t *testing.T) {
	c0 := color.RGBA{0, 0, 0, 255}
	c1 := color.RGBA{100, 100, 100, 255}
	c2 := color.RGBA{200, 200, 200, 255}

	// Mode = 5.
	// Palette len = 3.
	// bgIdx = 5 % 3 = 2. -> c2
	// fgIdx = (5 / 3) % 3 = 1 % 3 = 1. -> c1

	p := NewEightByEight(5, SetPalette(c0, c1, c2))

	// Mode 5 in binary: 00 00 01 01. (idx 0=1, idx 1=1, idx 2=0, idx 3=0).
	// south values:
	// idx 0 (x%4=0): val=1 -> sv=2^(4-1)=8.
	// idx 1 (x%4=1): val=1 -> sv=8.
	// idx 2 (x%4=2): val=0 -> sv=16.
	// idx 3 (x%4=3): val=0 -> sv=16.

	// Check x=0, y=0.
	// xpRaw=0.
	// useMultiColor=true. Mode cutoff ignored.
	// idx=0 -> sv=8.
	// dp = (0-0)%8 = 0. 0%8 == 0. Matches.
	// Should return fgIdx (c1).
	if p.At(0, 0) != c1 {
		t.Errorf("Expected c1 (fg) at 0,0")
	}

	// Check a pixel that fails the pattern check (sv check).
	// x=0, y=1. dp=(1-0)%8=1. 1%8 != 0.
	// Should return bgIdx (c2).
	if p.At(0, 1) != c2 {
		t.Errorf("Expected c2 (bg) at 0,1")
	}
}
