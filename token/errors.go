package token

import "errors"

var (
	// ErrInvalidMatch for token.
	ErrInvalidMatch = errors.New("token: invalid match")

	// ErrInvalidIssuer for service.
	ErrInvalidIssuer = errors.New("token: invalid issuer")

	// ErrInvalidAudience for service.
	ErrInvalidAudience = errors.New("token: invalid audience")

	// ErrInvalidTime for service.
	ErrInvalidTime = errors.New("token: invalid time")
)
