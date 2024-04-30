package errors

import (
	"fmt"
)

// Prefix an error.
func Prefix(p string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%v: %w", p, err)
}
