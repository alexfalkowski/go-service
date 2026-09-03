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

func TestGeneratorGeneratesKeyOrReturnsEntropyError(t *testing.T) {
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

func TestConfigIsEnabledUnlessNil(t *testing.T) {
	require.False(t, (*hmac.Config)(nil).IsEnabled())
	require.True(t, (&hmac.Config{}).IsEnabled())
	require.True(t, test.NewHMAC().IsEnabled())
}

func TestNewSignerSignsAndVerifiesOrDisablesWhenUnconfigured(t *testing.T) {
	signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
	require.NoError(t, err)
	require.NotNil(t, signer)

	mac, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.Len(t, mac, hmac.Size)
	require.NoError(t, signer.Verify(mac, strings.Bytes("test")))

	signer, err = hmac.NewSigner(nil, nil)
	require.NoError(t, err)
	require.Nil(t, signer)
}

func TestNewSignerRejectsMissingKey(t *testing.T) {
	t.Setenv("HMAC_EMPTY", "")

	tests := []struct {
		name        string
		config      *hmac.Config
		wantErr     error
		errContains string
	}{
		{name: "missing key", config: &hmac.Config{}, wantErr: errors.ErrMissingKey},
		{name: "empty key source", config: &hmac.Config{Key: "env:HMAC_EMPTY"}, wantErr: errors.ErrMissingKey},
		{name: "missing key source", config: &hmac.Config{Key: "env:HMAC_MISSING"}, errContains: "env:HMAC_MISSING"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signer, err := hmac.NewSigner(test.FS, tt.config)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.ErrorContains(t, err, tt.errContains)
			}
			require.Nil(t, signer)
		})
	}
}

func TestSignerRejectsInvalidSignature(t *testing.T) {
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
