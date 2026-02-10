package errors

import (
	"errors"
	"fmt"
)

// As is an alias for errors.As.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// New is an alias for errors.New.
func New(text string) error {
	return errors.New(text)
}

// Is is an alias for errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Join is an alias for errors.Join.
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Prefix wraps err with prefix using fmt.Errorf("%v: %w", prefix, err).
//
// If err is nil, Prefix returns nil.
func Prefix(prefix string, err error) error {
	if err != nil {
		return fmt.Errorf("%v: %w", prefix, err)
	}

	return nil
}
