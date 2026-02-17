package config

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrNoEncoder is returned when no encoder is registered for the configuration kind.
	ErrNoEncoder = errors.New("config: no encoder")

	// ErrEnvMissing is returned when the env config variable is missing or malformed.
	ErrEnvMissing = errors.New("config: env is missing")

	// ErrLocationMissing is returned when no configuration file is found in search locations.
	ErrLocationMissing = errors.New("config: location is missing")

	// ErrInvalidConfig is returned when a decoded config is empty or invalid.
	ErrInvalidConfig = errors.New("config: invalid format")
)
