package jwt_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		config *jwt.Config
		name   string
	}{
		{
			name: "missing issuer",
			config: &jwt.Config{
				KeyID:      "1234567890",
				Expiration: time.Hour,
			},
		},
		{
			name: "missing key id",
			config: &jwt.Config{
				Issuer:     "iss",
				Expiration: time.Hour,
			},
		},
		{
			name: "zero expiration",
			config: &jwt.Config{
				Issuer: "iss",
				KeyID:  "1234567890",
			},
		},
		{
			name: "negative expiration",
			config: &jwt.Config{
				Issuer:     "iss",
				KeyID:      "1234567890",
				Expiration: -time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, test.Validator.Struct(tt.config))
		})
	}

	cfg := &jwt.Config{Issuer: "iss", KeyID: "1234567890", Expiration: time.Hour}
	require.NoError(t, test.Validator.Struct(cfg))
}

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

	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "malformed token",
			value: "invalid",
		},
		{
			name:  "unexpected signing method",
			value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := token.Verify(tt.value, "hello")
			require.Error(t, err)
		})
	}

	t.Run("invalid audience", func(t *testing.T) {
		tkn, err := token.Generate("hello", test.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, tkn)

		sub, err := token.Verify(tkn, "test")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidAudience)
	})

	t.Run("invalid issuer", func(t *testing.T) {
		issuerToken := jwt.NewToken(&jwt.Config{Issuer: "test", Expiration: time.Hour, KeyID: "1234567890"}, signer, verifier, gen)

		tkn, err := issuerToken.Generate("hello", test.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, tkn)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidIssuer)
	})

	t.Run("disabled config", func(t *testing.T) {
		token := jwt.NewToken(nil, signer, verifier, gen)
		require.Nil(t, token)
	})
}

func TestInvalidSignature(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	token := jwt.NewToken(cfg.JWT, signer, verifier, gen)

	_, wrongPrivate, err := ed25519.NewGenerator(rand.NewGenerator(rand.NewReader())).Generate()
	require.NoError(t, err)

	wrongSigner, err := ed25519.NewSigner(test.PEM, &ed25519.Config{Private: wrongPrivate})
	require.NoError(t, err)

	wrong := jwt.NewToken(cfg.JWT, wrongSigner, verifier, gen)
	tkn, err := wrong.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := token.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
}

func TestInvalidKeyID(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	token := jwt.NewToken(cfg.JWT, signer, verifier, gen)

	t.Run("unexpected key id", func(t *testing.T) {
		wrong := jwt.NewToken(&jwt.Config{Issuer: cfg.JWT.Issuer, Expiration: time.Hour, KeyID: "test"}, signer, verifier, gen)

		tkn, err := wrong.Generate("hello", test.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, tkn)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidKeyID)
	})
}

func TestInvalidLifetimeExceedsConfig(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	generator := jwt.NewToken(&jwt.Config{Issuer: cfg.JWT.Issuer, Expiration: time.Hour, KeyID: cfg.JWT.KeyID}, signer, verifier, gen)
	verifierToken := jwt.NewToken(&jwt.Config{Issuer: cfg.JWT.Issuer, Expiration: time.Minute, KeyID: cfg.JWT.KeyID}, signer, verifier, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidTime)
}

func TestInvalidVerifyExpirationConfig(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	generator := jwt.NewToken(cfg.JWT, signer, verifier, gen)
	verifierToken := jwt.NewToken(&jwt.Config{Issuer: cfg.JWT.Issuer, KeyID: cfg.JWT.KeyID}, signer, verifier, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestInvalidConfigDoesNotPanic(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()

	t.Run("generate without signer", func(t *testing.T) {
		token := jwt.NewToken(cfg.JWT, nil, verifier, gen)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate without private key", func(t *testing.T) {
		token := jwt.NewToken(cfg.JWT, &ed25519.Signer{}, verifier, gen)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate without generator", func(t *testing.T) {
		token := jwt.NewToken(cfg.JWT, signer, verifier, nil)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify without verifier", func(t *testing.T) {
		valid := jwt.NewToken(cfg.JWT, signer, verifier, gen)
		tkn, err := valid.Generate("hello", test.UserID.String())
		require.NoError(t, err)

		token := jwt.NewToken(cfg.JWT, signer, nil, gen)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify without public key", func(t *testing.T) {
		valid := jwt.NewToken(cfg.JWT, signer, verifier, gen)
		tkn, err := valid.Generate("hello", test.UserID.String())
		require.NoError(t, err)

		token := jwt.NewToken(cfg.JWT, signer, &ed25519.Verifier{}, gen)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})
}

func TestInvalidPrivateKeyDoesNotPanic(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	token := jwt.NewToken(cfg.JWT, &ed25519.Signer{PrivateKey: []byte("short")}, verifier, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.Empty(t, tkn)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}
