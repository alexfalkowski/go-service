package ed25519_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	gen := ed25519.NewGenerator(rand.NewGenerator(rand.NewReader()))
	pub, pri, err := gen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, pub)
	require.NotEmpty(t, pri)

	gen = ed25519.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))
	pub, pri, err = gen.Generate()
	require.Error(t, err)
	require.Empty(t, pub)
	require.Empty(t, pri)
}

func TestValid(t *testing.T) {
	cfg := test.NewEd25519()

	signer, err := ed25519.NewSigner(test.PEM, cfg)
	require.NoError(t, err)
	require.NotNil(t, signer.PrivateKey)

	verifier, err := ed25519.NewVerifier(test.PEM, cfg)
	require.NoError(t, err)
	require.NotNil(t, verifier.PublicKey)

	e, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.NoError(t, verifier.Verify(e, strings.Bytes("test")))

	signer, err = ed25519.NewSigner(nil, nil)
	require.NoError(t, err)
	require.Nil(t, signer)

	verifier, err = ed25519.NewVerifier(nil, nil)
	require.NoError(t, err)
	require.Nil(t, verifier)
}

func TestInvalid(t *testing.T) {
	configs := []*ed25519.Config{
		{},
		{Public: test.FilePath("secrets/ed25519_public"), Private: test.FilePath("secrets/ed25519_private_invalid")},
	}

	for _, config := range configs {
		_, err := ed25519.NewSigner(test.PEM, config)
		require.Error(t, err)
	}

	configs = []*ed25519.Config{
		{},
		{Public: test.FilePath("secrets/ed25519_public_invalid"), Private: test.FilePath("secrets/ed25519_private")},
	}

	for _, config := range configs {
		_, err := ed25519.NewVerifier(test.PEM, config)
		require.Error(t, err)
	}

	cfg := test.NewEd25519()

	signer, err := ed25519.NewSigner(test.PEM, cfg)
	require.NoError(t, err)

	verifier, err := ed25519.NewVerifier(test.PEM, cfg)
	require.NoError(t, err)

	sig, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)

	sig = append(sig, byte('w'))
	require.Error(t, verifier.Verify(sig, strings.Bytes("test")))

	e, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.ErrorIs(t, verifier.Verify(e, strings.Bytes("bob")), errors.ErrInvalidMatch)

	_, err = ed25519.NewVerifier(
		test.PEM,
		&ed25519.Config{
			Public:  test.FilePath("secrets/rsa_public"),
			Private: test.FilePath("secrets/ed25519_private"),
		},
	)
	require.Error(t, err)

	_, err = ed25519.NewSigner(
		test.PEM,
		&ed25519.Config{
			Public:  test.FilePath("secrets/ed25519_public"),
			Private: test.FilePath("secrets/rsa_private"),
		},
	)
	require.Error(t, err)
}
