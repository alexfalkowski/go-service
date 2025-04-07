package ssh

import (
	"slices"

	"github.com/alexfalkowski/go-service/crypto/ssh"
)

// IsEnabled for ssh.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && (cfg.Key != nil || cfg.Keys != nil)
}

// Config for ssh.
type Config struct {
	Key  *Key `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
	Keys Keys `yaml:"keys,omitempty" json:"keys,omitempty" toml:"keys,omitempty"`
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
	index := slices.IndexFunc(c, func(k *Key) bool { return k.Name == name })
	if index == -1 {
		return nil
	}

	return c[index]
}
