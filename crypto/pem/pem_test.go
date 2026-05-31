package pem_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	t.Setenv("PEM_EMPTY", "")

	t.Run("valid block", func(t *testing.T) {
		decoded, err := test.PEM.Decode("-----BEGIN TEST KEY-----\nZmlyc3Q=\n-----END TEST KEY-----\n", "TEST KEY")
		require.NoError(t, err)
		require.Equal(t, []byte("first"), decoded)
	})

	t.Run("uses first matching block", func(t *testing.T) {
		decoded, err := test.PEM.Decode(
			"-----BEGIN TEST KEY-----\nZmlyc3Q=\n-----END TEST KEY-----\n"+
				"-----BEGIN OTHER KEY-----\nc2Vjb25k\n-----END OTHER KEY-----\n",
			"TEST KEY",
		)
		require.NoError(t, err)
		require.Equal(t, []byte("first"), decoded)
	})

	t.Run("missing source", func(t *testing.T) {
		_, err := test.PEM.Decode(test.FilePath("none"), "n/a")
		require.Error(t, err)
	})

	t.Run("empty path", func(t *testing.T) {
		_, err := test.PEM.Decode("", "n/a")
		require.ErrorIs(t, err, errors.ErrMissingKey)
	})

	t.Run("empty source", func(t *testing.T) {
		_, err := test.PEM.Decode("env:PEM_EMPTY", "n/a")
		require.ErrorIs(t, err, errors.ErrMissingKey)
	})

	t.Run("invalid block", func(t *testing.T) {
		_, err := test.PEM.Decode(test.FilePath("secrets/redis"), "n/a")
		require.ErrorIs(t, err, pem.ErrInvalidBlock)
	})

	t.Run("invalid kind", func(t *testing.T) {
		_, err := test.PEM.Decode(test.FilePath("secrets/rsa_public"), "what")
		require.ErrorIs(t, err, pem.ErrInvalidKind)
	})
}
