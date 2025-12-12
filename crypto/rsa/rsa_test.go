package rsa_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
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

	gen = rsa.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))
	pub, pri, err = gen.Generate()
	require.Error(t, err)
	require.Empty(t, pub)
	require.Empty(t, pri)
}

func TestValid(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())

	enc, err := rsa.NewEncryptor(rand, test.PEM, test.NewRSA())
	require.NoError(t, err)
	require.NotNil(t, enc)

	dec, err := rsa.NewDecryptor(rand, test.PEM, test.NewRSA())
	require.NoError(t, err)
	require.NotNil(t, dec)

	cfg := test.NewRSA()

	enc, err = rsa.NewEncryptor(rand, test.PEM, cfg)
	require.NoError(t, err)

	dec, err = rsa.NewDecryptor(rand, test.PEM, cfg)
	require.NoError(t, err)

	e, err := enc.Encrypt(strings.Bytes("test"))
	require.NoError(t, err)

	d, err := dec.Decrypt(e)
	require.NoError(t, err)
	require.Equal(t, "test", bytes.String(d))

	enc, err = rsa.NewEncryptor(rand, test.PEM, nil)
	require.NoError(t, err)
	require.Nil(t, enc)

	dec, err = rsa.NewDecryptor(rand, test.PEM, nil)
	require.NoError(t, err)
	require.Nil(t, dec)
}

func TestInvalid(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())

	enc, err := rsa.NewEncryptor(rand, test.PEM, &rsa.Config{})
	require.Error(t, err)
	require.Nil(t, enc)

	dec, err := rsa.NewDecryptor(rand, test.PEM, &rsa.Config{})
	require.Error(t, err)
	require.Nil(t, dec)

	cfg := test.NewRSA()

	enc, err = rsa.NewEncryptor(rand, test.PEM, cfg)
	require.NoError(t, err)

	dec, err = rsa.NewDecryptor(rand, test.PEM, cfg)
	require.NoError(t, err)

	e, err := enc.Encrypt(strings.Bytes("test"))
	require.NoError(t, err)

	e = append(e, byte('w'))
	_, err = dec.Decrypt(e)
	require.Error(t, err)

	_, err = dec.Decrypt(strings.Bytes("test"))
	require.Error(t, err)

	_, err = rsa.NewEncryptor(rand, test.PEM, &rsa.Config{
		Public:  test.FilePath("secrets/ed25519_public"),
		Private: test.FilePath("secrets/rsa_private"),
	})
	require.Error(t, err)

	_, err = rsa.NewDecryptor(rand, test.PEM, &rsa.Config{
		Public:  test.FilePath("secrets/rsa_public"),
		Private: test.FilePath("secrets/ed25519_private"),
	})
	require.Error(t, err)
}
