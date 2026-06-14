package time

// Config configures a network time provider.
//
// This configuration is consumed by NewNetwork, which selects a concrete Network
// implementation based on Kind and passes Address to that implementation.
//
// # Enablement
//
// Enablement is modeled by presence: a nil *[Config] indicates that network time is
// disabled. When disabled, NewNetwork returns (nil, nil). A non-nil *[Config]
// is enabled even when Kind is empty; empty or unknown kinds return
// [ErrNotFound].
//
// # Kind
//
// Kind selects the network time provider implementation. This package currently
// supports (via NewNetwork):
//
//   - "ntp": Network Time Protocol (NTP)
//   - "nts": Network Time Security (NTS)
//
// If Kind is empty or unrecognized, NewNetwork returns [ErrNotFound].
//
// # Address
//
// Address is the provider/server address passed to the selected implementation.
// The expected format is implementation-specific (for example a hostname or
// pool name). If Address is empty or invalid, the provider will typically return
// an error when Now is called.
type Config struct {
	// Kind selects the network time provider implementation (for example "ntp" or "nts").
	//
	// If Kind is empty or unknown, NewNetwork returns [ErrNotFound].
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// Address is the provider address passed to the selected implementation.
	//
	// The expected format is implementation-specific. For example, NTP may accept
	// pool hostnames, and NTS may accept hostnames of NTS-enabled servers.
	Address string `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`
}

// IsEnabled reports whether network time configuration is present.
//
// A nil receiver is considered disabled and returns false. IsEnabled does not
// validate Kind; a non-nil config with an empty Kind still returns true.
func (c *Config) IsEnabled() bool {
	return c != nil
}
