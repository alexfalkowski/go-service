package runtime

import (
	"errors"
	"fmt"
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
		return kind
	case string:
		return fmt.Errorf("%w: %s", ErrRecovered, kind)
	default:
		return fmt.Errorf("%w: %s", ErrRecovered, kind)
	}
}
