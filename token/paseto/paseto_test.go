package paseto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/stretchr/testify/require"
)

func TestConfigRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		config *paseto.Config
		name   string
	}{
		{
			name:   "missing issuer",
			config: &paseto.Config{Expiration: time.Hour},
		},
		{
			name:   "negative expiration",
			config: &paseto.Config{Issuer: "iss", Expiration: -time.Second},
		},
		{
			name:   "zero expiration",
			config: &paseto.Config{Issuer: "iss"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, test.Validator.Struct(tt.config))
		})
	}

	cfg := &paseto.Config{Issuer: "iss", Expiration: time.Hour}
	require.NoError(t, test.Validator.Struct(cfg))
}

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

func TestInvalidEmptySubject(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	token := paseto.NewToken(cfg.Paseto, signer, verifier, uuid.NewGenerator())

	tkn, err := token.Generate("hello", strings.Empty)
	require.NoError(t, err)

	sub, err := token.Verify(tkn, "hello")

	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidSubject)
}

func TestInvalidLifetimeExceedsConfig(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	generator := paseto.NewToken(&paseto.Config{Issuer: cfg.Paseto.Issuer, Expiration: time.Hour}, signer, verifier, gen)
	verifierToken := paseto.NewToken(&paseto.Config{Issuer: cfg.Paseto.Issuer, Expiration: time.Minute}, signer, verifier, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidTime)
}

func TestInvalidVerifyExpirationConfig(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	generator := paseto.NewToken(cfg.Paseto, signer, verifier, gen)
	verifierToken := paseto.NewToken(&paseto.Config{Issuer: cfg.Paseto.Issuer}, signer, verifier, gen)

	tkn, err := generator.Generate("hello", test.UserID.String())
	require.NoError(t, err)

	sub, err := verifierToken.Verify(tkn, "hello")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
}

func TestInvalidGenerateExpirationConfig(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	token := paseto.NewToken(&paseto.Config{Issuer: cfg.Paseto.Issuer}, signer, verifier, uuid.NewGenerator())

	tkn, err := token.Generate("hello", test.UserID.String())
	require.Empty(t, tkn)
	require.ErrorIs(t, err, errors.ErrInvalidConfig)
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

	token = paseto.NewToken(&paseto.Config{Issuer: "test", Expiration: time.Hour}, signer, verifier, gen)

	tkn, err = token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "hello")
	require.NoError(t, err)

	for _, tkn := range []string{"invalid"} {
		t.Run(tkn, func(t *testing.T) {
			cfg := test.NewToken("paseto")
			token := paseto.NewToken(cfg.Paseto, signer, verifier, gen)

			_, err := token.Verify(tkn, "aud")
			require.Error(t, err)
		})
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

func TestInvalidConfigDoesNotPanic(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()

	t.Run("generate without signer", func(t *testing.T) {
		token := paseto.NewToken(cfg.Paseto, nil, verifier, gen)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate without private key", func(t *testing.T) {
		token := paseto.NewToken(cfg.Paseto, &ed25519.Signer{}, verifier, gen)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate without generator", func(t *testing.T) {
		token := paseto.NewToken(cfg.Paseto, signer, verifier, nil)

		tkn, err := token.Generate("hello", test.UserID.String())
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify without verifier", func(t *testing.T) {
		valid := paseto.NewToken(cfg.Paseto, signer, verifier, gen)
		tkn, err := valid.Generate("hello", test.UserID.String())
		require.NoError(t, err)

		token := paseto.NewToken(cfg.Paseto, signer, nil, gen)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify without public key", func(t *testing.T) {
		valid := paseto.NewToken(cfg.Paseto, signer, verifier, gen)
		tkn, err := valid.Generate("hello", test.UserID.String())
		require.NoError(t, err)

		token := paseto.NewToken(cfg.Paseto, signer, &ed25519.Verifier{}, gen)

		sub, err := token.Verify(tkn, "hello")
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})
}
