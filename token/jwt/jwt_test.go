package jwt_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	v4 "github.com/golang-jwt/jwt/v4"
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

	tokens := []string{
		"invalid",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	for _, tkn := range tokens {
		t.Run(tkn, func(t *testing.T) {
			_, err := token.Verify(tkn, "hello")
			require.Error(t, err)
		})
	}

	tkn, err := token.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "test")
	require.Error(t, err)

	token = jwt.NewToken(&jwt.Config{Issuer: "test", Expiration: time.Hour, KeyID: "1234567890"}, signer, verifier, gen)

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
	wrong := jwt.NewToken(&jwt.Config{Issuer: cfg.JWT.Issuer, Expiration: time.Hour, KeyID: "test"}, signer, verifier, gen)

	tkn, err := wrong.Generate("hello", test.UserID.String())
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	token := jwt.NewToken(cfg.JWT, signer, verifier, gen)

	_, err = token.Verify(tkn, "hello")
	require.ErrorIs(t, err, errors.ErrInvalidKeyID)

	jwtToken := v4.NewWithClaims(v4.SigningMethodEdDSA, &v4.RegisteredClaims{})
	tkn, err = jwtToken.SignedString(signer.PrivateKey)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "hello")
	require.ErrorIs(t, err, errors.ErrInvalidKeyID)
}

func TestInvalidMissingClaims(t *testing.T) {
	cfg := test.NewToken("jwt")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	gen := uuid.NewGenerator()
	token := jwt.NewToken(cfg.JWT, signer, verifier, gen)

	now := time.Now()
	tests := []struct {
		err    error
		mutate func(*v4.RegisteredClaims)
		name   string
	}{
		{
			name:   "missing expiration",
			err:    errors.ErrInvalidTime,
			mutate: func(claims *v4.RegisteredClaims) { claims.ExpiresAt = nil },
		},
		{
			name:   "missing issued at",
			err:    errors.ErrInvalidTime,
			mutate: func(claims *v4.RegisteredClaims) { claims.IssuedAt = nil },
		},
		{
			name:   "missing not before",
			err:    errors.ErrInvalidTime,
			mutate: func(claims *v4.RegisteredClaims) { claims.NotBefore = nil },
		},
		{
			name:   "missing subject",
			err:    errors.ErrInvalidSubject,
			mutate: func(claims *v4.RegisteredClaims) { claims.Subject = strings.Empty },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := validClaims(cfg.JWT, gen, now)
			tt.mutate(claims)

			jwtToken := v4.NewWithClaims(v4.SigningMethodEdDSA, claims)
			jwtToken.Header["kid"] = cfg.JWT.KeyID

			tkn, err := jwtToken.SignedString(signer.PrivateKey)
			require.NoError(t, err)

			sub, err := token.Verify(tkn, "hello")
			require.Empty(t, sub)
			require.ErrorIs(t, err, tt.err)
		})
	}
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

func validClaims(cfg *jwt.Config, gen *uuid.Generator, now time.Time) *v4.RegisteredClaims {
	return &v4.RegisteredClaims{
		ExpiresAt: &v4.NumericDate{Time: now.Add(time.Hour.Duration())},
		ID:        gen.Generate(),
		IssuedAt:  &v4.NumericDate{Time: now},
		Issuer:    cfg.Issuer,
		NotBefore: &v4.NumericDate{Time: now},
		Audience:  []string{"hello"},
		Subject:   test.UserID.String(),
	}
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

	jwtToken := v4.NewWithClaims(v4.SigningMethodHS256, &v4.RegisteredClaims{})
	tkn, err = jwtToken.SignedString([]byte("secret"))
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	_, err = token.Verify(tkn, "hello")
	require.ErrorIs(t, err, errors.ErrInvalidAlgorithm)
}
