package access

// Config configures access control policy loading for the access Controller.
//
// The policy is consumed by NewController, which constructs a Casbin-backed
// controller/enforcer using the embedded ModelConfig and a policy adapter.
//
// # Source string convention
//
// Policy supports the go-service “source string” pattern:
//   - "env:NAME" to read from an environment variable
//   - "file:/path" to read from a file
//   - otherwise treated as the literal policy content
//
// Note: This package does not resolve source strings itself. If you configure Policy
// as a source string, resolve it before calling NewController (for example by using
// os.FS.ReadSource elsewhere in your config projection/wiring).
//
// # Enablement
//
// Enablement is modeled by presence: a nil *Config disables access control wiring,
// and NewController returns (nil, nil).
type Config struct {
	// Policy is the access control policy document source.
	//
	// It may be a go-service “source string” (see Config docs) or a literal policy.
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty" toml:"policy,omitempty"`
}

// IsEnabled reports whether access configuration is present.
//
// A nil receiver is considered disabled and returns false.
func (c *Config) IsEnabled() bool {
	return c != nil
}
