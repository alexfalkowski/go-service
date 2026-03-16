package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "none"} {
		cfg := test.NewToken(kind)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(test.PEM, ec)
		verifier, _ := ed25519.NewVerifier(test.PEM, ec)
		gen := uuid.NewGenerator()
		tkn := token.NewToken(test.Name, cfg, test.FS, signer, verifier, gen)

		_, err := tkn.Generate("hello", test.UserID.String())
		require.NoError(t, err)
	}
}

func TestVerify(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
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
	}

	for _, kind := range []string{"ssh", "none"} {
		cfg := test.NewToken(kind)
		tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

		gen, err := tkn.Generate(strings.Empty, strings.Empty)
		require.NoError(t, err)

		_, err = tkn.Verify(gen, strings.Empty)
		require.NoError(t, err)
	}
}

func TestNewTokenWithNilConfig(t *testing.T) {
	tkn := token.NewToken(test.Name, nil, test.FS, nil, nil, nil)
	require.Nil(t, tkn)
}

func TestInvalidKindConfig(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		cfg := &token.Config{Kind: kind}
		tkn := token.NewToken(test.Name, cfg, test.FS, nil, nil, nil)

		gen, err := tkn.Generate("hello", test.UserID.String())
		require.Nil(t, gen)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)

		sub, err := tkn.Verify([]byte("test"), "hello")
		require.Equal(t, strings.Empty, sub)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	}
}
