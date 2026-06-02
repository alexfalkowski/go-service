package bcrypt_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestSigner(t *testing.T) {
	signer := bcrypt.NewSigner()

	t.Run("sign and verify", func(t *testing.T) {
		hash, err := signer.Sign(strings.Bytes("test"))
		require.NoError(t, err)
		require.NotEmpty(t, hash)

		cost, err := bcrypt.Cost(hash)
		require.NoError(t, err)
		require.Equal(t, bcrypt.DefaultCost, cost)

		require.NoError(t, signer.Verify(hash, strings.Bytes("test")))
	})

	t.Run("wrong password", func(t *testing.T) {
		hash, err := signer.Sign(strings.Bytes("steve"))
		require.NoError(t, err)
		require.ErrorIs(t, signer.Verify(hash, strings.Bytes("bob")), errors.ErrInvalidMatch)
	})

	t.Run("malformed hash", func(t *testing.T) {
		require.ErrorIs(t, signer.Verify(strings.Bytes("steve"), strings.Bytes("bob")), errors.ErrInvalidMatch)
	})
}
