package pattern

// Seed configures the seed for a pattern.
type Seed struct {
	Seed int64
}

func (s *Seed) SetSeed(v int64) {
	s.Seed = v
}

func (s *Seed) SetSeedUint64(v uint64) {
	s.Seed = int64(v)
}

type hasSeed interface {
	SetSeed(int64)
}

type hasSeedUint64 interface {
	SetSeedUint64(uint64)
}

// SetSeed creates an option to set the seed.
func SetSeed(v int64) func(any) {
	return func(i any) {
		if h, ok := i.(hasSeed); ok {
			h.SetSeed(v)
		} else if h, ok := i.(hasSeedUint64); ok {
			h.SetSeedUint64(uint64(v))
		}
	}
}

// WithSeed creates an option to set the seed using uint64.
// This is the preferred option for consistency.
func WithSeed(v uint64) func(any) {
	return func(i any) {
		if h, ok := i.(hasSeedUint64); ok {
			h.SetSeedUint64(v)
		} else if h, ok := i.(hasSeed); ok {
			h.SetSeed(int64(v))
		}
	}
}
