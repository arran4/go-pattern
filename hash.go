package pattern

// StableHash is a stateless, deterministic hash function based on coordinates and a seed.
// It serves as the standard "random" source for stochastic patterns in this library,
// ensuring that the same coordinates and seed always produce the same output.
//
// The algorithm is a variant of SplitMix64, chosen for its good distribution and speed.
func StableHash(x, y int, seed uint64) uint64 {
	// Mix x, y, and Seed using a robust hash function.
	// We cast x and y to int64 to avoid overflow issues before multiplying,
	// though in Go int is usually 64-bit on modern systems.
	z := uint64(int64(x)*0x9e3779b9+int64(y)*0x632be59b) + seed
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	z = z ^ (z >> 31)
	return z
}
