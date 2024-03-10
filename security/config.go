package security

// Config for security.
type Config struct {
	Enabled  bool   `yaml:"enabled" json:"enabled" toml:"enabled"`
	CertFile string `yaml:"cert_file" json:"cert_file" toml:"cert_file"`
	KeyFile  string `yaml:"key_file" json:"key_file" toml:"key_file"`
}
