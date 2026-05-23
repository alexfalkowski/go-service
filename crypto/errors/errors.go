package errors

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrInvalidMatch indicates that a cryptographic verification check failed.
//
// It is used by verification helpers (for example signature verifiers) to report that a provided
// signature/hash does not match the expected value for the input message.
var ErrInvalidMatch = errors.New("crypto: invalid match")

// ErrMissingKey indicates that required cryptographic key material was not configured
// or resolved to empty bytes.
var ErrMissingKey = errors.New("crypto: missing key")

// ErrInvalidKeyType indicates that parsed key material was valid, but not of the
// key type expected by the caller.
//
// It is used by key-loading helpers that only support a specific algorithm
// (for example Ed25519) when the provided data decodes successfully as some
// other key type.
var ErrInvalidKeyType = errors.New("crypto: invalid key type")

// ErrInvalidKeySize indicates that parsed key material was valid, but its size
// does not satisfy the package policy for the cryptographic operation.
var ErrInvalidKeySize = errors.New("crypto: invalid key size")

// ErrInvalidKeyFormat indicates that parsed key material used a syntactic form
// that is not supported by the package.
var ErrInvalidKeyFormat = errors.New("crypto: invalid key format")
