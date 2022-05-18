package errors

import (
	"go.uber.org/multierr"
)

// Combine errors into a single error, making sure we do not have any nils.
func Combine(errs ...error) error {
	newErrs := make([]error, 0)

	for _, err := range errs {
		if err != nil {
			newErrs = append(newErrs, err)
		}
	}

	if len(newErrs) > 0 {
		return multierr.Combine(errs...)
	}

	return nil
}
