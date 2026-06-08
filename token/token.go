package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	token "github.com/alexfalkowski/go-service/v2/token/errors"
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
// A nil cfg is treated as disabled and returns nil. If Kind selects an implementation whose
// nested config is missing, Generate/Verify return [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig] instead of panicking.
//
// Unknown kinds are treated as invalid configuration by the facade methods: Generate and Verify
// return [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig].
func NewToken(name env.Name, cfg *Config, fs *os.FS, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}

	return &Token{
		name: name, cfg: cfg,
		jwt:    jwt.NewToken(cfg.JWT, fs, gen),
		paseto: paseto.NewToken(cfg.Paseto, fs, gen),
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
//   - "ssh": the token is minted for the provided audience (aud), while subject
//     is derived from the active trusted key id.
//
// If the configured kind is unknown, Generate returns [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig].
//
// The returned bytes should be treated as read-only. Generate may return a
// zero-allocation byte view over the generated token string. Callers that need
// to mutate the bytes should clone the returned slice first.
func (t *Token) Generate(aud, sub string) ([]byte, error) {
	switch t.cfg.Kind {
	case "jwt":
		if t.jwt == nil {
			return nil, token.ErrInvalidConfig
		}
		tkn, err := t.jwt.Generate(aud, sub)
		return strings.Bytes(tkn), err
	case "paseto":
		if t.paseto == nil {
			return nil, token.ErrInvalidConfig
		}
		tkn, err := t.paseto.Generate(aud, sub)
		return strings.Bytes(tkn), err
	case "ssh":
		if t.ssh == nil {
			return nil, token.ErrInvalidConfig
		}
		tkn, err := t.ssh.Generate(aud, sub)
		return strings.Bytes(tkn), err
	default:
		return nil, token.ErrInvalidConfig
	}
}

// Verify validates token for the configured kind and returns the subject identifier.
//
// Semantics by kind:
//
//   - "jwt" and "paseto": verifies the token for the provided audience (aud) and
//     returns the subject ("sub") claim.
//
//   - "ssh": verifies the token for the provided audience (aud), and the
//     returned string is the "sub" claim, which must match the signed key id.
//
// If the configured kind is unknown, Verify returns [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig].
func (t *Token) Verify(tkn []byte, aud string) (string, error) {
	switch t.cfg.Kind {
	case "jwt":
		if t.jwt == nil {
			return strings.Empty, token.ErrInvalidConfig
		}
		sub, err := t.jwt.Verify(bytes.String(tkn), aud)
		return sub, invalidMatch(err)
	case "paseto":
		if t.paseto == nil {
			return strings.Empty, token.ErrInvalidConfig
		}
		sub, err := t.paseto.Verify(bytes.String(tkn), aud)
		return sub, invalidMatch(err)
	case "ssh":
		if t.ssh == nil {
			return strings.Empty, token.ErrInvalidConfig
		}
		sub, err := t.ssh.Verify(bytes.String(tkn), aud)
		return sub, invalidMatch(err)
	default:
		return strings.Empty, token.ErrInvalidConfig
	}
}

func invalidMatch(err error) error {
	if err == nil {
		return nil
	}

	if isTokenSentinel(err) {
		return err
	}

	if isInvalidMatch(err) {
		return errors.Join(token.ErrInvalidMatch, err)
	}

	if isPasetoRuleError(err) || isJWTValidationError(err) {
		return err
	}

	return errors.Join(token.ErrInvalidMatch, err)
}

func isTokenSentinel(err error) bool {
	return errors.Is(err, token.ErrInvalidConfig) ||
		errors.Is(err, token.ErrInvalidIssuer) ||
		errors.Is(err, token.ErrInvalidAudience) ||
		errors.Is(err, token.ErrInvalidSubject) ||
		errors.Is(err, token.ErrInvalidAlgorithm) ||
		errors.Is(err, token.ErrInvalidKeyID) ||
		errors.Is(err, token.ErrInvalidTime)
}

func isInvalidMatch(err error) bool {
	return errors.Is(err, jwt.ErrTokenMalformed) ||
		errors.Is(err, jwt.ErrTokenSignatureInvalid) ||
		errors.Is(err, paseto.TokenError{})
}

func isPasetoRuleError(err error) bool {
	return errors.Is(err, paseto.RuleError{})
}

func isJWTValidationError(err error) bool {
	var validation *jwt.ValidationError
	return errors.As(err, &validation)
}
