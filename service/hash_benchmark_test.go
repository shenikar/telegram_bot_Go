package service

import (
	"testing"
)

func BenchmarkSingle(b *testing.B) {
	hash := "962012d09b8170d912f0669f6d7d9d07"
	service := NewHashService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetWord(hash)
	}
}

func BenchmarkMulti(b *testing.B) {
	hash := "962012d09b8170d912f0669f6d7d9d07"
	service := NewHashService()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetWordMulti(hash)
	}
}
