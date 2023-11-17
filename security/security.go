package security

import (
	"github.com/alexfalkowski/go-service/security/token"
)

// NewToken for security.
func NewToken(c *token.Config) (token.Generator, token.Verifier) {
	token.RegisterGenerator("none", token.NewGenerator())
	token.RegisterVerifier("none", token.NewVerifier())

	return c.Generator(), c.Verifier()
}
