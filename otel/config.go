package otel

// Config for otel.
type Config struct {
	Kind string `yaml:"kind" json:"kind" toml:"kind"`
	Host string `yaml:"host" json:"host" toml:"host"`
}

// IsJaeger config.
func (c *Config) IsJaeger() bool {
	return c.Kind == "jaeger"
}
