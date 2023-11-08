package meta

import (
	"context"
	"encoding/json"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/golang-jwt/jwt/v4"
)

const (
	registeredClaims = "jwt.claims"
)

// WithRegisteredClaims for jwt.
func WithRegisteredClaims(ctx context.Context, claims *jwt.RegisteredClaims) (context.Context, error) {
	c, err := json.Marshal(claims)
	if err != nil {
		return ctx, err
	}

	return meta.WithAttribute(ctx, registeredClaims, string(c)), nil
}

// RegisteredClaims for jwt.
func RegisteredClaims(ctx context.Context) (*jwt.RegisteredClaims, error) {
	var c jwt.RegisteredClaims
	if err := json.Unmarshal([]byte(meta.Attribute(ctx, registeredClaims)), &c); err != nil {
		return nil, err
	}

	return &c, nil
}
