package redis

// IsEnabled for redis.
func IsEnabled(c *Config) bool {
	return c != nil
}

// Config for redis.
type Config struct {
	Compressor string            `yaml:"compressor,omitempty" json:"compressor,omitempty" toml:"compressor,omitempty"`
	Marshaller string            `yaml:"marshaller,omitempty" json:"marshaller,omitempty" toml:"marshaller,omitempty"`
	Addresses  map[string]string `yaml:"addresses,omitempty" json:"addresses,omitempty" toml:"addresses,omitempty"`
	Username   string            `yaml:"username,omitempty" json:"username,omitempty" toml:"username,omitempty"`
	Password   string            `yaml:"password,omitempty" json:"password,omitempty" toml:"password,omitempty"`
	DB         int               `yaml:"db,omitempty" json:"db,omitempty" toml:"db,omitempty"`
}
