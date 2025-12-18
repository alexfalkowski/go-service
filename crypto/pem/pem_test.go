package pem_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	_, err := test.PEM.Decode(test.FilePath("none"), "n/a")
	require.Error(t, err)

	_, err = test.PEM.Decode(test.FilePath("secrets/redis"), "n/a")
	require.ErrorIs(t, err, pem.ErrInvalidBlock)

	_, err = test.PEM.Decode(test.FilePath("secrets/rsa_public"), "what")
	require.ErrorIs(t, err, pem.ErrInvalidKind)
}
