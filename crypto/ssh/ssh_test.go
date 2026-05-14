package ssh_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
	xssh "golang.org/x/crypto/ssh"
)

func TestGenerator(t *testing.T) {
	gen := ssh.NewGenerator(rand.NewGenerator(rand.NewReader()))
	pub, pri, err := gen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, pub)
	require.NotEmpty(t, pri)

	gen = ssh.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))
	pub, pri, err = gen.Generate()
	require.Error(t, err)
	require.Empty(t, pub)
	require.Empty(t, pri)
}

func TestValid(t *testing.T) {
	cfg := test.NewSSH("secrets/ssh_public", "secrets/ssh_private")

	signer, err := ssh.NewSigner(test.FS, cfg)
	require.NoError(t, err)

	verifier, err := ssh.NewVerifier(test.FS, cfg)
	require.NoError(t, err)

	e, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.NoError(t, verifier.Verify(e, strings.Bytes("test")))

	signer, err = ssh.NewSigner(nil, nil)
	require.NoError(t, err)
	require.Nil(t, signer)

	verifier, err = ssh.NewVerifier(nil, nil)
	require.NoError(t, err)
	require.Nil(t, verifier)
}

func TestInvalidConfig(t *testing.T) {
	_, err := ssh.NewSigner(test.FS, &ssh.Config{})
	require.ErrorIs(t, err, errors.ErrMissingKey)

	_, err = ssh.NewVerifier(test.FS, &ssh.Config{})
	require.ErrorIs(t, err, errors.ErrMissingKey)

	t.Setenv("SSH_EMPTY", "")

	_, err = ssh.NewSigner(test.FS, &ssh.Config{Private: "env:SSH_EMPTY"})
	require.ErrorIs(t, err, errors.ErrMissingKey)

	_, err = ssh.NewVerifier(test.FS, &ssh.Config{Public: "env:SSH_EMPTY"})
	require.ErrorIs(t, err, errors.ErrMissingKey)

	_, err = ssh.NewVerifier(test.FS, &ssh.Config{Public: test.FilePath("secrets/redis")})
	require.Error(t, err)

	_, err = ssh.NewSigner(test.FS, &ssh.Config{Private: test.FilePath("secrets/redis")})
	require.Error(t, err)

	_, err = ssh.NewSigner(
		test.FS,
		&ssh.Config{
			Public:  test.FilePath("secrets/ssh_public"),
			Private: test.FilePath("secrets/none"),
		},
	)
	require.Error(t, err)
}

func TestInvalidSignature(t *testing.T) {
	cfg := test.NewSSH("secrets/ssh_public", "secrets/ssh_private")

	signer, err := ssh.NewSigner(test.FS, cfg)
	require.NoError(t, err)

	verifier, err := ssh.NewVerifier(test.FS, cfg)
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

	block, _ := pem.Decode([]byte(public))
	require.NotNil(t, block)

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	require.NoError(t, err)

	sshPublicKey, err := xssh.NewPublicKey(publicKey)
	require.NoError(t, err)

	var verifierErr error
	require.NotPanics(t, func() {
		_, verifierErr = ssh.NewVerifier(test.FS, &ssh.Config{Public: string(xssh.MarshalAuthorizedKey(sshPublicKey))})
	})
	require.ErrorIs(t, verifierErr, errors.ErrInvalidKeyType)

	var signerErr error
	require.NotPanics(t, func() {
		_, signerErr = ssh.NewSigner(test.FS, &ssh.Config{Private: private})
	})
	require.ErrorIs(t, signerErr, errors.ErrInvalidKeyType)
}
