package opentracing

// Config for opentracing.
type Config struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
}

// IsDataDog config.
func (c *Config) IsDataDog() bool {
	return c.Type == "datadog"
}

// IsJaeger config.
func (c *Config) IsJaeger() bool {
	return c.Type == "jaeger"
}
