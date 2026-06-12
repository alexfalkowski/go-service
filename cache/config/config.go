package config

import "github.com/alexfalkowski/go-service/v2/bytes"

// DefaultMaxEntries is the default maximum number of entries for bounded in-memory cache drivers.
const DefaultMaxEntries = 1024

// Config configures the cache subsystem.
type Config struct {
	// Options contains implementation-specific configuration for the selected Kind.
	//
	// The interpretation of this map depends on the cache backend implementation.
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty" toml:"options,omitempty"`

	// Kind selects the cache backend implementation.
	//
	// The built-in driver kinds are "redis" and "ttlcache". Unknown kinds cause
	// [github.com/alexfalkowski/go-service/v2/cache/driver.NewDriver] to return
	// [github.com/alexfalkowski/go-service/v2/cache/driver.ErrNotFound].
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
	// A zero value applies [bytes.DefaultSize]. Values must be between 0 and [bytes.MaxConfigSize].
	MaxSize bytes.Size `yaml:"max_size" json:"max_size" toml:"max_size" validate:"config_size"`

	// MaxEntries limits the number of entries retained by bounded in-memory cache drivers.
	MaxEntries int `yaml:"max_entries" json:"max_entries" toml:"max_entries" validate:"gt=0"`
}

// IsEnabled reports whether caching is enabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetMaxSize returns the configured cache value limit.
//
// A nil receiver or a zero value falls back to [bytes.DefaultSize].
func (c *Config) GetMaxSize() bytes.Size {
	if c == nil || c.MaxSize == 0 {
		return bytes.DefaultSize
	}

	return c.MaxSize
}

// GetMaxEntries returns the configured cache entry limit.
//
// A nil receiver or a zero value falls back to [DefaultMaxEntries].
func (c *Config) GetMaxEntries() int {
	if c == nil || c.MaxEntries == 0 {
		return DefaultMaxEntries
	}

	return c.MaxEntries
}
