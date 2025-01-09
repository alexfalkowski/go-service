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
//
//nolint:err113
func ConvertRecover(rec any) error {
	switch recType := rec.(type) {
	case string:
		return errors.New(recType)
	case error:
		return recType
	default:
		return errors.New(fmt.Sprint(recType))
	}
}
