package oauth

import (
	"context"
	"crypto"

	"github.com/golang-jwt/jwt/v4"
)

// Certificator for OAuth.
type Certificator interface {
	Certificate(ctx context.Context, token *jwt.Token) (crypto.PublicKey, error)
}
