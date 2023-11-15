package security

// Config for security.
type Config struct {
	CertFile       string `yaml:"cert_file" json:"cert_file" toml:"cert_file"`
	KeyFile        string `yaml:"key_file" json:"key_file" toml:"key_file"`
	ClientCertFile string `yaml:"client_cert_file" json:"client_cert_file" toml:"client_cert_file"`
	ClientKeyFile  string `yaml:"client_key_file" json:"client_key_file" toml:"client_key_file"`
}

// IsEnabled security.
func (c *Config) IsEnabled() bool {
	return c.CertFile != "" && c.KeyFile != ""
}

// IsEnabled security.
func (c *Config) IsClientEnabled() bool {
	return c.ClientCertFile != "" && c.ClientKeyFile != ""
}
