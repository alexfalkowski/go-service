package errors

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrInvalidMatch is a sentinel error indicating a token comparison failed to match.
	//
	// Implementations typically use this for cryptographic or structural mismatches
	// (for example signature/MAC mismatch, invalid encoding, etc.) and should wrap
	// this error so callers can detect it via errors.Is while still preserving context.
	ErrInvalidMatch = errors.New("token: invalid match")

	// ErrInvalidIssuer is a sentinel error indicating the issuer claim is invalid.
	//
	// For claim-based token formats, this commonly corresponds to an "iss" claim
	// mismatch. Implementations should wrap this error to preserve additional context.
	ErrInvalidIssuer = errors.New("token: invalid issuer")

	// ErrInvalidAudience is a sentinel error indicating the audience claim is invalid.
	//
	// For claim-based token formats, this commonly corresponds to an "aud" claim
	// mismatch. Implementations should wrap this error to preserve additional context.
	ErrInvalidAudience = errors.New("token: invalid audience")

	// ErrInvalidAlgorithm is a sentinel error indicating the token used an unexpected algorithm.
	//
	// Implementations may return or wrap this when a token is signed/encrypted with an
	// algorithm that does not match the expected configuration.
	ErrInvalidAlgorithm = errors.New("token: invalid algorithm")

	// ErrInvalidKeyID is a sentinel error indicating a key identifier is missing or unexpected.
	//
	// This is commonly used for JWT verification when the "kid" header is missing or
	// does not match the expected configured key ID.
	ErrInvalidKeyID = errors.New("token: invalid key id")

	// ErrInvalidTime is a sentinel error indicating a token is not valid for the current time.
	//
	// This commonly covers expiration and not-before failures (for example expired tokens
	// or tokens that are not yet valid).
	ErrInvalidTime = errors.New("token: invalid time")
)
