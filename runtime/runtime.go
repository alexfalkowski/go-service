package runtime

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/errors"
)

// ErrRecovered is used to wrap panic values converted to an error.
var ErrRecovered = errors.New("recovered")

// Must panics if we have an error.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// ConvertRecover converts a recovered panic value into an error wrapped with ErrRecovered.
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
