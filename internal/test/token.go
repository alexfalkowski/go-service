package test

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
)

// WithWorldToken for test.
func WithWorldToken(generator token.Generator, verifier token.Verifier) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.generator = generator
		o.verifier = verifier
	})
}

// NewGenerator for test.
func NewGenerator(token string, err error) *Generator {
	return &Generator{token: token, err: err}
}

// Generator for test.
type Generator struct {
	err   error
	token string
}

func (g *Generator) Generate(_, _ string) ([]byte, error) {
	return strings.Bytes(g.token), g.err
}

// NewVerifier for test.
func NewVerifier(token string) *Verifier {
	return &Verifier{token: token}
}

// Verifier for test.
type Verifier struct {
	token string
}

func (v *Verifier) Verify(token []byte, aud string) (string, error) {
	if bytes.String(token) != v.token {
		return strings.Empty, ErrInvalid
	}

	return UserID.String(), nil
}
