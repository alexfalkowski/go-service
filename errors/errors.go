package errors

import (
	"errors"
	"fmt"
)

var (
	// As is an alias for errors.As.
	As = errors.As

	// New is an alias for errors.New.
	New = errors.New

	// Is is an alias for errors.Is.
	Is = errors.Is

	// Join is an alias for errors.Join.
	Join = errors.Join
)

// Prefix an error.
func Prefix(prefix string, err error) error {
	if err != nil {
		return fmt.Errorf("%v: %w", prefix, err)
	}

	return nil
}
