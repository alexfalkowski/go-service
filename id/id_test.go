package id_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestValidID(t *testing.T) {
	configs := []*id.Config{
		{Kind: "uuid"},
		{Kind: "ksuid"},
		{Kind: "nanoid"},
		{Kind: "ulid"},
		{Kind: "xid"},
	}

	for _, config := range configs {
		gen, err := id.NewGenerator(config, test.Generators)
		require.NoError(t, err)

		require.NotEmpty(t, gen.Generate())
	}
}

func TestNilID(t *testing.T) {
	gen, err := id.NewGenerator(nil, test.Generators)
	require.NoError(t, err)
	require.Nil(t, gen)
}

func TestInvalidID(t *testing.T) {
	_, err := id.NewGenerator(&id.Config{Kind: "invalid"}, test.Generators)
	require.Error(t, err)
}
