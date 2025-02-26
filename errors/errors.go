package errors

import "fmt"

// Prefix an error.
func Prefix(prefix string, err error) error {
	if !IsError(err) {
		return nil
	}

	return fmt.Errorf("%v: %w", prefix, err)
}

// IsError returns true if err != nil, otherwise false.
func IsError(err error) bool {
	return err != nil
}
