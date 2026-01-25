package pattern

import "testing"

func BenchmarkCryptoNoise_At(b *testing.B) {
	cn := &CryptoNoise{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cn.At(i, i)
	}
}
