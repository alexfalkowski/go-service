package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
)

// TokenParams for token.
type TokenParams struct {
	di.In
	Config *Config
	JWT    *jwt.Token
	Paseto *paseto.Token
	SSH    *ssh.Token
	Name   env.Name
}

// NewToken based on config.
func NewToken(params TokenParams) *Token {
	if !params.Config.IsEnabled() {
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

// Generate a token based on kind.
func (t *Token) Generate(aud, sub string) ([]byte, error) {
	switch {
	case t.cfg.IsSSH():
		token, err := t.ssh.Generate()

		return strings.Bytes(token), err
	case t.cfg.IsJWT():
		token, err := t.jwt.Generate(aud, sub)

		return strings.Bytes(token), err
	case t.cfg.IsPaseto():
		token, err := t.paseto.Generate(aud, sub)

		return strings.Bytes(token), err
	default:
		return nil, nil
	}
}

// Verify a token based on kind.
func (t *Token) Verify(token []byte, aud string) (string, error) {
	tkn := bytes.String(token)
	var (
		user string
		err  error
	)

	switch {
	case t.cfg.IsSSH():
		user, err = t.ssh.Verify(tkn)
	case t.cfg.IsJWT():
		user, err = t.jwt.Verify(tkn, aud)
	case t.cfg.IsPaseto():
		user, err = t.paseto.Verify(tkn, aud)
	}

	return user, err
}
