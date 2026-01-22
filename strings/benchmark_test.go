package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func BenchmarkStrings(b *testing.B) {
	b.ReportAllocs()

	cp := func(s string) []byte {
		return []byte(s)
	}

	gen := rand.NewGenerator(rand.NewReader())
	s, err := gen.GenerateText(1024)
	require.NoError(b, err)

	b.Run("copy", func(b *testing.B) {
		for n := 0; b.Loop(); n++ {
			_ = cp(s)
		}
	})

	b.Run("convert", func(b *testing.B) {
		for n := 0; b.Loop(); n++ {
			_ = strings.Bytes(s)
		}
	})
}
