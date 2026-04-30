package access

import "github.com/alexfalkowski/go-service/v2/os"

// Config configures access control model and policy loading for the access Controller.
//
// The config is consumed by NewController, which constructs a Casbin-backed
// controller/enforcer using the configured model and policy sources.
//
// # Model and policy sources
//
// Model and Policy are resolved with os.FS.ReadSource and loaded into Casbin
// from the resolved content.
//
// # Enablement
//
// Enablement is modeled by presence: a nil *Config disables access control wiring,
// and NewController returns (nil, nil).
type Config struct {
	// Model is the access control model source string.
	//
	// It supports the go-service source string pattern implemented by `os.FS.ReadSource`:
	//   - "env:NAME" to read model content from an environment variable,
	//   - "file:/path/to/model.conf" to read model content from a file, or
	//   - any other value treated as literal model content.
	Model string `yaml:"model,omitempty" json:"model,omitempty" toml:"model,omitempty"`

	// Policy is the access control policy source string.
	//
	// It supports the go-service source string pattern implemented by `os.FS.ReadSource`:
	//   - "env:NAME" to read policy content from an environment variable,
	//   - "file:/path/to/policy.csv" to read policy content from a file, or
	//   - any other value treated as literal policy content.
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty" toml:"policy,omitempty"`
}

// IsEnabled reports whether access configuration is present.
//
// A nil receiver is considered disabled and returns false.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetModel resolves and returns the configured model content.
//
// It delegates to `fs.ReadSource(c.Model)` and returns any read/resolve error
// from that operation.
func (c *Config) GetModel(fs *os.FS) (string, error) {
	model, err := fs.ReadSource(c.Model)
	return string(model), err
}

// GetPolicy resolves and returns the configured policy content.
//
// It delegates to `fs.ReadSource(c.Policy)` and returns any read/resolve error
// from that operation.
func (c *Config) GetPolicy(fs *os.FS) (string, error) {
	policy, err := fs.ReadSource(c.Policy)
	return string(policy), err
}
