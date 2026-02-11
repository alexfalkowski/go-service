package ssh

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/types/slices"
)

// Config configures SSH token key material.
//
// It supports a single signing key (Key) and a set of verification keys (Keys).
type Config struct {
	// Key is the signing key configuration used to mint SSH tokens.
	Key *Key `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`

	// Keys is the set of verification keys that may be used to validate SSH tokens.
	Keys Keys `yaml:"keys,omitempty" json:"keys,omitempty" toml:"keys,omitempty"`
}

// IsEnabled reports whether SSH token configuration is present and contains at least one key.
func (c *Config) IsEnabled() bool {
	return c != nil && (c.Key != nil || c.Keys != nil)
}

// Key describes an SSH key configuration along with its logical name.
//
// The embedded `ssh`.Config provides the key material source/loading configuration.
type Key struct {
	// Config contains the SSH key material configuration (public/private key sources).
	*ssh.Config `yaml:",inline" json:",inline" toml:",inline"`

	// Name is the logical key name used to select a key (for example via Keys.Get).
	Name string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
}

// Keys is a list of named SSH keys.
type Keys []*Key

// Get returns the key with the given name, or nil if no matching key exists.
func (c Keys) Get(name string) *Key {
	key, _ := slices.ElemFunc(c, func(k *Key) bool { return k.Name == name })

	return key
}
