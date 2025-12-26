package pattern

// Seed configures the seed for a pattern.
type Seed struct {
	Seed int64
}

func (s *Seed) SetSeed(v int64) {
	s.Seed = v
}

type hasSeed interface {
	SetSeed(int64)
}

// SetSeed creates an option to set the seed.
func SetSeed(v int64) func(any) {
	return func(i any) {
		if h, ok := i.(hasSeed); ok {
			h.SetSeed(v)
		}
	}
}
