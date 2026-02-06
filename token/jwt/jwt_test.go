package jwt_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	tokenerrors "github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	gjwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)

	cfg := test.NewToken("jwt")
	token := jwt.NewToken(cfg.JWT, signer, verifier, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, "hello")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestInvalid(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	token := jwt.NewToken(cfg.JWT, signer, verifier, gen)

	tokens := []string{
		"invalid",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	for _, tkn := range tokens {
		_, err := token.Verify(tkn, "hello")
		require.Error(t, err)
	}

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "test")
	require.Error(t, err)

	token = jwt.NewToken(&jwt.Config{Issuer: "test", Expiration: "1h", KeyID: "1234567890"}, signer, verifier, gen)

	tkn, err = token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	token = jwt.NewToken(cfg.JWT, signer, verifier, gen)
	_, err = token.Verify(tkn, "hello")
	require.Error(t, err)

	token = jwt.NewToken(nil, signer, verifier, gen)
	require.Nil(t, token)
}

func TestInvalidKeyID(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	wrong := jwt.NewToken(&jwt.Config{Issuer: cfg.JWT.Issuer, Expiration: "1h", KeyID: "test"}, signer, verifier, gen)

	tkn, err := wrong.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	token := jwt.NewToken(cfg.JWT, signer, verifier, gen)

	_, err = token.Verify(tkn, "hello")
	require.ErrorIs(t, err, tokenerrors.ErrInvalidKeyID)

	jwtToken := gjwt.NewWithClaims(gjwt.SigningMethodEdDSA, &gjwt.RegisteredClaims{})
	tkn, err = jwtToken.SignedString(signer.PrivateKey)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "hello")
	require.ErrorIs(t, err, tokenerrors.ErrInvalidKeyID)
}

func TestInvalidAlgorithm(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	token := jwt.NewToken(cfg.JWT, signer, verifier, gen)

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	jwtToken := gjwt.NewWithClaims(gjwt.SigningMethodHS256, &gjwt.RegisteredClaims{})
	tkn, err = jwtToken.SignedString([]byte("secret"))
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "hello")
	require.ErrorIs(t, err, tokenerrors.ErrInvalidAlgorithm)
}
