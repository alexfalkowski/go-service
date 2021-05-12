package auth0

import (
	"context"

	"github.com/form3tech-oss/jwt-go"
)

// Certificator for Auth0.
type Certificator interface {
	Certificate(ctx context.Context, token *jwt.Token) (string, error)
}
