package errors

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrTooLarge is returned when uncompressed data exceeds the configured size limit.
var ErrTooLarge = errors.New("compress: too large")
