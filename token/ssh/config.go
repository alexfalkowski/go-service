package ssh

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/types/slices"
)

// Config for ssh.
type Config struct {
	Key  *Key `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	Keys Keys `yaml:"keys,omitempty" json:"keys,omitempty" toml:"keys,omitempty"`
}

// IsEnabled for ssh.
func (c *Config) IsEnabled() bool {
	return c != nil && (c.Key != nil || c.Keys != nil)
}

// Key configuration with a name and ssh key.
type Key struct {
	*ssh.Config `yaml:",inline" json:",inline" toml:",inline"`
	Name        string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
}

// Keys configuration with a names and ssh keys.
type Keys []*Key

// Get a key by name.
func (c Keys) Get(name string) *Key {
	key, _ := slices.ElemFunc(c, func(k *Key) bool { return k.Name == name })

	return key
}
