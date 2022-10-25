package auth0

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

// Certificator for Auth0.
type Certificator interface {
	Certificate(ctx context.Context, token *jwt.Token) (string, error)
}
