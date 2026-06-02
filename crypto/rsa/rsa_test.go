package rsa_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	gen := rsa.NewGenerator(rand.NewGenerator(rand.NewReader()))
	pub, pri, err := gen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, pub)
	require.NotEmpty(t, pri)

	cfg := &rsa.Config{Public: pub, Private: pri}

	publicKey, err := cfg.PublicKey(test.PEM)
	require.NoError(t, err)
	require.Equal(t, rsa.KeySize, publicKey.N.BitLen())

	privateKey, err := cfg.PrivateKey(test.PEM)
	require.NoError(t, err)
	require.Equal(t, rsa.KeySize, privateKey.N.BitLen())
}

func TestValid(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	enc, err := rsa.NewEncryptor(gen, test.PEM, test.NewRSA())
	require.NoError(t, err)
	require.NotNil(t, enc)

	dec, err := rsa.NewDecryptor(gen, test.PEM, test.NewRSA())
	require.NoError(t, err)
	require.NotNil(t, dec)

	cfg := test.NewRSA()

	enc, err = rsa.NewEncryptor(gen, test.PEM, cfg)
	require.NoError(t, err)

	dec, err = rsa.NewDecryptor(gen, test.PEM, cfg)
	require.NoError(t, err)

	ciphertext, err := enc.Encrypt(strings.Bytes("test"))
	require.NoError(t, err)

	plaintext, err := dec.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, "test", bytes.String(plaintext))

	privateKey, err := cfg.PrivateKey(test.PEM)
	require.NoError(t, err)

	plaintext, err = rsa.DecryptOAEP(gen, privateKey, ciphertext)
	require.NoError(t, err)
	require.Equal(t, "test", bytes.String(plaintext))

	publicKey, err := cfg.PublicKey(test.PEM)
	require.NoError(t, err)

	ciphertext, err = rsa.EncryptOAEP(gen, publicKey, strings.Bytes("test"))
	require.NoError(t, err)

	plaintext, err = dec.Decrypt(ciphertext)
	require.NoError(t, err)
	require.Equal(t, "test", bytes.String(plaintext))

	enc, err = rsa.NewEncryptor(gen, test.PEM, nil)
	require.NoError(t, err)
	require.Nil(t, enc)

	dec, err = rsa.NewDecryptor(gen, test.PEM, nil)
	require.NoError(t, err)
	require.Nil(t, dec)
}

func TestInvalidConfig(t *testing.T) {
	t.Setenv("RSA_EMPTY", "")

	t.Run("missing public key", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		enc, err := rsa.NewEncryptor(gen, test.PEM, &rsa.Config{})
		require.ErrorIs(t, err, errors.ErrMissingKey)
		require.Nil(t, enc)
	})

	t.Run("missing private key", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		dec, err := rsa.NewDecryptor(gen, test.PEM, &rsa.Config{})
		require.ErrorIs(t, err, errors.ErrMissingKey)
		require.Nil(t, dec)
	})

	t.Run("empty public key source", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		enc, err := rsa.NewEncryptor(gen, test.PEM, &rsa.Config{Public: "env:RSA_EMPTY"})
		require.ErrorIs(t, err, errors.ErrMissingKey)
		require.Nil(t, enc)
	})

	t.Run("empty private key source", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		dec, err := rsa.NewDecryptor(gen, test.PEM, &rsa.Config{Private: "env:RSA_EMPTY"})
		require.ErrorIs(t, err, errors.ErrMissingKey)
		require.Nil(t, dec)
	})
}

func TestInvalid(t *testing.T) {
	t.Run("tampered ciphertext", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())
		cfg := test.NewRSA()

		enc, err := rsa.NewEncryptor(gen, test.PEM, cfg)
		require.NoError(t, err)

		dec, err := rsa.NewDecryptor(gen, test.PEM, cfg)
		require.NoError(t, err)

		ciphertext, err := enc.Encrypt(strings.Bytes("test"))
		require.NoError(t, err)

		ciphertext = append(ciphertext, byte('w'))
		_, err = dec.Decrypt(ciphertext)
		require.Error(t, err)
	})

	t.Run("short ciphertext", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		dec, err := rsa.NewDecryptor(gen, test.PEM, test.NewRSA())
		require.NoError(t, err)

		_, err = dec.Decrypt(strings.Bytes("test"))
		require.Error(t, err)
	})
}

func TestInvalidKey(t *testing.T) {
	t.Run("invalid public key", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		_, err := rsa.NewEncryptor(gen, test.PEM, &rsa.Config{
			Public:  test.FilePath("secrets/ed25519_public"),
			Private: test.FilePath("secrets/rsa_private"),
		})
		require.Error(t, err)
	})

	t.Run("invalid private key", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())

		_, err := rsa.NewDecryptor(gen, test.PEM, &rsa.Config{
			Public:  test.FilePath("secrets/rsa_public"),
			Private: test.FilePath("secrets/ed25519_private"),
		})
		require.Error(t, err)
	})

	t.Run("small public key", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())
		cfg := &rsa.Config{
			Public:  test.FilePath("secrets/rsa_2048_public"),
			Private: test.FilePath("secrets/rsa_2048_private"),
		}

		_, err := rsa.NewEncryptor(gen, test.PEM, cfg)
		require.ErrorIs(t, err, errors.ErrInvalidKeySize)
	})

	t.Run("small private key", func(t *testing.T) {
		gen := rand.NewGenerator(rand.NewReader())
		cfg := &rsa.Config{
			Public:  test.FilePath("secrets/rsa_2048_public"),
			Private: test.FilePath("secrets/rsa_2048_private"),
		}

		_, err := rsa.NewDecryptor(gen, test.PEM, cfg)
		require.ErrorIs(t, err, errors.ErrInvalidKeySize)
	})
}

func TestInvalidConfigParse(t *testing.T) {
	t.Run("public key", func(t *testing.T) {
		_, err := (&rsa.Config{Public: test.FilePath("secrets/rsa_public_invalid")}).PublicKey(test.PEM)
		require.Error(t, err)
	})

	t.Run("private key", func(t *testing.T) {
		_, err := (&rsa.Config{Private: test.FilePath("secrets/rsa_private_invalid")}).PrivateKey(test.PEM)
		require.Error(t, err)
	})
}
