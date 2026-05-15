package tracer

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrNotFound is returned when the configured tracer kind is unknown.
var ErrNotFound = errors.New("tracer: not found")

func prefix(err error) error {
	return errors.Prefix("tracer", err)
}
