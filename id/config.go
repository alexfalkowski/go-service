package id

// Config configures ID generation for a go-service based service.
//
// Services typically set Kind to select which generator implementation to use (for example "uuid",
// "ksuid", "ulid", "nanoid", or "xid"), depending on which generators are compiled/registered by the
// service.
//
// # Optional pointers and "enabled" semantics
//
// This config is intentionally optional. By convention across go-service configuration types, a nil
// *Config is treated as "ID generation disabled". When disabled, wiring that depends on this config
// often returns (nil, nil) rather than failing.
type Config struct {
	// Kind selects the ID generator implementation.
	//
	// The set of supported kinds depends on what has been wired into the application (see the id module
	// wiring). If Kind is set to an unknown value, generator selection typically returns ErrNotFound.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsEnabled reports whether ID configuration is present.
//
// By convention, a nil *Config is treated as "disabled".
func (c *Config) IsEnabled() bool {
	return c != nil
}
