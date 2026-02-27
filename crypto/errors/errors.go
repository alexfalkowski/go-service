package errors

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrInvalidMatch indicates that a cryptographic verification check failed.
//
// It is used by verification helpers (for example signature verifiers) to report that a provided
// signature/hash does not match the expected value for the input message.
var ErrInvalidMatch = errors.New("crypto: invalid match")
