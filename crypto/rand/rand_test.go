package rand_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestValidRand(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	c, err := gen.GenerateBytes(5)
	require.NoError(t, err)
	require.Len(t, c, 5)

	s, err := gen.GenerateText(32)
	require.NoError(t, err)
	require.Len(t, s, 32)
}

func TestInvalidRand(t *testing.T) {
	gen := rand.NewGenerator(&test.ErrReaderCloser{})
	_, err := gen.GenerateText(5)
	require.Error(t, err)
}
