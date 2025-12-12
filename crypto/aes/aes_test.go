package aes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
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

	gen = aes.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))
	key, err = gen.Generate()
	require.Error(t, err)
	require.Empty(t, key)
}

func TestValidCipher(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())

	cipher, err := aes.NewCipher(rand, test.FS, test.NewAES())
	require.NoError(t, err)
	require.NotNil(t, cipher)

	enc, err := cipher.Encrypt(strings.Bytes("test"))
	require.NoError(t, err)

	d, err := cipher.Decrypt(enc)
	require.NoError(t, err)
	require.Equal(t, "test", bytes.String(d))

	cipher, err = aes.NewCipher(nil, nil, nil)
	require.NoError(t, err)
	require.Nil(t, cipher)
}

func TestInvalidCipher(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	cipher, err := aes.NewCipher(gen, test.FS, &aes.Config{Key: test.FilePath("secrets/aes_invalid")})
	require.NoError(t, err)

	_, err = cipher.Encrypt(strings.Bytes("test"))
	require.Error(t, err)

	_, err = cipher.Decrypt(strings.Bytes("test"))
	require.Error(t, err)

	gen = rand.NewGenerator(&test.ErrReaderCloser{})

	cipher, err = aes.NewCipher(gen, test.FS, test.NewAES())
	require.NoError(t, err)

	_, err = cipher.Encrypt(strings.Bytes("test"))
	require.Error(t, err)

	rand := rand.NewGenerator(rand.NewReader())

	cipher, err = aes.NewCipher(rand, test.FS, test.NewAES())
	require.NoError(t, err)

	enc, err := cipher.Encrypt(strings.Bytes("test"))
	require.NoError(t, err)
	enc = append(enc, byte('w'))

	_, err = cipher.Decrypt(enc)
	require.Error(t, err)

	cipher, err = aes.NewCipher(rand, test.FS, test.NewAES())
	require.NoError(t, err)

	_, err = cipher.Decrypt(strings.Bytes("test"))
	require.Error(t, err)
}
