package meta

import (
	"context"

	"github.com/form3tech-oss/jwt-go"
)

type contextKey string

var (
	token = contextKey("token")
)

// WithToken for security.
func WithToken(ctx context.Context, tkn *jwt.Token) context.Context {
	return context.WithValue(ctx, token, tkn)
}

// Token for security.
func Token(ctx context.Context) *jwt.Token {
	token, ok := ctx.Value(token).(*jwt.Token)
	if ok {
		return token
	}

	return nil
}
