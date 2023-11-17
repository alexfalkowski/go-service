package security

import (
	"github.com/alexfalkowski/go-service/security/token"
)

// NewToken for security.
func NewToken(c *token.Config) (token.Generator, token.Verifier) {
	return c.Generator(), c.Verifier()
}
