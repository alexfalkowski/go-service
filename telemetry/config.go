package telemetry

// Config for telemetry.
type Config struct {
	Host   string `yaml:"host" json:"host" toml:"host"`
	Secure bool   `yaml:"secure" json:"secure" toml:"secure"`
}
