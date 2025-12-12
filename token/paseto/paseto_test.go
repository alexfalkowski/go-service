package paseto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	paseto := paseto.NewToken(cfg.Paseto, signer, verifier, uuid.NewGenerator())

	token, err := paseto.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sub, err := paseto.Verify(token, "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestInvalid(t *testing.T) {
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()

	cfg := test.NewToken("paseto")
	token := paseto.NewToken(cfg.Paseto, signer, verifier, gen)

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "test")
	require.Error(t, err)

	token = paseto.NewToken(&paseto.Config{Issuer: "test", Expiration: "1h"}, signer, verifier, gen)

	tkn, err = token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "hello")
	require.NoError(t, err)

	for _, tkn := range []string{"invalid"} {
		cfg := test.NewToken("paseto")
		token := paseto.NewToken(cfg.Paseto, signer, verifier, gen)

		_, err := token.Verify(tkn, "aud")
		require.Error(t, err)
	}

	cfg = test.NewToken("paseto")

	token = paseto.NewToken(cfg.Paseto, &ed25519.Signer{}, verifier, gen)
	_, err = token.Generate("hello", test.UserID.String())
	require.Error(t, err)

	token = paseto.NewToken(cfg.Paseto, signer, &ed25519.Verifier{}, gen)
	_, err = token.Verify(strings.Empty, "aud")
	require.Error(t, err)

	token = paseto.NewToken(nil, signer, verifier, gen)
	require.Nil(t, token)
}
