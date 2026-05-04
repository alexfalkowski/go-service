package id_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func BenchmarkGenerators(b *testing.B) {
	benchmarks := []struct {
		name string
		kind string
	}{
		{name: "ksuid", kind: "ksuid"},
		{name: "nanoid", kind: "nanoid"},
		{name: "ulid", kind: "ulid"},
		{name: "uuid", kind: "uuid"},
		{name: "xid", kind: "xid"},
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			gen, err := id.NewGenerator(&id.Config{Kind: benchmark.kind}, test.Generators)
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
