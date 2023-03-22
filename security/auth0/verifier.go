package auth0

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var (
	// ErrInvalidAudience for Auth0.
	ErrInvalidAudience = errors.New("invalid audience")

	// ErrInvalidIssuer for Auth0.
	ErrInvalidIssuer = errors.New("invalid issuer")

	// ErrInvalidAlgorithm for Auth0.
	ErrInvalidAlgorithm = errors.New("invalid algorithm")
)

type verifier struct {
	cfg  *Config
	cert Certificator
}

func (v *verifier) Verify(_ context.Context, token []byte) (*jwt.Token, *jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	t, err := jwt.ParseWithClaims(string(token), claims, v.key)
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

func (v *verifier) key(token *jwt.Token) (any, error) {
	cert, err := v.cert.Certificate(context.Background(), token)
	if err != nil {
		return token, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	if err != nil {
		return token, err
	}

	return key, nil
}
