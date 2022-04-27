package auth0

import (
	"context"
	"errors"
	"sync"

	"github.com/form3tech-oss/jwt-go"
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
	ctx  context.Context
	mux  sync.Mutex
}

func (v *verifier) Verify(ctx context.Context, token []byte) (*jwt.Token, error) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.ctx = ctx

	parsedToken, err := jwt.Parse(string(token), v.validate)
	if err != nil {
		return nil, err
	}

	if parsedToken.Header["alg"] != v.cfg.Algorithm {
		return nil, ErrInvalidAlgorithm
	}

	return parsedToken, nil
}

// nolint:forcetypeassert
func (v *verifier) validate(token *jwt.Token) (any, error) {
	claims := token.Claims.(jwt.MapClaims)

	checkAud := claims.VerifyAudience(v.cfg.Audience, true)
	if !checkAud {
		return token, ErrInvalidAudience
	}

	checkIss := claims.VerifyIssuer(v.cfg.Issuer, true)
	if !checkIss {
		return token, ErrInvalidIssuer
	}

	cert, err := v.cert.Certificate(v.ctx, token)
	if err != nil {
		return token, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	if err != nil {
		return token, err
	}

	return key, nil
}
