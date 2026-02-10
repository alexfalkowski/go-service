package errors

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrInvalidMatch is returned when a token comparison fails to match (for example, signature mismatch).
	ErrInvalidMatch = errors.New("token: invalid match")

	// ErrInvalidIssuer is returned when a token issuer claim does not match the expected issuer.
	ErrInvalidIssuer = errors.New("token: invalid issuer")

	// ErrInvalidAudience is returned when a token audience claim does not match the expected audience.
	ErrInvalidAudience = errors.New("token: invalid audience")

	// ErrInvalidAlgorithm is returned when a token is signed with an unexpected algorithm.
	ErrInvalidAlgorithm = errors.New("token: invalid algorithm")

	// ErrInvalidKeyID is returned when a token key id (kid) header is missing or does not match the expected key id.
	ErrInvalidKeyID = errors.New("token: invalid key id")

	// ErrInvalidTime is returned when a token is not valid for the current time (for example expired or not yet valid).
	ErrInvalidTime = errors.New("token: invalid time")
)
