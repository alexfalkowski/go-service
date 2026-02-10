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

// NewToken returns a Token that generates and verifies tokens according to cfg.
//
// Supported kinds are "jwt", "paseto", and "ssh". If an unsupported kind is configured,
// Generate/Verify return nil/empty results.
func NewToken(name env.Name, cfg *Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	return &Token{
		name: name, cfg: cfg,
		jwt:    jwt.NewToken(cfg.JWT, sig, ver, gen),
		paseto: paseto.NewToken(cfg.Paseto, sig, ver, gen),
		ssh:    ssh.NewToken(cfg.SSH, fs),
	}
}

// Token generates and verifies tokens using the implementation selected by configuration.
type Token struct {
	cfg    *Config
	jwt    *jwt.Token
	paseto *paseto.Token
	ssh    *ssh.Token
	name   env.Name
}

// Generate creates a token for the configured kind.
//
// For "jwt" and "paseto" the token is created for the provided audience and subject.
// For "ssh" the audience and subject are ignored.
//
// If the configured kind is unknown, it returns (nil, nil).
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

// Verify validates token for the configured kind and returns the subject.
//
// For "ssh" the audience is ignored and the returned string is the key name.
//
// If the configured kind is unknown, it returns (strings.Empty, nil).
func (t *Token) Verify(token []byte, aud string) (string, error) {
	switch t.cfg.Kind {
	case "jwt":
		return t.jwt.Verify(bytes.String(token), aud)
	case "paseto":
		return t.paseto.Verify(bytes.String(token), aud)
	case "ssh":
		return t.ssh.Verify(bytes.String(token))
	default:
		return strings.Empty, nil
	}
}
