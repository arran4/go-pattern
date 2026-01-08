package pattern

// FalloffCurve configures an exponent used for radial falloff calculations.
type FalloffCurve struct {
	FalloffCurve float64
}

func (f *FalloffCurve) SetFalloffCurve(v float64) {
	f.FalloffCurve = v
}

type hasFalloffCurve interface {
	SetFalloffCurve(float64)
}

// SetFalloffCurve creates an option to set the radial falloff curve strength.
func SetFalloffCurve(v float64) func(any) {
	return func(i any) {
		if h, ok := i.(hasFalloffCurve); ok {
			h.SetFalloffCurve(v)
		}
	}
}
