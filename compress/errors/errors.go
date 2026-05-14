package errors

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrTooLarge is returned when data exceeds the configured size.
var ErrTooLarge = errors.New("compress: too large")
