package access

// Config configures access control policy loading for the access Controller.
//
// The policy is consumed by NewController, which constructs a Casbin-backed
// controller/enforcer using the embedded ModelConfig and a policy adapter.
//
// # Policy path
//
// Policy is passed directly to Casbin's file adapter, so it must be a real
// filesystem path. This package does not resolve go-service source strings such
// as "env:" or "file:", and it does not accept literal policy payloads.
//
// # Enablement
//
// Enablement is modeled by presence: a nil *Config disables access control wiring,
// and NewController returns (nil, nil).
type Config struct {
	// Policy is the access control policy file path.
	//
	// The path is passed directly to Casbin's file adapter.
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty" toml:"policy,omitempty"`
}

// IsEnabled reports whether access configuration is present.
//
// A nil receiver is considered disabled and returns false.
func (c *Config) IsEnabled() bool {
	return c != nil
}
