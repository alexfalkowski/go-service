package config

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrNoEncoder for encoding.
	ErrNoEncoder = errors.New("config: no encoder")

	// ErrLocationMissing for cmd.
	ErrLocationMissing = errors.New("config: location is missing")

	// ErrInvalidConfig when decoding fails.
	ErrInvalidConfig = errors.New("config: invalid format")
)
