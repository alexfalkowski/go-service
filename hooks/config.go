package hooks

// Config configures Standard Webhooks secret loading.
type Config struct {
	// Secrets contains named Standard Webhooks secrets trusted for verification.
	//
	// Key selects the active secret used for signing.
	// Verification accepts signatures from every configured secret.
	Secrets Secrets `yaml:"secrets,omitempty" json:"secrets,omitempty" toml:"secrets,omitempty"`

	// Key is the active secret id used for signing.
	//
	// The selected entry must exist in Secrets.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled reports whether hooks configuration is present.
//
// By convention across go-service config types, a nil *[Config] is treated as "disabled".
func (c *Config) IsEnabled() bool {
	return c != nil
}

// Secrets maps secret ids to Standard Webhooks secret source strings.
type Secrets map[string]string

// Get returns the secret with id.
func (s Secrets) Get(id string) string {
	return s[id]
}
