package security

// Config for security.
type Config struct {
	Enabled  bool   `yaml:"enabled,omitempty" json:"enabled,omitempty" toml:"enabled,omitempty"`
	CertFile string `yaml:"cert_file,omitempty" json:"cert_file,omitempty" toml:"cert_file,omitempty"`
	KeyFile  string `yaml:"key_file,omitempty" json:"key_file,omitempty" toml:"key_file,omitempty"`
}
