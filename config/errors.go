package config

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrNoEncoder is returned when no encoder is provided.
	ErrNoEncoder = errors.New("config: no encoder")

	// ErrEnvMissing is returned when env is missing.
	ErrEnvMissing = errors.New("config: env is missing")

	// ErrLocationMissing is returned when location is missing.
	ErrLocationMissing = errors.New("config: location is missing")

	// ErrInvalidConfig is returned when decoding fails.
	ErrInvalidConfig = errors.New("config: invalid format")
)
