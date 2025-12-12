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

	gen = hmac.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))
	key, err = gen.Generate()
	require.Error(t, err)
	require.Empty(t, key)
}

func TestValidSigner(t *testing.T) {
	signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
	require.NoError(t, err)
	require.NotNil(t, signer)

	signer, err = hmac.NewSigner(test.FS, test.NewHMAC())
	require.NoError(t, err)

	e, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.NoError(t, signer.Verify(e, strings.Bytes("test")))

	signer, err = hmac.NewSigner(nil, nil)
	require.NoError(t, err)
	require.Nil(t, signer)
}

func TestInvalidSigner(t *testing.T) {
	signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
	require.NoError(t, err)

	sign, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)

	sign = append(sign, byte('w'))
	require.Error(t, signer.Verify(sign, strings.Bytes("test")))

	signer, err = hmac.NewSigner(test.FS, test.NewHMAC())
	require.NoError(t, err)

	e, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.ErrorIs(t, signer.Verify(e, strings.Bytes("bob")), errors.ErrInvalidMatch)
}
