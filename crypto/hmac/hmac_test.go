package hmac_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	gen := hmac.NewGenerator(rand.NewGenerator(rand.NewReader()))
	key, err := gen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, key)
	require.Len(t, key, 32)

	gen = hmac.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))
	key, err = gen.Generate()
	require.Error(t, err)
	require.ErrorContains(t, err, "hmac")
	require.Empty(t, key)
}

func TestIsEnabled(t *testing.T) {
	require.False(t, (*hmac.Config)(nil).IsEnabled())
	require.True(t, (&hmac.Config{}).IsEnabled())
	require.True(t, test.NewHMAC().IsEnabled())
}

func TestValidSigner(t *testing.T) {
	signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
	require.NoError(t, err)
	require.NotNil(t, signer)

	signer, err = hmac.NewSigner(test.FS, test.NewHMAC())
	require.NoError(t, err)

	mac, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.Len(t, mac, hmac.Size)
	require.NoError(t, signer.Verify(mac, strings.Bytes("test")))

	signer, err = hmac.NewSigner(nil, nil)
	require.NoError(t, err)
	require.Nil(t, signer)
}

func TestSignerMissingKey(t *testing.T) {
	t.Setenv("HMAC_EMPTY", "")

	t.Run("missing key", func(t *testing.T) {
		signer, err := hmac.NewSigner(test.FS, &hmac.Config{})
		require.ErrorIs(t, err, errors.ErrMissingKey)
		require.Nil(t, signer)
	})

	t.Run("empty key source", func(t *testing.T) {
		signer, err := hmac.NewSigner(test.FS, &hmac.Config{Key: "env:HMAC_EMPTY"})
		require.ErrorIs(t, err, errors.ErrMissingKey)
		require.Nil(t, signer)
	})

	t.Run("missing key source", func(t *testing.T) {
		signer, err := hmac.NewSigner(test.FS, &hmac.Config{Key: "env:HMAC_MISSING"})
		require.Error(t, err)
		require.ErrorContains(t, err, "env:HMAC_MISSING")
		require.Nil(t, signer)
	})
}

func TestInvalidSigner(t *testing.T) {
	t.Run("tampered signature", func(t *testing.T) {
		signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
		require.NoError(t, err)

		mac, err := signer.Sign(strings.Bytes("test"))
		require.NoError(t, err)

		mac = append(mac, byte('w'))
		require.Error(t, signer.Verify(mac, strings.Bytes("test")))
	})

	t.Run("wrong message", func(t *testing.T) {
		signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
		require.NoError(t, err)

		mac, err := signer.Sign(strings.Bytes("test"))
		require.NoError(t, err)
		require.ErrorIs(t, signer.Verify(mac, strings.Bytes("bob")), errors.ErrInvalidMatch)
	})
}
