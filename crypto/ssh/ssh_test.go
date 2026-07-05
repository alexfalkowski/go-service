package ssh_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
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

	cfg := &ssh.Config{Public: pub, Private: pri}

	signer, err := ssh.NewSigner(test.FS, cfg)
	require.NoError(t, err)

	verifier, err := ssh.NewVerifier(test.FS, cfg)
	require.NoError(t, err)

	sig, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.NoError(t, verifier.Verify(sig, strings.Bytes("test")))

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

	sig, err := signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.NoError(t, verifier.Verify(sig, strings.Bytes("test")))

	signer, err = ssh.NewSigner(nil, nil)
	require.NoError(t, err)
	require.Nil(t, signer)

	verifier, err = ssh.NewVerifier(nil, nil)
	require.NoError(t, err)
	require.Nil(t, verifier)
}

func TestInvalidConfig(t *testing.T) {
	t.Setenv("SSH_EMPTY", "")

	t.Run("missing signer private key", func(t *testing.T) {
		_, err := ssh.NewSigner(test.FS, &ssh.Config{})
		require.ErrorIs(t, err, errors.ErrMissingKey)
	})

	t.Run("missing verifier public key", func(t *testing.T) {
		_, err := ssh.NewVerifier(test.FS, &ssh.Config{})
		require.ErrorIs(t, err, errors.ErrMissingKey)
	})

	t.Run("empty signer private key source", func(t *testing.T) {
		_, err := ssh.NewSigner(test.FS, &ssh.Config{Private: "env:SSH_EMPTY"})
		require.ErrorIs(t, err, errors.ErrMissingKey)
	})

	t.Run("empty verifier public key source", func(t *testing.T) {
		_, err := ssh.NewVerifier(test.FS, &ssh.Config{Public: "env:SSH_EMPTY"})
		require.ErrorIs(t, err, errors.ErrMissingKey)
	})
}

func TestInvalidPublicKeyConfig(t *testing.T) {
	data, err := test.FS.ReadSource(test.FilePath("secrets/ssh_public"))
	require.NoError(t, err)
	public := bytes.String(data)

	t.Run("invalid verifier public key", func(t *testing.T) {
		_, err := ssh.NewVerifier(test.FS, &ssh.Config{Public: test.FilePath("secrets/redis")})
		require.Error(t, err)
	})

	t.Run("verifier public key with comment", func(t *testing.T) {
		_, err := ssh.NewVerifier(test.FS, &ssh.Config{Public: public + " generated@example"})
		require.NoError(t, err)
	})

	t.Run("verifier public key with options", func(t *testing.T) {
		_, err := ssh.NewVerifier(test.FS, &ssh.Config{Public: `from="10.0.0.0/8" ` + public})
		require.ErrorIs(t, err, errors.ErrInvalidKeyFormat)
	})

	t.Run("verifier public key with trailing entry", func(t *testing.T) {
		_, err := ssh.NewVerifier(test.FS, &ssh.Config{Public: public + "\n" + public})
		require.ErrorIs(t, err, errors.ErrInvalidKeyFormat)
	})

	t.Run("verifier security key public key", func(t *testing.T) {
		_, err := ssh.NewVerifier(test.FS, &ssh.Config{
			Public: "sk-ssh-ed25519@openssh.com AAAAGnNrLXNzaC1lZDI1NTE5QG9wZW5zc2guY29tAAAAIJjzc2a20RjCvN/0ibH6UpGuN9F9hDvD7x182bOesNhHAAAABHNzaDo= user@host",
		})
		require.ErrorIs(t, err, errors.ErrInvalidKeyType)
	})
}

func TestInvalidPrivateKeyConfig(t *testing.T) {
	t.Run("invalid signer private key", func(t *testing.T) {
		_, err := ssh.NewSigner(test.FS, &ssh.Config{Private: test.FilePath("secrets/redis")})
		require.Error(t, err)
	})

	t.Run("missing signer private key file", func(t *testing.T) {
		_, err := ssh.NewSigner(
			test.FS,
			&ssh.Config{
				Public:  test.FilePath("secrets/ssh_public"),
				Private: test.FilePath("secrets/none"),
			},
		)
		require.Error(t, err)
	})
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

	sig, err = signer.Sign(strings.Bytes("test"))
	require.NoError(t, err)
	require.ErrorIs(t, verifier.Verify(sig, strings.Bytes("bob")), errors.ErrInvalidMatch)
}

func TestInvalidSignerPrivateKey(t *testing.T) {
	tests := []struct {
		signer *ssh.Signer
		name   string
	}{
		{name: "nil signer", signer: nil},
		{name: "zero value signer", signer: &ssh.Signer{}},
		{name: "short private key", signer: &ssh.Signer{PrivateKey: []byte("short")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				sig []byte
				err error
			)
			require.NotPanics(t, func() {
				sig, err = tt.signer.Sign(strings.Bytes("test"))
			})
			require.Nil(t, sig)
			require.ErrorIs(t, err, errors.ErrInvalidKeySize)
		})
	}
}

func TestInvalidVerifierPublicKey(t *testing.T) {
	tests := []struct {
		verifier *ssh.Verifier
		name     string
	}{
		{name: "nil verifier", verifier: nil},
		{name: "zero value verifier", verifier: &ssh.Verifier{}},
		{name: "short public key", verifier: &ssh.Verifier{PublicKey: []byte("short")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			require.NotPanics(t, func() {
				err = tt.verifier.Verify(strings.Bytes("sig"), strings.Bytes("test"))
			})
			require.ErrorIs(t, err, errors.ErrInvalidKeySize)
		})
	}
}

func TestInvalidKeyType(t *testing.T) {
	var verifierErr error
	require.NotPanics(t, func() {
		_, verifierErr = ssh.NewVerifier(test.FS, &ssh.Config{
			Public: test.FilePath("secrets/rsa_ssh_public"),
		})
	})
	require.ErrorIs(t, verifierErr, errors.ErrInvalidKeyType)

	var signerErr error
	require.NotPanics(t, func() {
		_, signerErr = ssh.NewSigner(test.FS, &ssh.Config{Private: test.FilePath("secrets/rsa_private")})
	})
	require.ErrorIs(t, signerErr, errors.ErrInvalidKeyType)
}
