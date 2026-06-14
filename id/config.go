package id

// Config configures ID generation for a go-service based service.
//
// Services typically set Kind to select which generator implementation to use (for example "uuid",
// "ksuid", "ulid", "nanoid", or "xid"), depending on which generators are compiled/registered by the
// service.
//
// # Optional pointers and "enabled" semantics
//
// This config is intentionally optional. Only a nil *[Config] falls back to the "uuid" kind in the
// generator registry; an enabled config with an empty Kind is selected as the empty kind and typically
// returns [ErrNotFound].
type Config struct {
	// Kind selects the ID generator implementation.
	//
	// The set of supported kinds depends on what has been wired into the application (see [Module]).
	// If Kind is empty or unknown, generator selection typically returns [ErrNotFound].
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`
}

// IsEnabled reports whether ID configuration is present.
//
// By convention, a nil *[Config] is treated as "disabled".
func (c *Config) IsEnabled() bool {
	return c != nil
}
