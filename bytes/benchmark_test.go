package bytes_test

import (
	"crypto/rand"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/stretchr/testify/require"
)

func BenchmarkBytes(b *testing.B) {
	b.ReportAllocs()

	b.Run("copy", func(b *testing.B) {
		for n := 0; b.Loop(); n++ {
			bs := make([]byte, n)
			_, err := rand.Read(bs)
			require.NoError(b, err)
			_ = string(bs)
		}
	})

	b.Run("convert", func(b *testing.B) {
		for n := 0; b.Loop(); n++ {
			bs := make([]byte, n)
			_, err := rand.Read(bs)
			require.NoError(b, err)
			_ = bytes.String(bs)
		}
	})
}
