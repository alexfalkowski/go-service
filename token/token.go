package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
)

// NewToken based on config.
func NewToken(name env.Name, cfg *Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	return &Token{
		name: name, cfg: cfg,
		jwt:    jwt.NewToken(cfg.JWT, sig, ver, gen),
		paseto: paseto.NewToken(cfg.Paseto, sig, ver, gen),
		ssh:    ssh.NewToken(cfg.SSH, fs),
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
	switch t.cfg.Kind {
	case "jwt":
		token, err := t.jwt.Generate(aud, sub)
		return strings.Bytes(token), err
	case "paseto":
		token, err := t.paseto.Generate(aud, sub)
		return strings.Bytes(token), err
	case "ssh":
		token, err := t.ssh.Generate()
		return strings.Bytes(token), err
	default:
		return nil, nil
	}
}

// Verify a token based on kind.
func (t *Token) Verify(token []byte, aud string) (string, error) {
	switch t.cfg.Kind {
	case "jwt":
		return t.jwt.Verify(bytes.String(token), aud)
	case "paseto":
		return t.paseto.Verify(bytes.String(token), aud)
	case "ssh":
		return t.ssh.Verify(bytes.String(token))
	default:
		return "", nil
	}
}
