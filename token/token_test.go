package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			ec := test.NewEd25519()
			signer, _ := ed25519.NewSigner(test.PEM, ec)
			verifier, _ := ed25519.NewVerifier(test.PEM, ec)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

			_, err := tkn.Generate("hello", test.UserID.String())
			require.NoError(t, err)
		})
	}
}

func TestVerify(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			ec := test.NewEd25519()
			signer, _ := ed25519.NewSigner(test.PEM, ec)
			verifier, _ := ed25519.NewVerifier(test.PEM, ec)
			gen := uuid.NewGenerator()
			tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

			token, err := tkn.Generate("hello", test.UserID.String())
			require.NoError(t, err)

			sub, err := tkn.Verify(token, "hello")
			require.NoError(t, err)
			require.Equal(t, test.UserID.String(), sub)

			sub, err = tkn.Verify(token, "other")
			require.Equal(t, strings.Empty, sub)
			require.Error(t, err)
			require.NotErrorIs(t, err, errors.ErrInvalidMatch)
		})
	}

	t.Run("ssh", func(t *testing.T) {
		cfg := test.NewToken("ssh")
		tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

		gen, err := tkn.Generate("hello", strings.Empty)
		require.NoError(t, err)

		sshToken := ssh.NewToken(cfg.SSH, test.FS)
		sub, err := sshToken.Verify(bytes.String(gen), "hello")
		require.NoError(t, err)
		require.Equal(t, test.UserID.String(), sub)

		sub, err = tkn.Verify(gen, "hello")
		require.NoError(t, err)
		require.Equal(t, test.UserID.String(), sub)

		_, err = tkn.Verify(gen, "other")
		require.ErrorIs(t, err, errors.ErrInvalidAudience)
	})
}

func TestVerifyDispatchesJWTToConcreteImplementation(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

	raw, err := tkn.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	jwtToken := jwt.NewToken(cfg.JWT, nil, verifier, gen)
	sub, err := jwtToken.Verify(bytes.String(raw), "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestVerifyDispatchesPasetoToConcreteImplementation(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

	raw, err := tkn.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	pasetoToken := paseto.NewToken(cfg.Paseto, nil, verifier, gen)
	sub, err := pasetoToken.Verify(bytes.String(raw), "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestVerifyDispatchesSSHToConcreteImplementation(t *testing.T) {
	cfg := test.NewToken("ssh")
	tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

	raw, err := tkn.Generate("hello", strings.Empty)
	require.NoError(t, err)

	sshToken := ssh.NewToken(cfg.SSH, test.FS)
	sub, err := sshToken.Verify(bytes.String(raw), "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestVerifyRejectsPasetoEmptySubject(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

	raw, err := tkn.Generate("hello", strings.Empty)
	require.NoError(t, err)

	sub, err := tkn.Verify(raw, "hello")
	require.Equal(t, strings.Empty, sub)
	require.ErrorIs(t, err, errors.ErrInvalidSubject)
}

func TestUnknownKindConfig(t *testing.T) {
	cfg := test.NewToken("none")
	tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

	gen, err := tkn.Generate("hello", test.UserID.String())
	require.Nil(t, gen)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)

	sub, err := tkn.Verify([]byte("test"), "hello")
	require.Equal(t, strings.Empty, sub)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestNewTokenWithNilConfig(t *testing.T) {
	tkn := token.NewToken(test.Name, nil, test.FS, nil, nil, nil)
	require.Nil(t, tkn)
}

func TestInvalidKindConfig(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		t.Run(kind, func(t *testing.T) {
			cfg := &token.Config{Kind: kind}
			tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

			gen, err := tkn.Generate("hello", test.UserID.String())
			require.Nil(t, gen)
			require.ErrorIs(t, err, errors.ErrInvalidConfig)

			sub, err := tkn.Verify([]byte("test"), "hello")
			require.Equal(t, strings.Empty, sub)
			require.ErrorIs(t, err, errors.ErrInvalidConfig)
		})
	}
}

func TestInvalidDependenciesDoNotPanic(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

			gen, err := tkn.Generate("hello", test.UserID.String())
			require.Nil(t, gen)
			require.ErrorIs(t, err, errors.ErrInvalidConfig)

			sub, err := tkn.Verify([]byte("test"), "hello")
			require.Equal(t, strings.Empty, sub)
			require.ErrorIs(t, err, errors.ErrInvalidConfig)
		})
	}
}

func TestInvalidMatchClassification(t *testing.T) {
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()

	for _, kind := range []string{"jwt", "paseto"} {
		t.Run(kind, func(t *testing.T) {
			cfg := test.NewToken(kind)
			tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

			sub, err := tkn.Verify([]byte("test"), "hello")
			require.Equal(t, strings.Empty, sub)
			require.ErrorIs(t, err, errors.ErrInvalidMatch)
		})
	}

	t.Run("ssh", func(t *testing.T) {
		cfg := test.NewToken("ssh")
		tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

		sub, err := tkn.Verify([]byte("test"), "hello")
		require.Equal(t, strings.Empty, sub)
		require.ErrorIs(t, err, errors.ErrInvalidMatch)
	})
}
