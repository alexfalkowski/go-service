package config

import "github.com/alexfalkowski/go-service/v2/errors"

var (
	// ErrNoEncoder is returned when the configuration kind cannot be decoded because there is no
	// registered encoder/decoder for it.
	//
	// It may be returned by:
	//   - the file decoder when the file extension does not map to a known kind, or
	//   - the env decoder when the env-provided kind is not registered.
	ErrNoEncoder = errors.New("config: no encoder")

	// ErrEnvMissing is returned when env-based configuration is missing or malformed.
	//
	// It is returned when the configured environment variable is unset or when its value does not
	// match the expected "<kind>:<base64-content>" format (for example, missing the kind or data).
	ErrEnvMissing = errors.New("config: env is missing")

	// ErrLocationMissing is returned when default lookup cannot find a configuration file in any of
	// the search locations.
	//
	// The default lookup searches for "<serviceName>.{yaml,yml,toml,json}" under common directories
	// (e.g. executable dir, user config dir, and /etc).
	ErrLocationMissing = errors.New("config: location is missing")

	// ErrInvalidConfig is returned when a decoded configuration value is considered empty.
	//
	// It is returned by NewConfig[T] when the decoded configuration does not populate any fields
	// (as determined by structs.IsEmpty), which helps prevent accidentally starting with a zero-value
	// configuration.
	ErrInvalidConfig = errors.New("config: invalid format")
)
