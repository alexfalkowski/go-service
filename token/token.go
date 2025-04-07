package token

import (
	"bytes"
	"context"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token/errors"
	"github.com/alexfalkowski/go-service/token/ssh"
	"go.uber.org/fx"
)

// Params for token.
type Params struct {
	fx.In

	Config *Config
	JWT    *JWT
	Paseto *Paseto
	Opaque *Opaque
	SSH    *ssh.Token
	Name   env.Name
}

// NewToken based on config.
func NewToken(params Params) *Token {
	if !IsEnabled(params.Config) {
		return nil
	}

	return &Token{
		cfg:    params.Config,
		name:   params.Name,
		jwt:    params.JWT,
		paseto: params.Paseto,
		opaque: params.Opaque,
		ssh:    params.SSH,
	}
}

// Token will generate and verify based on what is defined in the config.
type Token struct {
	cfg    *Config
	jwt    *JWT
	paseto *Paseto
	opaque *Opaque
	ssh    *ssh.Token
	name   env.Name
}

func (t *Token) Generate(ctx context.Context) (context.Context, []byte, error) {
	switch {
	case t.cfg.IsOpaque():
		b, err := os.ReadFile(t.cfg.Secret)

		return ctx, b, err

	case t.cfg.IsSSH():
		token, err := t.ssh.Generate()

		return ctx, []byte(token), err
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
	switch {
	case t.cfg.IsOpaque():
		b, err := os.ReadFile(t.cfg.Secret)
		if err != nil {
			return ctx, err
		}

		if !bytes.Equal(b, token) {
			return ctx, errors.ErrInvalidMatch
		}

		return ctx, t.opaque.Verify(string(token))
	case t.cfg.IsSSH():
		return ctx, t.ssh.Verify(string(token))
	case t.cfg.IsJWT():
		subject, err := t.jwt.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

		return WithSubject(ctx, meta.String(subject)), err
	case t.cfg.IsPaseto():
		subject, err := t.paseto.Verify(string(token), t.cfg.Audience, t.cfg.Issuer)

		return WithSubject(ctx, meta.String(subject)), err
	}

	return ctx, nil
}
