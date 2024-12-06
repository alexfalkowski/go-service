package token

import (
	"context"
	"errors"
	"time"

	st "github.com/alexfalkowski/go-service/time"
)

var (
	// ErrInvalidAlgorithm for service.
	ErrInvalidAlgorithm = errors.New("invalid algorithm")

	// ErrInvalidIssuer for service.
	ErrInvalidIssuer = errors.New("invalid issuer")

	// ErrInvalidAudience for service.
	ErrInvalidAudience = errors.New("invalid audience")

	// ErrInvalidTime for service.
	ErrInvalidTime = errors.New("invalid time")
)

type (
	// Generator allows the implementation of different types generators.
	Generator interface {
		// Generate a new token or error.
		Generate(ctx context.Context) (context.Context, []byte, error)
	}

	// Verifier allows the implementation of different types of verifiers.
	Verifier interface {
		// Verify a token or error.
		Verify(ctx context.Context, token []byte) (context.Context, error)
	}

	token interface {
		Generate(sub, aud, iss string, exp time.Duration) (string, error)
		Verify(token, aud, iss string) (string, error)
	}

	// Token will generate and verify based on what is defined in the config.
	Token struct {
		cfg   *Config
		token token
	}
)

// NewToken based on config.
func NewToken(cfg *Config, jwt *JWT, pas *Paseto) *Token {
	if !IsEnabled(cfg) {
		return &Token{}
	}

	var token token

	switch {
	case cfg.IsJWT():
		token = jwt
	case cfg.IsPaseto():
		token = pas
	}

	return &Token{cfg: cfg, token: token}
}

func (t *Token) Generate(ctx context.Context) (context.Context, []byte, error) {
	if t.token == nil {
		return ctx, nil, nil
	}

	token, err := t.token.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, st.MustParseDuration(t.cfg.Expiration))

	return ctx, []byte(token), err
}

func (t *Token) Verify(ctx context.Context, token []byte) (context.Context, error) {
	if t.token == nil {
		return ctx, nil
	}

	_, err := t.token.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

	return ctx, err
}
