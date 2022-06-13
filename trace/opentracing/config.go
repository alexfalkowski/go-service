package opentracing

// Config for opentracing.
type Config struct {
	Kind string `yaml:"kind"`
	Host string `yaml:"host"`
}

// IsDataDog config.
func (c *Config) IsDataDog() bool {
	return c.Kind == "datadog"
}

// IsJaeger config.
func (c *Config) IsJaeger() bool {
	return c.Kind == "jaeger"
}
