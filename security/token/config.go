package token

// Config for token.
type Config struct {
	Kind string `yaml:"kind" json:"kind" toml:"kind"`
}

// Generator for token.
func (c *Config) Generator() Generator {
	return genRegister[c.Kind]
}

// Verifier for token.
func (c *Config) Verifier() Verifier {
	return verRegister[c.Kind]
}
