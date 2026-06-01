package config

import "github.com/alexfalkowski/go-service/v2/bytes"

// Config configures the cache subsystem.
type Config struct {
	// Options contains implementation-specific configuration for the selected Kind.
	//
	// The interpretation of this map depends on the cache backend implementation.
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty" toml:"options,omitempty"`

	// Kind selects the cache backend implementation.
	//
	// The built-in driver kinds are "redis" and "sync". Unknown kinds cause
	// cache/driver.NewDriver to return ErrNotFound.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// Compressor selects the compression algorithm used for cached values (if supported by the implementation).
	Compressor string `yaml:"compressor,omitempty" json:"compressor,omitempty" toml:"compressor,omitempty"`

	// Encoder selects the value encoding used when storing objects in the cache (if applicable).
	Encoder string `yaml:"encoder,omitempty" json:"encoder,omitempty" toml:"encoder,omitempty"`

	// MaxSize limits encoded cache value size before compression, after compression, and after
	// decompression.
	//
	// In config files it is encoded as a human-readable SI size string (for example "64B", "2MB", "4GB").
	//
	// A zero value applies bytes.DefaultSize. Values must be between 0 and bytes.MaxConfigSize.
	MaxSize bytes.Size `yaml:"max_size,omitempty" json:"max_size,omitempty" toml:"max_size,omitempty" validate:"config_size"`
}

// IsEnabled reports whether caching is enabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetMaxSize returns the configured cache value limit.
//
// A nil receiver or a zero value falls back to bytes.DefaultSize.
func (c *Config) GetMaxSize() bytes.Size {
	if c == nil || c.MaxSize == 0 {
		return bytes.DefaultSize
	}

	return c.MaxSize
}
