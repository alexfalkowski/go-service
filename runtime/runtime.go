package runtime

import (
	"errors"
	"fmt"
)

// Must panics if we have an error.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// ConvertRecover to an error.
func ConvertRecover(value any) error {
	switch kind := value.(type) {
	case string:
		return errors.New(kind)
	case error:
		return kind
	default:
		return errors.New(fmt.Sprint(kind))
	}
}
