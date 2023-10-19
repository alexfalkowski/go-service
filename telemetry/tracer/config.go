package tracer

// Config for tracer.
type Config struct {
	Host   string `yaml:"host" json:"host" toml:"host"`
	Secure bool   `yaml:"secure" json:"secure" toml:"secure"`
}
