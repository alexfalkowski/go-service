package hooks

import "github.com/alexfalkowski/go-service/v2/time"

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

	// Leeway is the optional clock-skew tolerance applied to the Webhook-Timestamp
	// freshness check during verification.
	//
	// A zero value keeps the Standard Webhooks library's fixed 5-minute freshness
	// window. A non-zero value replaces that fixed window with this configured
	// tolerance, matching the clock-skew Leeway already exposed by the JWT,
	// PASETO, and SSH token verifiers.
	Leeway time.Duration `yaml:"leeway,omitempty" json:"leeway,omitempty" toml:"leeway,omitempty" validate:"omitempty,duration_second_precision"`
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
