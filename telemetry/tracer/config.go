package tracer

// Config for tracer.
type Config struct {
	Enabled bool   `yaml:"enabled" json:"enabled" toml:"enabled"`
	Kind    string `yaml:"kind" json:"kind" toml:"kind"`
	Host    string `yaml:"host" json:"host" toml:"host"`
	Key     string `yaml:"key" json:"key" toml:"key"`
	Secure  bool   `yaml:"secure" json:"secure" toml:"secure"`
}

// IsBaselime configuration.
func (c *Config) IsBaselime() bool {
	return c.Kind == "baselime"
}
