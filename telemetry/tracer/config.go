package tracer

// Config for otel.
type Config struct {
	Host   string `yaml:"host" json:"host" toml:"host"`
	Secure bool   `yaml:"secure" json:"secure" toml:"secure"`
}
