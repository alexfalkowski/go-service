package errors

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrInvalidType when we can't encode the type provided.
var ErrInvalidType = errors.New("encoding: invalid type")
