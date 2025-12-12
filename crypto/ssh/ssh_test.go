package ssh_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
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

func TestInvalid(t *testing.T) {
	_, err := ssh.NewSigner(test.FS, &ssh.Config{})
	require.Error(t, err)

	_, err = ssh.NewVerifier(test.FS, &ssh.Config{})
	require.Error(t, err)

	cfg := test.NewSSH("secrets/ssh_public", "secrets/ssh_private")

	signer, err := ssh.NewSigner(test.FS, cfg)
	require.NoError(t, err)

	verifier, err := ssh.NewVerifier(test.FS, cfg)
	require.NoError(t, err)

	sig, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)

	sig = append(sig, byte('w'))
	require.Error(t, verifier.Verify(sig, strings.Bytes("test")))

	cfg = test.NewSSH("secrets/ssh_public", "secrets/ssh_private")

	signer, err = ssh.NewSigner(test.FS, cfg)
	require.NoError(t, err)

	verifier, err = ssh.NewVerifier(test.FS, cfg)
	require.NoError(t, err)

	e, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.ErrorIs(t, verifier.Verify(e, strings.Bytes("bob")), errors.ErrInvalidMatch)

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
