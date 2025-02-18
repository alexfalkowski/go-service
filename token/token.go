package token

import (
	"bytes"
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/fx"
)

// Params for token.
type Params struct {
	fx.In

	Config *Config
	JWT    *JWT
	Paseto *Paseto
	Opaque *Opaque
	Name   env.Name
}

// NewToken based on config.
func NewToken(params Params) *Token {
	return &Token{
		cfg:    params.Config,
		name:   params.Name,
		jwt:    params.JWT,
		paseto: params.Paseto,
		opaque: params.Opaque,
	}
}

// Token will generate and verify based on what is defined in the config.
type Token struct {
	cfg    *Config
	jwt    *JWT
	paseto *Paseto
	opaque *Opaque
	name   env.Name
}

func (t *Token) Generate(ctx context.Context) (context.Context, []byte, error) {
	if t.cfg == nil {
		return ctx, nil, nil
	}

	switch {
	case t.cfg.IsToken():
		d, err := os.ReadFile(t.cfg.Secret)

		return ctx, []byte(d), err
	case t.cfg.IsJWT():
		token, err := t.jwt.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, time.MustParseDuration(t.cfg.Expiration))

		return ctx, []byte(token), err
	case t.cfg.IsPaseto():
		token, err := t.paseto.Generate(t.cfg.Subject, t.cfg.Audience, t.cfg.Issuer, time.MustParseDuration(t.cfg.Expiration))

		return ctx, []byte(token), err
	}

	return ctx, nil, nil
}

func (t *Token) Verify(ctx context.Context, token []byte) (context.Context, error) {
	if t.cfg == nil {
		return ctx, nil
	}

	switch {
	case t.cfg.IsToken():
		d, err := os.ReadFile(t.cfg.Secret)
		if err != nil {
			return ctx, err
		}

		if !bytes.Equal([]byte(d), token) {
			return ctx, ErrInvalidMatch
		}

		return ctx, t.opaque.Verify(string(token))
	case t.cfg.IsJWT():
		_, err := t.jwt.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

		return ctx, err
	case t.cfg.IsPaseto():
		_, err := t.paseto.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

		return ctx, err
	}

	return ctx, nil
}
