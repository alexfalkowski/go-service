package errors

import "errors"

// ErrInvalidType when we can't encode the type provided.
var ErrInvalidType = errors.New("encoding: invalid type")
