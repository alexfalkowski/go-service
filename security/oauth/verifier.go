package oauth

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/security/oauth/meta"
	"github.com/golang-jwt/jwt/v4"
)

var (
	// ErrInvalidAudience for OAuth.
	ErrInvalidAudience = errors.New("invalid audience")

	// ErrInvalidIssuer for OAuth.
	ErrInvalidIssuer = errors.New("invalid issuer")

	// ErrInvalidAlgorithm for OAuth.
	ErrInvalidAlgorithm = errors.New("invalid algorithm")
)

type verifier struct {
	cfg  *Config
	cert Certificator
}

func (v *verifier) Verify(ctx context.Context, token []byte) (context.Context, error) {
	claims := &jwt.RegisteredClaims{}

	t, err := jwt.ParseWithClaims(string(token), claims, v.validate)
	if err != nil {
		return ctx, err
	}

	if t.Header["alg"] != v.cfg.Algorithm {
		return ctx, ErrInvalidAlgorithm
	}

	if !claims.VerifyIssuer(v.cfg.Issuer, true) {
		return ctx, ErrInvalidIssuer
	}

	if !claims.VerifyAudience(v.cfg.Audience, true) {
		return ctx, ErrInvalidAudience
	}

	return meta.WithRegisteredClaims(ctx, claims)
}

func (v *verifier) validate(token *jwt.Token) (any, error) {
	return v.cert.Certificate(context.Background(), token)
}
