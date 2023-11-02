package oauth

import (
	"context"
	"errors"

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

func (v *verifier) Verify(_ context.Context, token []byte) (*jwt.Token, *jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	t, err := jwt.ParseWithClaims(string(token), claims, v.validate)
	if err != nil {
		return nil, claims, err
	}

	if t.Header["alg"] != v.cfg.Algorithm {
		return t, claims, ErrInvalidAlgorithm
	}

	if !claims.VerifyIssuer(v.cfg.Issuer, true) {
		return t, claims, ErrInvalidIssuer
	}

	if !claims.VerifyAudience(v.cfg.Audience, true) {
		return t, claims, ErrInvalidAudience
	}

	return t, claims, nil
}

func (v *verifier) validate(token *jwt.Token) (any, error) {
	return v.cert.Certificate(context.Background(), token)
}
