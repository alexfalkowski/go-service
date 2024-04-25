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
func ConvertRecover(r any) error {
	switch x := r.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	default:
		return errors.New(fmt.Sprint(x))
	}
}
