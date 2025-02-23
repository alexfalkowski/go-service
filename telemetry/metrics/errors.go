package metrics

import "github.com/alexfalkowski/go-service/errors"

func prefix(err error) error {
	return errors.Prefix("metrics", err)
}
