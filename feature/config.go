package feature

// Config for feature.
type Config struct {
	Kind string `yaml:"kind" json:"kind" toml:"kind"`
	Host string `yaml:"host" json:"host" toml:"host"`
}
