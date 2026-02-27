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

// NewToken constructs a Token facade that can generate and verify tokens for multiple kinds.
//
// The facade delegates to the implementation selected by cfg.Kind:
//
//   - "jwt": token/jwt
//   - "paseto": token/paseto
//   - "ssh": token/ssh
//
// The underlying implementations are constructed eagerly from the corresponding nested
// configuration blocks (cfg.JWT, cfg.Paseto, cfg.SSH). Individual implementations may
// be nil when their nested configuration is nil.
//
// Important: NewToken does not validate cfg or enforce that the selected kind has a
// non-nil nested config. If cfg.Kind selects an implementation whose constructor returned
// nil, calling Generate/Verify for that kind will typically panic due to a nil receiver.
// Ensure your configuration is consistent with the selected kind.
//
// Unknown kinds are treated as "disabled" by the facade methods: Generate returns (nil, nil)
// and Verify returns (strings.Empty, nil).
func NewToken(name env.Name, cfg *Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	return &Token{
		name: name, cfg: cfg,
		jwt:    jwt.NewToken(cfg.JWT, sig, ver, gen),
		paseto: paseto.NewToken(cfg.Paseto, sig, ver, gen),
		ssh:    ssh.NewToken(cfg.SSH, fs),
	}
}

// Token is a facade that generates and verifies tokens using the implementation selected by configuration.
//
// It standardizes the call sites for token issuance and validation, while allowing the actual
// token format and crypto scheme to be chosen by configuration.
type Token struct {
	cfg    *Config
	jwt    *jwt.Token
	paseto *paseto.Token
	ssh    *ssh.Token
	name   env.Name
}

// Generate creates a token for the configured kind.
//
// Semantics by kind:
//
//   - "jwt" and "paseto": the token is minted for the provided audience (aud) and
//     subject (sub).
//
//   - "ssh": audience and subject are ignored; the SSH token kind uses its own
//     encoding/signature format and typically identifies a key rather than a subject.
//
// If the configured kind is unknown, Generate returns (nil, nil).
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

// Verify validates token for the configured kind and returns the subject identifier.
//
// Semantics by kind:
//
//   - "jwt" and "paseto": verifies the token for the provided audience (aud) and
//     returns the subject ("sub") claim.
//
//   - "ssh": audience is ignored and the returned string is the selected key name
//     (not a JWT/PASETO "sub" claim).
//
// If the configured kind is unknown, Verify returns (strings.Empty, nil).
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
