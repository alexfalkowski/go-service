package test

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
)

// WithWorldToken overrides the token generator and verifier used by world clients and servers.
func WithWorldToken(generator token.Generator, verifier token.Verifier) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.generator = generator
		o.verifier = verifier
	})
}

// NewGenerator returns a token generator test double that always yields the configured token and error.
func NewGenerator(token string, err error) *Generator {
	return &Generator{token: token, err: err}
}

// Generator is a token.Generator test double with fixed output.
type Generator struct {
	err   error
	token string
}

// Generate implements token.Generator and returns the configured token and error.
func (g *Generator) Generate(_, _ string) ([]byte, error) {
	return strings.Bytes(g.token), g.err
}

// NewVerifier returns a token verifier test double that accepts exactly one token value.
func NewVerifier(token string) *Verifier {
	return &Verifier{token: token}
}

// Verifier is a token.Verifier test double that validates a single expected token.
type Verifier struct {
	token string
}

// Verify implements token.Verifier and validates the token matches the configured value.
func (v *Verifier) Verify(token []byte, aud string) (string, error) {
	if bytes.String(token) != v.token {
		return strings.Empty, ErrInvalid
	}

	return UserID.String(), nil
}
