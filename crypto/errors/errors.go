package errors

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrInvalidMatch is returned when a signature/hash comparison fails to match.
var ErrInvalidMatch = errors.New("crypto: invalid match")
