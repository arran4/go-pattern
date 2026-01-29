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
