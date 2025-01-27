package token

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
)

var (
	// ErrInvalidMatch for token.
	ErrInvalidMatch = errors.New("token: invalid match")

	// ErrInvalidIssuer for service.
	ErrInvalidIssuer = errors.New("token: invalid issuer")

	// ErrInvalidAudience for service.
	ErrInvalidAudience = errors.New("token: invalid audience")

	// ErrInvalidTime for service.
	ErrInvalidTime = errors.New("token: invalid time")
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
)

// NewToken based on config.
func NewToken(cfg *Config, jwt *JWT, pas *Paseto) *Token {
	return &Token{cfg: cfg, jwt: jwt, pas: pas}
}

// Token will generate and verify based on what is defined in the config.
type Token struct {
	cfg *Config
	jwt *JWT
	pas *Paseto
}

func (t *Token) Generate(ctx context.Context) (context.Context, []byte, error) {
	if t.cfg == nil {
		return ctx, nil, nil
	}

	switch {
	case t.cfg.IsKey():
		d, err := os.ReadFile(t.cfg.Key)

		return ctx, []byte(d), err
	case t.cfg.IsJWT():
		token, err := t.jwt.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, time.MustParseDuration(t.cfg.Expiration))

		return ctx, []byte(token), err
	case t.cfg.IsPaseto():
		token, err := t.pas.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, time.MustParseDuration(t.cfg.Expiration))

		return ctx, []byte(token), err
	}

	return ctx, nil, nil
}

func (t *Token) Verify(ctx context.Context, token []byte) (context.Context, error) {
	if t.cfg == nil {
		return ctx, nil
	}

	switch {
	case t.cfg.IsKey():
		file, err := os.ReadFile(t.cfg.Key)
		if err != nil {
			return ctx, err
		}

		if err := VerifyKey(file); err != nil {
			return ctx, fmt.Errorf("%w: %w", err, ErrInvalidMatch)
		}

		if !bytes.Equal([]byte(file), token) {
			return ctx, ErrInvalidMatch
		}

		return ctx, nil
	case t.cfg.IsJWT():
		_, err := t.jwt.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

		return ctx, err
	case t.cfg.IsPaseto():
		_, err := t.pas.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

		return ctx, err
	}

	return ctx, nil
}
