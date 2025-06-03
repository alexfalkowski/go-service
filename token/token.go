package token

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"go.uber.org/fx"
)

// Options for token.
type Options struct {
	// Path that is used as an audience, such as users/1 or package.service/method
	Path string

	// UserID is the current user.
	UserID string
}

// Params for token.
type Params struct {
	fx.In

	Config *Config
	JWT    *jwt.Token
	Paseto *paseto.Token
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
		ssh:    params.SSH,
	}
}

// Token will generate and verify based on what is defined in the config.
type Token struct {
	cfg    *Config
	jwt    *jwt.Token
	paseto *paseto.Token
	ssh    *ssh.Token
	name   env.Name
}

func (t *Token) Generate(ctx context.Context, opts Options) (context.Context, []byte, error) {
	switch {
	case t.cfg.IsSSH():
		token, err := t.ssh.Generate()

		return ctx, strings.Bytes(token), err
	case t.cfg.IsJWT():
		token, err := t.jwt.Generate(opts.Path, opts.UserID)

		return ctx, strings.Bytes(token), err
	case t.cfg.IsPaseto():
		token, err := t.paseto.Generate(opts.Path, opts.UserID)

		return ctx, strings.Bytes(token), err
	}

	return ctx, nil, nil
}

func (t *Token) Verify(ctx context.Context, token []byte, opts Options) (context.Context, error) {
	tkn := bytes.String(token)

	switch {
	case t.cfg.IsSSH():
		return ctx, t.ssh.Verify(tkn)
	case t.cfg.IsJWT():
		subject, err := t.jwt.Verify(tkn, opts.Path)

		return WithSubject(ctx, meta.String(subject)), err
	case t.cfg.IsPaseto():
		subject, err := t.paseto.Verify(tkn, opts.Path)

		return WithSubject(ctx, meta.String(subject)), err
	}

	return ctx, nil
}
