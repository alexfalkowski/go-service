package errors

import (
	"errors"
	"fmt"
)

// As reports whether any error in err's chain matches target, and if so sets target to that error value.
//
// This is a thin wrapper around the standard library errors.As, provided so go-service code can depend
// on a stable import path.
//
// See: https://pkg.go.dev/errors#As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// New returns an error that formats as the given text.
//
// This is a thin wrapper around the standard library errors.New. It is typically used to define
// sentinel errors for comparisons with Is.
//
// See: https://pkg.go.dev/errors#New
func New(text string) error {
	return errors.New(text)
}

// Is reports whether any error in err's chain matches target.
//
// This is a thin wrapper around the standard library errors.Is.
//
// See: https://pkg.go.dev/errors#Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Join returns an error that wraps the given errors.
//
// This is a thin wrapper around the standard library errors.Join. A nil error in errs is ignored.
// If all errors are nil, Join returns nil.
//
// See: https://pkg.go.dev/errors#Join
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Prefix wraps err with a component prefix while preserving it for unwrapping.
//
// It formats the returned error using:
//
//	fmt.Errorf("%v: %w", prefix, err)
//
// If err is nil, Prefix returns nil. This makes it convenient to use in return statements without
// additional nil checks.
//
// Example:
//
//	return errors.Prefix("database", err)
func Prefix(prefix string, err error) error {
	if err != nil {
		return fmt.Errorf("%v: %w", prefix, err)
	}

	return nil
}
