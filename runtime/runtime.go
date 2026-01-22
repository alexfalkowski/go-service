package runtime

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/errors"
)

var ErrRecovered = errors.New("recovered")

// Must panics if we have an error.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// ConvertRecover to an error.
func ConvertRecover(value any) error {
	switch kind := value.(type) {
	case error:
		return fmt.Errorf("%w: %w", ErrRecovered, kind)
	case string:
		return fmt.Errorf("%w: %s", ErrRecovered, kind)
	default:
		return fmt.Errorf("%w: %v", ErrRecovered, kind)
	}
}
