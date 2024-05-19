package hooks

import (
	"encoding/base64"

	"github.com/alexfalkowski/go-service/crypto/rand"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// Generate for hooks.
func Generate() (Secret, error) {
	s, err := rand.GenerateBytes(32)

	return Secret(base64.StdEncoding.EncodeToString(s)), err
}

// New hook from config.
func New(cfg *Config) (*hooks.Webhook, error) {
	if cfg == nil {
		return hooks.NewWebhookRaw(nil)
	}

	return hooks.NewWebhook(string(cfg.Secret))
}
