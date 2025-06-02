package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token/context"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"go.uber.org/fx"
)

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

func (t *Token) Generate(ctx context.Context) (context.Context, []byte, error) {
	switch {
	case t.cfg.IsSSH():
		token, err := t.ssh.Generate(ctx)

		return ctx, strings.Bytes(token), err
	case t.cfg.IsJWT():
		token, err := t.jwt.Generate(ctx)

		return ctx, strings.Bytes(token), err
	case t.cfg.IsPaseto():
		token, err := t.paseto.Generate(ctx)

		return ctx, strings.Bytes(token), err
	}

	return ctx, nil, nil
}

func (t *Token) Verify(ctx context.Context, token []byte) (context.Context, error) {
	tkn := bytes.String(token)

	switch {
	case t.cfg.IsSSH():
		return ctx, t.ssh.Verify(ctx, tkn)
	case t.cfg.IsJWT():
		ctx, err := t.jwt.Verify(ctx, tkn)

		return ctx, err
	case t.cfg.IsPaseto():
		ctx, err := t.paseto.Verify(ctx, tkn)

		return ctx, err
	}

	return ctx, nil
}
