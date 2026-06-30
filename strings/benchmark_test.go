package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

// benchmarkBytesSink keeps conversion results observable so benchmarks report real allocation costs.
var benchmarkBytesSink []byte

// BenchmarkStrings measures the allocation difference between ordinary string-to-byte copies and the unsafe helper.
func BenchmarkStrings(b *testing.B) {
	b.ReportAllocs()

	cp := func(s string) []byte {
		return []byte(s)
	}

	gen := rand.NewGenerator(rand.NewReader())
	s, err := gen.GenerateText(1024)
	require.NoError(b, err)

	b.Run("copy", func(b *testing.B) {
		for b.Loop() {
			benchmarkBytesSink = cp(s)
		}
	})

	b.Run("convert", func(b *testing.B) {
		for b.Loop() {
			benchmarkBytesSink = strings.Bytes(s)
		}
	})
}
