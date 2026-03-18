package ed25519_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
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

func TestInvalidConfig(t *testing.T) {
	signerTests := []struct {
		config *ed25519.Config
		name   string
	}{
		{name: "empty signer config", config: &ed25519.Config{}},
		{name: "invalid signer private key", config: &ed25519.Config{Public: test.FilePath("secrets/ed25519_public"), Private: test.FilePath("secrets/ed25519_private_invalid")}},
	}

	for _, tt := range signerTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ed25519.NewSigner(test.PEM, tt.config)
			require.Error(t, err)
		})
	}

	verifierTests := []struct {
		config *ed25519.Config
		name   string
	}{
		{name: "empty verifier config", config: &ed25519.Config{}},
		{name: "invalid verifier public key", config: &ed25519.Config{Public: test.FilePath("secrets/ed25519_public_invalid"), Private: test.FilePath("secrets/ed25519_private")}},
	}

	for _, tt := range verifierTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ed25519.NewVerifier(test.PEM, tt.config)
			require.Error(t, err)
		})
	}
}

func TestInvalidSignature(t *testing.T) {
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
}

func TestInvalidKeyType(t *testing.T) {
	public, private, err := rsa.NewGenerator(rand.NewGenerator(rand.NewReader())).Generate()
	require.NoError(t, err)

	publicBlock, _ := pem.Decode([]byte(public))
	require.NotNil(t, publicBlock)

	publicKey, err := x509.ParsePKCS1PublicKey(publicBlock.Bytes)
	require.NoError(t, err)

	marshaledPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	require.NoError(t, err)

	privateBlock, _ := pem.Decode([]byte(private))
	require.NotNil(t, privateBlock)

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	require.NoError(t, err)

	marshaledPrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)

	var verifierErr error
	require.NotPanics(t, func() {
		_, verifierErr = ed25519.NewVerifier(
			test.PEM,
			&ed25519.Config{
				Public:  string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: marshaledPublicKey})),
				Private: test.FilePath("secrets/ed25519_private"),
			},
		)
	})
	require.ErrorIs(t, verifierErr, errors.ErrInvalidKeyType)

	var signerErr error
	require.NotPanics(t, func() {
		_, signerErr = ed25519.NewSigner(
			test.PEM,
			&ed25519.Config{
				Public:  test.FilePath("secrets/ed25519_public"),
				Private: string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey})),
			},
		)
	})
	require.ErrorIs(t, signerErr, errors.ErrInvalidKeyType)
}
