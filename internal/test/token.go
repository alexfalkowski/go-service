package test

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
)

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

// SequenceGenerator is a token.Generator test double that returns token-N values.
type SequenceGenerator struct {
	next int
}

// Generate implements token.Generator and returns the next token value.
func (g *SequenceGenerator) Generate(_, _ string) ([]byte, error) {
	g.next++

	return []byte("token-" + strconv.Itoa(g.next)), nil
}

// AcceptingVerifier is a token.Verifier test double that accepts any token.
type AcceptingVerifier struct{}

// Verify implements token.Verifier and returns the shared test user ID.
func (AcceptingVerifier) Verify([]byte, string) (string, error) {
	return UserID.String(), nil
}
