package id_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func BenchmarkGenerators(b *testing.B) {
	kinds := []string{
		"ksuid",
		"nanoid",
		"ulid",
		"uuid",
		"xid",
	}

	for _, kind := range kinds {
		b.Run(kind, func(b *testing.B) {
			gen, err := id.NewGenerator(&id.Config{Kind: kind}, test.Generators)
			require.NoError(b, err)

			b.ReportAllocs()

			var value string
			for b.Loop() {
				value = gen.Generate()
			}

			require.NotEmpty(b, value)
		})
	}
}
