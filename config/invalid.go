package config

import "github.com/alexfalkowski/go-service/v3/errors"

type invalidSourceDecoder struct {
	source string
}

func (d invalidSourceDecoder) Decode(any) error {
	return errors.Prefix("source "+d.source, ErrInvalidSource)
}
