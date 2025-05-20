package logger

import "github.com/alexfalkowski/go-service/v2/errors"

func prefix(err error) error {
	return errors.Prefix("logger", err)
}
