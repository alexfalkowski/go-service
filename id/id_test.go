package id_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestValidID(t *testing.T) {
	kinds := []string{
		"uuid",
		"ksuid",
		"nanoid",
		"ulid",
		"xid",
	}

	for _, kind := range kinds {
		t.Run(kind, func(t *testing.T) {
			gen, err := id.NewGenerator(&id.Config{Kind: kind}, test.Generators)
			require.NoError(t, err)

			require.NotEmpty(t, gen.Generate())
		})
	}
}

func TestDefaultID(t *testing.T) {
	gen, err := id.NewGenerator(nil, test.Generators)
	require.NoError(t, err)
	require.NotNil(t, gen)
	require.NotEmpty(t, gen.Generate())
}

func TestInvalidID(t *testing.T) {
	_, err := id.NewGenerator(&id.Config{Kind: "invalid"}, test.Generators)
	require.ErrorIs(t, err, id.ErrNotFound)
}
