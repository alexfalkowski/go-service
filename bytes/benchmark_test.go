package bytes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/stretchr/testify/require"
)

// benchmarkStringSink keeps conversion results observable so benchmarks report real allocation costs.
var benchmarkStringSink string

func BenchmarkBytes(b *testing.B) {
	b.ReportAllocs()

	gen := rand.NewGenerator(rand.NewReader())
	bs, err := gen.GenerateBytes(1024)
	require.NoError(b, err)

	b.Run("copy", func(b *testing.B) {
		for n := 0; b.Loop(); n++ {
			benchmarkStringSink = string(bs)
		}
	})

	b.Run("convert", func(b *testing.B) {
		for n := 0; b.Loop(); n++ {
			benchmarkStringSink = bytes.String(bs)
		}
	})
}
