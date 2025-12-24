package pattern

import (
	"image/color"
	"testing"
)

func TestNewScreenTone(t *testing.T) {
	st := NewScreenTone().(*ScreenTone)
	if st.Radius.Radius != 2 {
		t.Errorf("Expected default radius 2, got %d", st.Radius.Radius)
	}
	if st.Spacing.Spacing != 10 {
		t.Errorf("Expected default spacing 10, got %d", st.Spacing.Spacing)
	}
	if st.Angle.Angle != 45.0 {
		t.Errorf("Expected default angle 45.0, got %f", st.Angle.Angle)
	}
}

func TestScreenTone_Options(t *testing.T) {
	st := NewScreenTone(
		SetRadius(5),
		SetSpacing(20),
		SetAngle(30),
		SetFillColor(color.White),
		SetSpaceColor(color.Black),
	).(*ScreenTone)

	if st.Radius.Radius != 5 {
		t.Errorf("Expected radius 5, got %d", st.Radius.Radius)
	}
	if st.Spacing.Spacing != 20 {
		t.Errorf("Expected spacing 20, got %d", st.Spacing.Spacing)
	}
	if st.Angle.Angle != 30.0 {
		t.Errorf("Expected angle 30.0, got %f", st.Angle.Angle)
	}
	if st.FillColor.FillColor != color.White {
		t.Errorf("Expected FillColor White, got %v", st.FillColor.FillColor)
	}
	if st.SpaceColor.SpaceColor != color.Black {
		t.Errorf("Expected SpaceColor Black, got %v", st.SpaceColor.SpaceColor)
	}
}

func TestScreenTone_At_Angle0(t *testing.T) {
	st := NewScreenTone(
		SetRadius(2),
		SetSpacing(10),
		SetAngle(0),
	).(*ScreenTone)

	// Center of first dot at (5, 5)
	if st.At(5, 5) != color.Black {
		t.Errorf("Expected dot at (5, 5) for Angle 0")
	}
	// Corner should be empty
	if st.At(0, 0) != color.White {
		t.Errorf("Expected space at (0, 0) for Angle 0")
	}
}

func TestScreenTone_At_Angle90(t *testing.T) {
	st := NewScreenTone(
		SetRadius(2),
		SetSpacing(10),
		SetAngle(90),
	).(*ScreenTone)

	// u = y, v = -x
	// At (5, 5): u=5, v=-5.
	// du = 5. dv = -5 mod 10 = 5.
	// Center is 5. distU=0, distV=0. Inside.
	if st.At(5, 5) != color.Black {
		t.Errorf("Expected dot at (5, 5) for Angle 90")
	}
}

func TestScreenTone_At_Angle45(t *testing.T) {
	st := NewScreenTone(
		SetRadius(2),
		SetSpacing(10),
		SetAngle(45),
	).(*ScreenTone)

	// At (0, 0): u=0, v=0. du=0, dv=0. Center=5. Dist=5^2+5^2=50. Radius^2=4. Outside.
	if st.At(0, 0) != color.White {
		t.Errorf("Expected space at (0, 0)")
	}

	// We need to find a point that should be inside.
	// u = (x+y)/sqrt(2), v = (-x+y)/sqrt(2)
	// We want u ~ 5 + 10k, v ~ 5 + 10j
	// Let's try to find x,y for u=5, v=5.
	// x+y = 5*sqrt(2)
	// -x+y = 5*sqrt(2)
	// => 2y = 10*sqrt(2) => y = 5*sqrt(2) approx 7.07
	// => 2x = 0 => x = 0
	// So (0, 7)
	// Let's check (0, 7)
	// u = 7/sqrt(2) = 4.95 -> close to 5.
	// v = 7/sqrt(2) = 4.95 -> close to 5.
	// Should be inside.
	if st.At(0, 7) != color.Black {
		t.Errorf("Expected dot at (0, 7) for Angle 45")
	}
}
