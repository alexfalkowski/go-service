package config

// Config configures the cache subsystem.
type Config struct {
	// Options contains implementation-specific configuration for the selected Kind.
	//
	// The interpretation of this map depends on the cache backend implementation.
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty" toml:"options,omitempty"`

	// Kind selects the cache backend implementation (for example "redis", "valkey", or "noop"),
	// depending on which implementations are compiled/registered by the service.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// Compressor selects the compression algorithm used for cached values (if supported by the implementation).
	Compressor string `yaml:"compressor,omitempty" json:"compressor,omitempty" toml:"compressor,omitempty"`

	// Encoder selects the value encoding used when storing objects in the cache (if applicable).
	Encoder string `yaml:"encoder,omitempty" json:"encoder,omitempty" toml:"encoder,omitempty"`
}

// IsEnabled for cache.
func (c *Config) IsEnabled() bool {
	return c != nil
}
