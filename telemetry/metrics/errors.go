package metrics

import "github.com/alexfalkowski/go-service/v3/errors"

func prefix(err error) error {
	return errors.Prefix("metrics", err)
}
