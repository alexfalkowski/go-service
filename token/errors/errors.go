package errors

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrInvalidMatch for token.
	ErrInvalidMatch = errors.New("token: invalid match")

	// ErrInvalidIssuer for token.
	ErrInvalidIssuer = errors.New("token: invalid issuer")

	// ErrInvalidAudience for token.
	ErrInvalidAudience = errors.New("token: invalid audience")

	// ErrInvalidAlgorithm for token.
	ErrInvalidAlgorithm = errors.New("token: invalid algorithm")

	// ErrInvalidKeyID for token.
	ErrInvalidKeyID = errors.New("token: invalid key id")

	// ErrInvalidTime for token.
	ErrInvalidTime = errors.New("token: invalid time")
)
