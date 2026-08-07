package aes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto"
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	gen := aes.NewGenerator(rand.NewGenerator(rand.NewReader()))
	key, err := gen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, key)
	require.Len(t, key, 32)

	gen = aes.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))
	key, err = gen.Generate()
	require.Error(t, err)
	require.Empty(t, key)
}

func TestValidCipher(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	cipher, err := aes.NewCipher(gen, test.FS, test.NewAES())
	require.NoError(t, err)
	require.NotNil(t, cipher)

	enc, err := cipher.Encrypt(crypto.Message{Data: strings.Bytes("test")})
	require.NoError(t, err)

	d, err := cipher.Decrypt(crypto.Message{Data: enc})
	require.NoError(t, err)
	require.Equal(t, "test", bytes.String(d))

	cipher, err = aes.NewCipher(nil, nil, nil)
	require.NoError(t, err)
	require.Nil(t, cipher)
}

func TestInvalidCipherConfig(t *testing.T) {
	t.Setenv("AES_EMPTY", "")

	tests := []struct {
		name        string
		config      *aes.Config
		wantErr     error
		errContains string
	}{
		{name: "missing key", config: &aes.Config{}, wantErr: errors.ErrMissingKey},
		{name: "empty key source", config: &aes.Config{Key: "env:AES_EMPTY"}, wantErr: errors.ErrMissingKey},
		{name: "missing key source", config: &aes.Config{Key: "env:AES_MISSING"}, errContains: "env:AES_MISSING"},
		{name: "invalid key", config: &aes.Config{Key: test.FilePath("secrets/aes_invalid")}, wantErr: errors.ErrInvalidKeySize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipher, err := aes.NewCipher(rand.NewGenerator(rand.NewReader()), test.FS, tt.config)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.ErrorContains(t, err, tt.errContains)
			}
			require.Nil(t, cipher)
		})
	}
}

func TestInvalidCipher(t *testing.T) {
	t.Run("nonce generation error", func(t *testing.T) {
		gen := rand.NewGenerator(&test.ErrReaderCloser{})

		cipher, err := aes.NewCipher(gen, test.FS, test.NewAES())
		require.NoError(t, err)

		_, err = cipher.Encrypt(crypto.Message{Data: strings.Bytes("test")})
		require.Error(t, err)
	})

	t.Run("tampered ciphertext", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		cipher, err := aes.NewCipher(gen, test.FS, test.NewAES())
		require.NoError(t, err)

		enc, err := cipher.Encrypt(crypto.Message{Data: strings.Bytes("test")})
		require.NoError(t, err)
		enc = append(enc, byte('w'))

		_, err = cipher.Decrypt(crypto.Message{Data: enc})
		require.Error(t, err)
	})

	t.Run("short ciphertext", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		cipher, err := aes.NewCipher(gen, test.FS, test.NewAES())
		require.NoError(t, err)

		_, err = cipher.Decrypt(crypto.Message{Data: strings.Bytes("test")})
		require.ErrorIs(t, err, aes.ErrInvalidLength)
	})
}

func TestCipherAuthenticatesMetadata(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	cipher, err := aes.NewCipher(gen, test.FS, test.NewAES())
	require.NoError(t, err)

	enc, err := cipher.Encrypt(crypto.Message{
		Data: strings.Bytes("test"),
		Meta: strings.Bytes("tenant=acme"),
	})
	require.NoError(t, err)

	d, err := cipher.Decrypt(crypto.Message{
		Data: enc,
		Meta: strings.Bytes("tenant=acme"),
	})
	require.NoError(t, err)
	require.Equal(t, "test", bytes.String(d))

	_, err = cipher.Decrypt(crypto.Message{
		Data: enc,
		Meta: strings.Bytes("tenant=other"),
	})
	require.Error(t, err)
}

func TestEncryptUsesRawNonceBytes(t *testing.T) {
	gen := rand.NewGenerator(bytes.NewReader(make([]byte, 12)))

	cipher, err := aes.NewCipher(gen, test.FS, test.NewAES())
	require.NoError(t, err)

	enc, err := cipher.Encrypt(crypto.Message{Data: strings.Bytes("test")})
	require.NoError(t, err)
	require.Equal(t, make([]byte, 12), enc[:12])
}
