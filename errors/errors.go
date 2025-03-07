package errors

import "fmt"

// Prefix an error.
func Prefix(prefix string, err error) error {
	if err != nil {
		return fmt.Errorf("%v: %w", prefix, err)
	}

	return nil
}
