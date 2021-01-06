package internal

import (
	"testing"
	"time"
)

func BenchmarkNow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Now()
	}
}
func BenchmarkTimeNow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}
func BenchmarkNowUnixNano(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Now().UnixNano()
	}
}

func BenchmarkTimeNowUnixNano(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Now().UnixNano()
	}
}

func TestTimeNow(t *testing.T) {
	Now().UnixNano()
}
